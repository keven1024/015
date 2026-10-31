package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// 测试切片列表的排序、空目录、非法文件名和路径穿越处理
func TestGetFileSliceList(t *testing.T) {
	uploadPath := t.TempDir()
	root, err := os.OpenRoot(uploadPath)
	require.NoError(t, err)
	defer root.Close() //nolint:errcheck
	fileId := "file-2"
	slicesPath := filepath.Join(uploadPath, fileId+"_tmp")
	require.NoError(t, os.MkdirAll(slicesPath, 0755))

	// 乱序创建切片文件，验证返回结果按索引升序排列
	for _, name := range []string{"2", "0", "10", "1"} {
		require.NoError(t, os.WriteFile(filepath.Join(slicesPath, name), []byte(name), 0600))
	}

	list, err := GetFileSliceList(fileId, root)
	require.NoError(t, err)
	require.Equal(t, []int{0, 1, 2, 10}, list)

	// 目录不存在时应返回空列表而非错误
	list, err = GetFileSliceList("not-exist", root)
	require.NoError(t, err)
	require.Empty(t, list)

	// 非法切片文件名应被忽略
	badPath := filepath.Join(uploadPath, "bad-file_tmp")
	require.NoError(t, os.MkdirAll(badPath, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(badPath, "abc"), []byte("x"), 0600))
	list, err = GetFileSliceList("bad-file", root)
	require.NoError(t, err)
	require.Empty(t, list)

	// 文件 ID 路径穿越时应返回错误
	_, err = GetFileSliceList("../outside", root)
	require.Error(t, err)
}
