package controllers

import (
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	filemodel "pkg/models/file"
	sharemodel "pkg/models/share"
	statmodel "pkg/models/stat"
	u "pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
)

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
	claims, err := services.ParseDownloadToken(token)
	if err != nil {
		if errors.Is(err, services.ErrInvalidDownloadToken) {
			err = ErrInvalidRequest
		}
		return utils.HTTPErrorHandler(c, err)
	}
	shareInfo, err := sharemodel.GetRedisShareInfo(claims.ShareId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo == nil {
		return utils.HTTPErrorHandler(c, ErrShareNotFound)
	}
	if shareInfo.Type == sharemodel.ShareTypeFile {
		if len(fileIds) > 0 {
			shareInfo.Files = lo.Filter(shareInfo.Files, func(file sharemodel.ShareFileData, _ int) bool {
				return lo.ContainsBy(fileIds, func(fileID string) bool {
					return fileID == utils.PublicFileID(file.Id)
				})
			})
			if len(shareInfo.Files) != len(fileIds) {
				return utils.HTTPErrorHandler(c, ErrInvalidShareFileData)
			}
		}
		if len(shareInfo.Files) == 0 {
			return utils.HTTPErrorHandler(c, ErrInvalidShareFileData)
		}
		root, err := u.GetUploadDir()
		if err != nil {
			return err
		}
		defer root.Close() //nolint:errcheck
		for index, file := range shareInfo.Files {
			fileInfo, err := filemodel.GetRedisFileInfo(file.Id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			if fileInfo == nil {
				return utils.HTTPErrorHandler(c, ErrShareFileNotFound)
			}
			shareInfo.Files[index].Id, err = u.GetFileId(fileInfo.FileHash, fileInfo.FileSize)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
		}
		if len(shareInfo.Files) == 1 {
			return c.Attachment(shareInfo.Files[0].Id, shareInfo.Files[0].FileName)
		}
		target := c.FormValue("target")
		if !lo.Contains([]string{"zip", "tar.gz", "tar.zst", "tar.s2", "tar.snappy"}, target) {
			target = "zip"
		}
		compressFullName, err := services.GenerateCompressFiles(claims.ShareId, shareInfo.Files, root, target)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		defer root.Remove(compressFullName) //nolint:errcheck
		return c.Attachment(compressFullName, fmt.Sprintf("%s.%s", claims.ShareId, target))
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
	if shareInfo == nil || shareInfo.ExpireAt <= time.Now().Unix() {
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
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if shareInfo == nil {
			return utils.HTTPErrorHandler(c, ErrShareNotFound)
		}
		if shareInfo.ExpireAt <= time.Now().Unix() {
			return utils.HTTPErrorHandler(c, ErrShareNotFound)
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
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, services.DownloadShareClaims{
			ShareId: r.ShareId,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.GetConfig().Share.DownloadWindow) * time.Hour)),
			},
		})

		// Sign and get the complete encoded token as a string using the secret
		downloadToken, err := token.SignedString([]byte(u.GetConfig().Share.DownloadSecret))
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

		// 统计下载数
		statData, err := statmodel.DomainDownloadStatDataToJson(&statmodel.DownloadStatData{
			ShareID: r.ShareId,
			IP:      c.RealIP(),
		})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		err = statmodel.AppendRedisStat(statmodel.DownloadZSetKey, time.Now().Unix(), statData)
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
