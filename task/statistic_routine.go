package task

import (
	"log"
	"os"
	"path/filepath"
)

/**
	FileStatisticTask statistic total items and size
 */
type FileStatisticTask struct {
	Signal 						*RoutineSignal

	// statistic fields
	SourcePath					string
	IsDir						bool

	// dynamic fields
	ItemsFound					int64
	BytesCount					int64

	// events
	OnItemsFoundChanged			func(int64)
	OnBytesCountChanged			func(int64)
	OnExit						func(err error, size int64, item int64)
}


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
	var err error = nil
	if statisticTask.IsDir {
		// size and items update is included in statisticTask.CalDirSize()
		_ , err = statisticTask.CalDirSize(statisticTask.SourcePath)

	} else {
		// single file size calculate can be done in O(1)
		statisticTask.ItemsFound = 1
		statisticTask.OnItemsFoundChanged(1)

		var size int64 = 0
		size, err = calFileSize(statisticTask.SourcePath)
		statisticTask.BytesCount = size
		statisticTask.OnBytesCountChanged(size)
	}

	statisticTask.OnExit(err, statisticTask.BytesCount, statisticTask.ItemsFound)
}


func (statisticTask *FileStatisticTask) CalDirSize(path string) (int64, error) {
	var size int64
	var count int64 = 0
	signal := statisticTask.Signal

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {

		if abort := signal.CheckSignal(); abort {
			log.Println("FileStatisticTask received exit signal, return")
			return AbortError
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