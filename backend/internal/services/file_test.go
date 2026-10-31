package services

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	sharemodel "pkg/models/share"
	"strings"
	"testing"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
)

// 测试生成不同格式的压缩文件并验证解压内容
func TestGenerateCompressedTarFiles(t *testing.T) {
	uploadPath := t.TempDir()
	root, err := os.OpenRoot(uploadPath)
	require.NoError(t, err)
	defer root.Close() //nolint:errcheck
	require.NoError(t, os.WriteFile(filepath.Join(uploadPath, "stored-file"), []byte("archive content"), 0600))
	files := []sharemodel.ShareFileData{{Id: "stored-file", FileName: "download.txt"}}

	tests := []struct {
		name      string
		target    string
		newReader func(io.Reader) (io.Reader, error)
	}{
		{
			name:   "zstd",
			target: "tar.zst",
			newReader: func(r io.Reader) (io.Reader, error) {
				return zstd.NewReader(r)
			},
		},
		{
			name:   "s2",
			target: "tar.s2",
			newReader: func(r io.Reader) (io.Reader, error) {
				return s2.NewReader(r), nil
			},
		},
		{
			name:   "snappy",
			target: "tar.snappy",
			newReader: func(r io.Reader) (io.Reader, error) {
				return snappy.NewReader(r), nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			archiveName, err := GenerateCompressFiles("share", files, root, test.target)
			require.NoError(t, err)
			require.True(t, strings.HasSuffix(archiveName, "."+test.target))

			archiveFile, err := os.Open(filepath.Join(uploadPath, archiveName))
			require.NoError(t, err)
			defer archiveFile.Close() //nolint:errcheck

			reader, err := test.newReader(archiveFile)
			require.NoError(t, err)
			tarReader := tar.NewReader(reader)
			header, err := tarReader.Next()
			require.NoError(t, err)
			require.Equal(t, "download.txt", header.Name)
			content, err := io.ReadAll(tarReader)
			require.NoError(t, err)
			require.Equal(t, "archive content", string(content))
		})
	}
}

// 测试文件切片的创建和重复写入
func TestCreateFileSlice(t *testing.T) {
	uploadPath := t.TempDir()
	root, err := os.OpenRoot(uploadPath)
	require.NoError(t, err)
	defer root.Close() //nolint:errcheck
	fileId := "file-1"

	filePath, err := CreateFileSlice(fileId, root, bytes.NewReader([]byte("slice-0")), 0)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(uploadPath, fileId+"_tmp"), filePath)

	content, err := os.ReadFile(filepath.Join(filePath, "0"))
	require.NoError(t, err)
	require.Equal(t, "slice-0", string(content))

	// 再次写入不同索引，目录应已存在且可继续写入
	_, err = CreateFileSlice(fileId, root, bytes.NewReader([]byte("slice-1")), 1)
	require.NoError(t, err)
	content, err = os.ReadFile(filepath.Join(filePath, "1"))
	require.NoError(t, err)
	require.Equal(t, "slice-1", string(content))
}

// 测试文件切片合并及临时目录清理
func TestMergeFileSlices(t *testing.T) {
	uploadPath := t.TempDir()
	root, err := os.OpenRoot(uploadPath)
	require.NoError(t, err)
	defer root.Close() //nolint:errcheck
	fileId := "file-3"
	slicesPath := filepath.Join(uploadPath, fileId+"_tmp")
	require.NoError(t, os.MkdirAll(slicesPath, 0755))

	slices := map[string]string{"0": "hello ", "1": "merge ", "2": "world"}
	for name, content := range slices {
		require.NoError(t, os.WriteFile(filepath.Join(slicesPath, name), []byte(content), 0600))
	}

	err = MergeFileSlices(fileId, root, []int{0, 1, 2})
	require.NoError(t, err)

	merged, err := os.ReadFile(filepath.Join(uploadPath, fileId))
	require.NoError(t, err)
	require.Equal(t, "hello merge world", string(merged))

	// 合并完成后临时切片目录应被清理
	_, err = os.Stat(slicesPath)
	require.True(t, os.IsNotExist(err))
}

// 测试文件 ID 路径穿越时被拒绝
func TestFileIDPathTraversal(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)
	defer root.Close() //nolint:errcheck

	fileId := "../outside"
	_, err = CreateFileSlice(fileId, root, bytes.NewReader([]byte("data")), 0)
	require.Error(t, err)

	err = MergeFileSlices(fileId, root, []int{0})
	require.Error(t, err)
}
