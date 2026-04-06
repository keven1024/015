package services

import (
	"fmt"
	"os"
	"path/filepath"
	"pkg/utils"
	"runtime"
	"testing"
	"worker/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCompressPNGHappyPath(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "test.png")
	// 从 test/resource 复制真实 PNG，路径基于当前测试文件位置
	_, self, _, _ := runtime.Caller(0)
	srcPath := filepath.Join(filepath.Dir(self), "..", "resource", "test.png")
	err := utils.CopyFile(srcPath, filePath)
	if err != nil {
		t.Fatal(err)
	}
	got, err := services.CompressImage(filePath, "image/png")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, filePath+"_compressed")
	origInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	compInfo, err := os.Stat(got)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, origInfo.Size(), compInfo.Size())
	fmt.Printf("原图: %d | 压缩后: %d | 压缩率: %f%%\n", origInfo.Size(), compInfo.Size(), float64(compInfo.Size())/float64(origInfo.Size())*100)
}

func TestCompressJPEGHappyPath(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "test.jpg")
	_, self, _, _ := runtime.Caller(0)
	srcPath := filepath.Join(filepath.Dir(self), "..", "resource", "test.jpg")
	err := utils.CopyFile(srcPath, filePath)
	if err != nil {
		t.Fatal(err)
	}
	got, err := services.CompressImage(filePath, "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, filePath+"_compressed")
	origInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	compInfo, err := os.Stat(got)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, origInfo.Size(), compInfo.Size())
	fmt.Printf("原图: %d | 压缩后: %d | 压缩率: %f%%\n", origInfo.Size(), compInfo.Size(), float64(compInfo.Size())/float64(origInfo.Size())*100)
}

func TestConvertImageWithMagickA2B(t *testing.T) {
	tmp := t.TempDir()
	testList := [][]string{
		{"jpg", "png"},
		{"png", "jpg"},
		{"jpg", "webp"},
		{"png", "webp"},
		{"webp", "jpg"},
		{"webp", "png"},
	}
	_, self, _, _ := runtime.Caller(0)
	for _, test := range testList {
		filePath := filepath.Join(tmp, uuid.New().String())
		srcPath := filepath.Join(filepath.Dir(self), "..", "resource", fmt.Sprintf("test.%s", test[0]))
		err := utils.CopyFile(srcPath, filePath)
		if err != nil {
			t.Fatal(err)
		}
		got, err := services.ConvertImageWithMagick(filePath, fmt.Sprintf("image/%s", test[0]), test[1])
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, got, filePath+"_converted."+test[1])
		info, err := os.Stat(got)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, info.Size() > 0, "转换后的文件应该有内容")
		fmt.Printf("%s -> %s 转换成功: %s (大小: %d bytes)\n", test[0], test[1], filepath.Base(got), info.Size())
	}
}

func TestConvertImageWithMagickInvalidExt(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "test.png")
	_, self, _, _ := runtime.Caller(0)
	srcPath := filepath.Join(filepath.Dir(self), "..", "resource", "test.png")
	err := utils.CopyFile(srcPath, filePath)
	if err != nil {
		t.Fatal(err)
	}

	// 测试非法扩展名（防止注入）
	_, err = services.ConvertImageWithMagick(filePath, "image/png", "exe")
	assert.Error(t, err, "应该返回错误")
	assert.Equal(t, services.ErrUnsupportedMimeType, err, "应该返回 ErrUnsupportedMimeType 错误")
}
