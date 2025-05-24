package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"backend/middleware"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type DownloadShareClaims struct {
	ShareId string `json:"share_id"`
	jwt.RegisteredClaims
}

func DownloadShare(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	token := cc.FormValue("token")
	if token == "" {
		return utils.HTTPErrorHandler(c, errors.New("缺少token"))
	}
	claims := DownloadShareClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetEnv("download_secret")), nil
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if !t.Valid {
		return utils.HTTPErrorHandler(c, errors.New("token格式错误"))
	}
	shareInfo, _ := models.GetRedisShareInfo(claims.ShareId)

	if shareInfo.Type == models.ShareTypeFile {
		fileInfo, _ := models.GetRedisFileInfo(shareInfo.Data)
		uploadPath, err := utils.GetUploadDirPath()
		if err != nil {
			return err
		}
		return cc.Attachment(fmt.Sprintf("%s/%s", uploadPath, utils.GetFileId(fileInfo.FileHash, fileInfo.FileSize)), shareInfo.FileName)
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"data": shareInfo.Data,
	})
}

type VaildateShareProps struct {
	ShareId  string `json:"share_id"`
	Password string `json:"password"`
}

func VaildateShare(c echo.Context) error {
	cc := c.(*middleware.CustomContext)

	r := new(VaildateShareProps)
	if err := cc.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.ShareId == "" {
		return utils.HTTPErrorHandler(c, errors.New("缺少分享ID"))
	}

	shareInfo, err := models.GetRedisShareInfo(r.ShareId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo == nil {
		return utils.HTTPErrorHandler(c, errors.New("分享不存在"))
	}
	if shareInfo.Password != "" && shareInfo.Password != r.Password {
		return utils.HTTPErrorHandler(c, errors.New("分享密码错误"))
	}
	// 如果下载次数为0，则设置为-1 防止空值问题
	if shareInfo.ViewNum < 1 {
		return utils.HTTPErrorHandler(c, errors.New("下载次数不足"))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, DownloadShareClaims{
		ShareId: r.ShareId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	downloadToken, err := token.SignedString([]byte(utils.GetEnv("download_secret")))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo.Type == models.ShareTypeFile {
		fileInfo, err := models.GetRedisFileInfo(shareInfo.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if fileInfo == nil {
			return utils.HTTPErrorHandler(c, errors.New("分享文件不存在"))
		}
		if fileInfo.FileType != models.FileTypeUpload {
			return utils.HTTPErrorHandler(c, errors.New("分享文件状态错误"))
		}
	}
	// download_nums 必须放在创建token的时候减掉，不然多线程下载会导致多次减掉
	latestViewNum := shareInfo.ViewNum - 1
	// 如果下载次数为0，则设置为-1 防止空值问题
	if latestViewNum < 1 {
		latestViewNum = -1
	}
	models.SetRedisShareInfo(r.ShareId, models.RedisShareInfo{
		ViewNum: latestViewNum,
	})
	if shareInfo.Type == models.ShareTypeFile {
		return utils.HTTPSuccessHandler(c, map[string]any{
			"token": downloadToken,
		})
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"token": downloadToken,
	})
}
