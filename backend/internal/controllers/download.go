package controllers

import (
	"backend/internal/utils"
	"fmt"
	"pkg/models"
	u "pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/spf13/cast"
)

type DownloadShareClaims struct {
	ShareId string `json:"share_id"`
	jwt.RegisteredClaims
}

func DownloadShare(c *echo.Context) error {
	token := c.FormValue("token")
	if token == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	claims := DownloadShareClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.GetEnv("share.download_secret")), nil
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if !t.Valid {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	shareInfo, _ := models.GetRedisShareInfo(claims.ShareId)
	if shareInfo == nil {
		return utils.HTTPErrorHandler(c, ErrShareNotFound)
	}

	if shareInfo.Type == models.ShareTypeFile {
		fileInfo, _ := models.GetRedisFileInfo(shareInfo.Data)
		uploadPath, err := u.GetUploadDirPath()
		if err != nil {
			return err
		}
		return c.Attachment(fmt.Sprintf("%s/%s", uploadPath, u.GetFileId(fileInfo.FileHash, fileInfo.FileSize)), shareInfo.FileName)
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"data": shareInfo.Data,
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

	shareInfo, err := models.GetRedisShareInfo(r.ShareId)
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
	// 如果下载次数为0，则设置为-1 防止空值问题
	if shareInfo.ViewNum < 1 {
		return utils.HTTPErrorHandler(c, ErrInsufficientDownloadQuota)
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
	if shareInfo.Type == models.ShareTypeFile {
		fileInfo, err := models.GetRedisFileInfo(shareInfo.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if fileInfo == nil {
			return utils.HTTPErrorHandler(c, ErrShareFileNotFound)
		}
		if fileInfo.FileType != models.FileTypeUpload {
			return utils.HTTPErrorHandler(c, ErrInvalidShareFileState)
		}
	}
	// download_nums 必须放在创建token的时候减掉，不然多线程下载会导致多次减掉
	latestViewNum := shareInfo.ViewNum - 1
	// 如果下载次数为0，则设置为-1 防止空值问题
	if latestViewNum < 1 {
		latestViewNum = -1
	}
	err = models.SetRedisShareInfo(r.ShareId, models.RedisShareInfo{
		ViewNum: latestViewNum,
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	// 统计分享数
	currentDate := time.Now().Format("2006-01-02")
	err = models.SetRedisStat(currentDate, func(stat *models.StatData) *models.StatData {
		stat.DownloadNum += 1
		return stat
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if shareInfo.Type == models.ShareTypeFile {
		return utils.HTTPSuccessHandler(c, map[string]any{
			"token": downloadToken,
		})
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"token": downloadToken,
	})
}
