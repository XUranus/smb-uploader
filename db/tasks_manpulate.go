package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func LoadInactiveUploadTasksFromDB(status UploadTaskRecordStatus) ([]*UploadTaskRecord, error) {
	db, err := sql.Open("sqlite3", "C:\\Users\\A\\Desktop\\uploader\\bin\\data.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		"SELECT " +
		"id, local_path, target_path, is_dir, bytes_copied, bytes_total_calculated, bytes_total, items_copied, items_total_calculated, items_total, start_time, finish_time, error_msg" +
		" FROM upload " +
		" WHERE status = ? and deleted = 0" +
			" ORDER BY start_time asc", status)
	if err != nil {
		return nil, err
	}

	uploadTaskList := make([]*UploadTaskRecord, 0, 10)

	for rows.Next() {
		var	taskId		string
		var localPath	string
		var targetPath	string
		var isDir		int
		var bytesCopied	int64
		var bytesCalculated int64
		var bytesTotal	int64
		var itemsCopied int64
		var itemsCalculated int64
		var itemsTotal 	int64
		var startTime	int64
		var finishTime 	int64
		var errMsg		string

		err = rows.Scan(&taskId, &localPath, &targetPath, &isDir,
			&bytesCopied, &bytesCalculated, &bytesTotal,
			&itemsCopied, &itemsCalculated, &itemsTotal,
			&startTime, &finishTime, &errMsg)
		if err != nil {
			return nil, err
		}

		taskItem := &UploadTaskRecord{
			TaskId: taskId,
			Status: status,
			LocalPath: localPath,
			TargetPath: targetPath,
			IsDir: isDir == 1,
			BytesCopied: bytesCopied,
			BytesCalculated: bytesCalculated == 1,
			BytesTotal: bytesTotal,
			ItemsCopied: itemsCopied,
			ItemsCalculated: itemsCalculated == 1,
			ItemsTotal: itemsTotal,
			StartTime: startTime,
			FinishTime: finishTime,
			ErrorMessage: errMsg,
		}

		uploadTaskList = append(uploadTaskList, taskItem)
	}

	return uploadTaskList, nil
}
