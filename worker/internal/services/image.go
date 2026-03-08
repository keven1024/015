package services

import (
	"errors"
	u "pkg/utils"
	"worker/internal/utils"

	"github.com/samber/lo"
)

var (
	ErrUnsupportedMimeType = errors.New("UnsupportedMimeType")
)

func CompressImage(filePath string, mimeType string) (string, error) {
	compressedPath := filePath + "_compressed"
	switch mimeType {
	case "image/png":
		_, err := utils.RunCommand("pngquant", "--output", compressedPath, filePath)
		if err != nil {
			return "", err
		}
	case "image/jpeg":
		err := u.CopyFile(filePath, compressedPath)
		if err != nil {
			return "", err
		}
		_, err = utils.RunCommand("jpegoptim", "-m", "80", "--strip-all", compressedPath)
		if err != nil {
			return "", err
		}
	default:
		return "", ErrUnsupportedMimeType
	}
	return compressedPath, nil
}

// 支持的图片扩展名列表
var supportedImageExts = []string{
	"jpg", "jpeg", "png", "gif", "webp",
}

func ConvertImageWithMagick(filePath, mimeType, targetExt string) (string, error) {
	// 验证目标扩展名是否合法
	if !lo.Contains(supportedImageExts, targetExt) {
		return "", ErrUnsupportedMimeType
	}

	outputPath := filePath + "_converted." + targetExt

	// JPG 不支持透明通道，透明 PNG 转 JPG 前需压平到背景色（白底），否则透明处会变黑
	args := []string{filePath}
	if lo.Contains([]string{"jpg", "jpeg"}, targetExt) {
		args = append(args, "-background", "white", "-flatten")
	}
	args = append(args, outputPath)

	_, err := utils.RunCommand("convert", args...)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
