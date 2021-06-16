package task

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FileStatisticTask struct {
	// statistic
	SourcePath					string
	IsDir						bool
	UploadTaskRef				*UploadTask

	// dynamic
	ItemsFound					int64
	BytesCount					int64

	OnItemsFoundChanged		func(int64)
	OnBytesCountChanged		func(int64)
	OnCompleted				func(int64, int64)
	OnFailed				func(err error)
}

/*
	statistic routine only support abort action
 */

func (statisticTask *FileStatisticTask) Start(async bool) {
	if async {
		go func() {
			statisticTask.BlockStart()
		}()
	} else {
		statisticTask.BlockStart()
	}
}


func (statisticTask *FileStatisticTask)	BlockStart()  {
	if statisticTask.IsDir {
		_ , err := statisticTask.CalDirSize(statisticTask.SourcePath)
		if err != nil {
			statisticTask.OnFailed(err)
			return
		}
	} else {
		statisticTask.ItemsFound = 1
		statisticTask.OnItemsFoundChanged(1)

		size, err := calFileSize(statisticTask.SourcePath)
		if err != nil {
			statisticTask.OnFailed(err)
			return
		}
		statisticTask.BytesCount = size
		statisticTask.OnBytesCountChanged(size)
	}

	statisticTask.UploadTaskRef.BytesTotal = statisticTask.BytesCount
	statisticTask.UploadTaskRef.ItemsTotal = statisticTask.ItemsFound
	statisticTask.UploadTaskRef.BytesCalculated = true
	statisticTask.UploadTaskRef.ItemsCalculated = true

	statisticTask.OnCompleted(statisticTask.ItemsFound, statisticTask.BytesCount)
}


func (statisticTask *FileStatisticTask) CalDirSize(path string) (int64, error) {
	var size int64
	var count int64 = 0
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {

		// check action
		select {
		case _ = <-statisticTask.UploadTaskRef.RuntimeChannel.AbortChan:
			log.Println("FileStatisticRoutine Receive Abort Signal, Ready To Return")
			return errors.New("abort signal found")
		default:
		}

		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
			count += 1
			statisticTask.BytesCount = size
			statisticTask.ItemsFound = count
			statisticTask.OnBytesCountChanged(size)
			statisticTask.OnItemsFoundChanged(count)
		}
		return err
	})
	return size, err
}


func calFileSize(path string) (int64, error) {
	fi,err:=os.Stat(path)
	if err == nil {
		return fi.Size(), err
	} else {
		return 0, err
	}
}