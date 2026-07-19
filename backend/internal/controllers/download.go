package controllers

import (
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	filemodel "pkg/models/file"
	sharemodel "pkg/models/share"
	statmodel "pkg/models/stat"
	u "pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type DownloadShareClaims struct {
	ShareId string `json:"share_id"`
	jwt.RegisteredClaims
}

func DownloadShare(c *echo.Context) error {
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return err
	}
	token := req.Form.Get("token")
	fileIds := req.Form["file_ids"]
	if token == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	claims := DownloadShareClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.GetEnv("share.download_secret")), nil
	})
	if err != nil || !t.Valid {
		return utils.HTTPErrorHandler(c, lo.Ternary(err != nil, err, ErrInvalidRequest))
	}
	if claims.ShareId == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	shareInfo, err := sharemodel.GetRedisShareInfo(claims.ShareId)
	if err != nil || shareInfo == nil {
		return utils.HTTPErrorHandler(c, lo.Ternary(err != nil, err, ErrShareNotFound))
	}
	if shareInfo.Type == sharemodel.ShareTypeFile {
		shareFiles := shareInfo.Files
		if len(fileIds) > 0 {
			shareFiles = lo.Filter(shareFiles, func(file sharemodel.ShareFileData, _ int) bool {
				return lo.Contains(fileIds, file.Id)
			})
			if len(shareFiles) != len(fileIds) {
				return utils.HTTPErrorHandler(c, ErrInvalidShareFileData)
			}
		}
		if len(shareFiles) == 0 {
			return utils.HTTPErrorHandler(c, ErrInvalidShareFileData)
		}
		uploadPath, err := u.GetUploadDirPath()
		if err != nil {
			return err
		}
		for index, file := range shareFiles {
			fileInfo, err := filemodel.GetRedisFileInfo(file.Id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			if fileInfo == nil {
				return utils.HTTPErrorHandler(c, ErrShareFileNotFound)
			}
			shareFiles[index].Id = u.GetFileId(fileInfo.FileHash, fileInfo.FileSize)
		}
		if len(shareFiles) == 1 {
			return c.Attachment(fmt.Sprintf("%s/%s", uploadPath, shareFiles[0].Id), shareFiles[0].FileName)
		}
		target := c.FormValue("target")
		if !lo.Contains([]string{"zip", "tar.gz"}, target) {
			target = "zip"
		}
		compressPath, err := services.GenerateCompressFiles(claims.ShareId, shareFiles, uploadPath, target)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		defer os.Remove(compressPath) //nolint:errcheck
		return c.Attachment(compressPath, fmt.Sprintf("%s.%s", claims.ShareId, target))
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"text": shareInfo.Text,
	})
}

type VaildateShareProps struct {
	ShareId  string `json:"share_id"`
	Password string `json:"password"`
}

func VaildateShare(c *echo.Context) error {
	r := new(VaildateShareProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.ShareId == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	shareInfo, err := sharemodel.GetRedisShareInfo(r.ShareId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo == nil {
		return utils.HTTPErrorHandler(c, ErrShareNotFound)
	}
	if shareInfo.Password != "" {
		if r.Password == "" {
			return utils.HTTPErrorHandler(c, ErrInvalidRequest)
		}
		hash, err := utils.GeneratePasswordHash(r.Password)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if hash != shareInfo.Password {
			return utils.HTTPErrorHandler(c, ErrInvalidSharePassword)
		}
	}
	return u.WithLocker(context.Background(), "015:shareInfoMap:"+r.ShareId, 0, func(ctx context.Context) error {
		shareInfo, err := sharemodel.GetRedisShareInfo(r.ShareId)
		if err != nil || shareInfo == nil {
			return utils.HTTPErrorHandler(c, lo.Ternary(err != nil, err, ErrShareNotFound))
		}
		if shareInfo.ViewNum < 1 {
			return utils.HTTPErrorHandler(c, ErrInsufficientDownloadQuota)
		}
		if shareInfo.Type == sharemodel.ShareTypeFile {
			for _, file := range shareInfo.Files {
				fileInfo, err := filemodel.GetRedisFileInfo(file.Id)
				if err != nil {
					return utils.HTTPErrorHandler(c, err)
				}
				if fileInfo == nil {
					return utils.HTTPErrorHandler(c, ErrShareFileNotFound)
				}
				if fileInfo.FileType != filemodel.FileTypeUpload {
					return utils.HTTPErrorHandler(c, ErrInvalidShareFileState)
				}
			}
		}
		downloadWindow := u.GetEnvWithDefault("share.download_window", "12")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, DownloadShareClaims{
			ShareId: r.ShareId,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(cast.ToDuration(downloadWindow + "h"))),
			},
		})

		// Sign and get the complete encoded token as a string using the secret
		downloadToken, err := token.SignedString([]byte(u.GetEnv("share.download_secret")))
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		// download_nums 必须放在创建token的时候减掉，不然多线程下载会导致多次减掉
		_, err = sharemodel.SetRedisShareInfo(r.ShareId, func(shareInfo *sharemodel.RedisShareInfo) *sharemodel.RedisShareInfo {
			shareInfo.ViewNum -= 1
			return shareInfo
		})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}

		// 统计分享数
		currentDate := time.Now().Format("2006-01-02")
		_, err = statmodel.SetRedisStat(currentDate, func(stat *statmodel.StatData) *statmodel.StatData {
			stat.DownloadNum += 1
			return stat
		})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}

		if len(shareInfo.NotifyEmails) > 0 || len(shareInfo.NotifyWebhooks) > 0 {
			payload, err := json.Marshal(map[string]string{
				"share_id": r.ShareId,
				"ip":       c.RealIP(),
			})
			if err == nil {
				_, _ = u.GetQueueClient().Enqueue(asynq.NewTask("share:notify", payload))
			}
		}

		if shareInfo.Type == sharemodel.ShareTypeFile {
			return utils.HTTPSuccessHandler(c, map[string]any{
				"token": downloadToken,
			})
		}
		return utils.HTTPSuccessHandler(c, map[string]any{
			"token": downloadToken,
		})
	})
}
