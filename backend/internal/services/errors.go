package services

import "errors"

var (
	ErrCreateMergeFile         = errors.New("failed to create merged file")
	ErrMergeFileSlices         = errors.New("failed to merge file slices")
	ErrUnsupportedCompressType = errors.New("unsupported compress type")
)
