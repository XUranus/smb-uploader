package task

import "uploader/db"

func UploadTaskToUploadTaskRecord(uploadTask *UploadTask) db.UploadTaskRecord {
	status := db.Active
	if uploadTask.Status == UploadStatusFailed {
		status = db.Failed
	} else if uploadTask.Status == UploadStatusSucceed {
		status = db.Succeed
	}

	errMsg := ""
	if uploadTask.Error != nil {
		errMsg = uploadTask.Error.Error()
	}

	return db.UploadTaskRecord{
		TaskId: uploadTask.TaskId,
		Status: status,
		LocalPath: uploadTask.LocalPath,
		TargetPath: uploadTask.TargetPath,
		IsDir: uploadTask.IsDir,

		BytesCopied: uploadTask.BytesCopied,
		BytesCalculated: uploadTask.BytesCalculated,
		BytesTotal: uploadTask.BytesTotal,

		ItemsTotal: uploadTask.ItemsTotal,
		ItemsCalculated: uploadTask.ItemsCalculated,
		ItemsCopied: uploadTask.ItemsCopied,

		StartTime: uploadTask.StartTime,
		FinishTime: uploadTask.FinishTime,
		ErrorMessage: errMsg,
	}
}