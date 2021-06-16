package db

/*
	static records of upload tasks
 */

type UploadTaskRecordStatus	string

const (
	Active		UploadTaskRecordStatus = "ACTIVE"
	Succeed     UploadTaskRecordStatus = "SUCCESS"
	Failed      UploadTaskRecordStatus = "FAILED"
)


type UploadTaskRecord struct {
	// Unique Identifier
	TaskId				string

	Status 				UploadTaskRecordStatus

	// Path
	LocalPath 			string // path of source folder or File
	TargetPath 			string // path of target folder
	IsDir				bool

	BytesCopied			int64 // n bytes
	BytesCalculated		bool
	BytesTotal			int64 // n bytes

	// when IsDir == false, ItemsTotal = 1, ItemsCopied = 0 or 1
	ItemsTotal			int64 // total number of items, usually used in folder copy
	ItemsCalculated		bool
	ItemsCopied			int64 // number of items copied, usually used in folder copy

	StartTime   		int64 // sec timestamp of task start
	FinishTime			int64 // sec timestamp of task finished

	ErrorMessage		string // nil unless UploadTask is failed
}

