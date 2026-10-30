package services

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	sharemodel "pkg/models/share"
	"sort"
	"time"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func CreateFileSlice(fileId string, root *os.Root, fileSlice io.Reader, fileIndex int64) (string, error) {
	filePath := fmt.Sprintf("%s_%s", fileId, "tmp")
	if err := root.MkdirAll(filePath, 0755); err != nil {
		return "", err
	}

	dst, err := root.Create(filepath.Join(filePath, fmt.Sprintf("%d", fileIndex)))
	if err != nil {
		return "", err
	}
	defer dst.Close() //nolint:errcheck

	if _, err = io.Copy(dst, fileSlice); err != nil {
		return "", err
	}

	return filepath.Join(root.Name(), filePath), nil
}

func GetFileSliceList(fileId string, root *os.Root) ([]int, error) {
	files, err := fs.ReadDir(root.FS(), fmt.Sprintf("%s_%s", fileId, "tmp"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %w", ErrReadFileSliceDir, err)
	}
	fileSliceList := lo.FilterMap(files, func(file os.DirEntry, _ int) (int, bool) {
		index, err := cast.ToIntE(file.Name())
		return index, err == nil
	})
	sort.Ints(fileSliceList)
	return fileSliceList, nil
}

// MergeFileSlices 合并文件切片
func MergeFileSlices(fileId string, root *os.Root, fileSliceList []int) error {
	slicesPath := fmt.Sprintf("%s_%s", fileId, "tmp")
	defer root.RemoveAll(slicesPath) //nolint:errcheck
	// 创建最终文件
	destFile, err := root.Create(fileId)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateMergeFile, err)
	}
	defer destFile.Close() //nolint:errcheck

	// 合并文件
	buffer := make([]byte, 4*1024*1024) // 4MB buffer
	if !lo.EveryBy(fileSliceList, func(index int) bool {
		sliceFilePath := filepath.Join(slicesPath, fmt.Sprintf("%d", index))
		sf, err := root.Open(sliceFilePath)
		if err != nil {
			zap.L().Error("failed to open file slice", zap.Int("index", index), zap.Error(err))
			return false
		}

		if _, err := io.CopyBuffer(destFile, sf, buffer); err != nil {
			zap.L().Error("failed to merge file slices", zap.Int("index", index), zap.Error(err))
			if closeErr := sf.Close(); closeErr != nil {
				zap.L().Error("failed to close file slice", zap.Int("index", index), zap.Error(closeErr))
			}
			return false
		}
		if err := sf.Close(); err != nil {
			zap.L().Error("failed to close file slice", zap.Int("index", index), zap.Error(err))
			return false
		}
		return true
	}) {
		return ErrMergeFileSlices
	}
	return nil
}

func GenerateCompressFiles(shareId string, shareFileData []sharemodel.ShareFileData, root *os.Root, target string) (string, error) {
	compressFileName := fmt.Sprintf("%s%d", shareId, time.Now().Unix())
	switch target {
	case "zip":
		return createZipFile(compressFileName, shareFileData, root)
	case "tar.gz":
		return createTarCompressedFile(compressFileName, "tar.gz", shareFileData, root, func(w io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriter(w), nil
		})
	case "tar.zst":
		return createTarCompressedFile(compressFileName, "tar.zst", shareFileData, root, func(w io.Writer) (io.WriteCloser, error) {
			return zstd.NewWriter(w)
		})
	case "tar.s2":
		return createTarCompressedFile(compressFileName, "tar.s2", shareFileData, root, func(w io.Writer) (io.WriteCloser, error) {
			return s2.NewWriter(w), nil
		})
	case "tar.snappy":
		return createTarCompressedFile(compressFileName, "tar.snappy", shareFileData, root, func(w io.Writer) (io.WriteCloser, error) {
			return snappy.NewBufferedWriter(w), nil
		})
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedCompressType, target)
	}
}

func createZipFile(compressFileName string, shareFileData []sharemodel.ShareFileData, root *os.Root) (string, error) {
	compressFullName := fmt.Sprintf("%s.zip", compressFileName)
	out, err := root.Create(compressFullName)
	if err != nil {
		return "", err
	}
	defer out.Close() //nolint:errcheck
	zw := zip.NewWriter(out)
	defer zw.Close() //nolint:errcheck
	for _, file := range shareFileData {
		f, err := root.Open(file.Id)
		if err != nil {
			return "", err
		}

		info, err := f.Stat()
		if err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}
		// zip 内部的文件名
		header.Name = file.FileName
		// 使用 Deflate 压缩
		header.Method = zip.Deflate

		w, err := zw.CreateHeader(header)
		if err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}

		_, err = io.Copy(w, f)
		if closeErr := f.Close(); closeErr != nil {
			return "", closeErr
		}
		if err != nil {
			return "", err
		}
	}
	return compressFullName, nil
}

type compressWriterFactory func(io.Writer) (io.WriteCloser, error)

func createTarCompressedFile(
	compressFileName string,
	extension string,
	shareFileData []sharemodel.ShareFileData,
	root *os.Root,
	newCompressWriter compressWriterFactory,
) (string, error) {
	compressFullName := fmt.Sprintf("%s.%s", compressFileName, extension)
	out, err := root.Create(compressFullName)
	if err != nil {
		return "", err
	}
	defer out.Close() //nolint:errcheck
	compressWriter, err := newCompressWriter(out)
	if err != nil {
		return "", err
	}
	defer compressWriter.Close() //nolint:errcheck

	tw := tar.NewWriter(compressWriter)
	defer tw.Close() //nolint:errcheck

	for _, file := range shareFileData {
		f, err := root.Open(file.Id)
		if err != nil {
			return "", err
		}

		info, err := f.Stat()
		if err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}
		header.Name = file.FileName
		if err := tw.WriteHeader(header); err != nil {
			f.Close() //nolint:errcheck
			return "", err
		}
		_, err = io.Copy(tw, f)
		if closeErr := f.Close(); closeErr != nil {
			return "", closeErr
		}
		if err != nil {
			return "", err
		}
	}
	return compressFullName, nil
}
