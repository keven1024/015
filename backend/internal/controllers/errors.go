package controllers

import "errors"

var (
	// 通用错误（参数校验失败）
	ErrInvalidRequest = errors.New("InvalidRequest") // 调用接口参数错误

	// 任务相关
	ErrTaskNotFound = errors.New("TaskNotFound") // 任务不存在
	ErrTaskExpired  = errors.New("TaskExpired")  // 任务已过期

	// 文件上传相关
	ErrInsufficientStorage    = errors.New("InsufficientStorage")    // 存储空间不足
	ErrUploadTaskExpired      = errors.New("UploadTaskExpired")      // 上传任务已过期
	ErrInvalidUploadTaskState = errors.New("InvalidUploadTaskState") // 上传任务状态错误
	ErrInvalidFileSliceIndex  = errors.New("InvalidFileSliceIndex")  // 文件切片索引错误
	ErrInvalidFileSliceSize   = errors.New("InvalidFileSliceSize")   // 文件切片大小错误
	ErrIncompleteFileSlices   = errors.New("IncompleteFileSlices")   // 文件切片不完整
	ErrFileHashMismatch       = errors.New("FileHashMismatch")       // 文件Hash不一致

	// 分享相关
	ErrShareFileNotFound     = errors.New("ShareFileNotFound")     // 分享文件不存在
	ErrInvalidShareFileState = errors.New("InvalidShareFileState") // 分享文件状态错误
	ErrShareNotFound         = errors.New("ShareNotFound")         // 分享不存在

	// 下载相关
	ErrInvalidSharePassword      = errors.New("InvalidSharePassword")      // 分享密码错误
	ErrInsufficientDownloadQuota = errors.New("InsufficientDownloadQuota") // 下载次数不足
)
