package task

import (
	"fmt"
	"log"
	"sync"
	"time"
	"uploader/gui"
)

type UploadStatus int32

const (
	// Active State
	UploadStatusPending 	UploadStatus = 0
	UploadStatusRunning     UploadStatus = 1
	UploadStatusSuspend     UploadStatus = 2
	UploadStatusRecovering  UploadStatus = 3

	// Inactive State
	UploadStatusSucceed     UploadStatus = 4
	UploadStatusFailed      UploadStatus = 5
)



type UploadTask struct {
	// Unique Identifier
	TaskId				string

	Status 				UploadStatus

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
	TimeLeft			int64 // sec timestamp of estimated time left
	CurrentCopyItemPath	string // absolute path of file copying

	Error				error // nil unless UploadTask is failed

	// RuntimeFlag
	GUIPanel 						*gui.ActiveUploadTaskPanel
	FileCopyGUIRoutineSignal		*RoutineSignal
	FileStatisticGUIRoutineSignal	*RoutineSignal
	FileCopyTaskSignal				*RoutineSignal
	FileStatisticTaskSignal			*RoutineSignal
}


func NewUploadTask(taskId string, localPath string, targetPath string, isDir bool) (uploadTask *UploadTask) {
	uploadTask = &UploadTask{
		TaskId:      taskId,
		Status:      UploadStatusPending,
		LocalPath:   localPath,
		TargetPath:  targetPath,
		IsDir:       isDir,
		BytesCopied: 0,
		BytesCalculated: false,
		BytesTotal:  0,
		ItemsTotal:  0,
		ItemsCalculated: false,
		ItemsCopied: 0,
		StartTime:   time.Now().Unix(),
		FinishTime:  0,
		TimeLeft:    0,
		CurrentCopyItemPath: "",

		Error: nil,
	}
	//register in global map
	taskMap[taskId] = uploadTask
	return
}


// start schedule new task
func (uploadTask *UploadTask) Start() {
	// sync init state to gui
	var panel *gui.ActiveUploadTaskPanel
	var panelCreateLock sync.WaitGroup
	panelCreateLock.Add(1)
	gui.GetMyMainWindow().ActiveTaskScrollView.Synchronize(func() {
		panel, _ = gui.AddActiveUploadTaskPanel(uploadTask.TaskId, uploadTask.LocalPath, uploadTask.TargetPath, uploadTask.IsDir)
		panelCreateLock.Done()
	})
	panelCreateLock.Wait()


	uploadTask.GUIPanel = panel

	// when all done, activate PostMissionRoutine to make task inactive
	activeLock := &sync.WaitGroup{}
	activeLock.Add(4)

	// make FileCopyGUIRoutine start after FileStatisticGUIRoutine
	guiSeqLock := &sync.WaitGroup{}
	guiSeqLock.Add(1)

	// init two gui routine control block
	fileStatisticGUIRoutineSignal := NewRoutineSignal()
	fileCopyGUIRoutineSignal := NewRoutineSignal()
	statisticTaskSignal := NewRoutineSignal()
	copyTaskSignal := NewRoutineSignal()


	// 1. statistic routine
	statisticTask := &FileStatisticTask{
		Signal: statisticTaskSignal,

		SourcePath: uploadTask.LocalPath,
		IsDir: uploadTask.IsDir,

		ItemsFound: 0,
		BytesCount: 0,

		OnBytesCountChanged: func(n int64) {},
		OnItemsFoundChanged: func(n int64) {},
		OnExit: func(err error, size int64, item int64) {
			if err != nil {
				uploadTask.Error = err
				uploadTask.Status = UploadStatusFailed
			} else {
				uploadTask.Status = UploadStatusRunning
				uploadTask.BytesTotal = size
				uploadTask.ItemsTotal = item
			}

			fileStatisticGUIRoutineSignal.SendAbortSignal()
			log.Println("statisticTask Lock Release")
			activeLock.Done()
		},
	}


	// 2. copy routine
	fileCopyTask := &FileCopyTask{
		Signal: copyTaskSignal,

		SourcePath: uploadTask.LocalPath,
		TargetPath: uploadTask.TargetPath,
		IsDir: uploadTask.IsDir,
		UploadTaskRef: uploadTask,

		ItemsCopied: 0,
		BytesCopied: 0,
		CurrentCopyItemPath: "",

		OnItemsCopiedChanged: func(n int64) {
			uploadTask.ItemsCopied = n
		},
		OnBytesCopiedChanged: func(n int64) {
			uploadTask.BytesCopied = n
		},
		OnCurrentCopyItemChanged: func(p string) {
			uploadTask.CurrentCopyItemPath = p
		},
		OnExit: func(err error) {
			uploadTask.Error = err
			if err != nil {
				uploadTask.Status = UploadStatusFailed
				uploadTask.Error = err
			}

			fileCopyGUIRoutineSignal.SendAbortSignal()
			log.Println("fileCopyTask Lock Release")
			activeLock.Done()
		},
	}

	// 3. statistic status gui async routine (implemented by timer, this is a compromise for the sake of compromise)
	fileStatisticGUIRoutine := FileStatisticGUIRoutine{
		Signal: fileStatisticGUIRoutineSignal,

		StatisticTask: statisticTask,
		LocalPath: uploadTask.LocalPath,
		TargetPath: uploadTask.TargetPath,
		Panel: panel,
		OnExit: func(err error) {
			guiSeqLock.Done()
			log.Println("fileStatisticGUIRoutine Lock Release")
			activeLock.Done()
		},
	}


	// 4. copy status gui async routine (implemented by timer, this is a compromise for the sake of performance)
	fileCopyGUIRoutine := &FileCopyGUIRoutine{
		Signal: fileCopyGUIRoutineSignal,

		UploadTaskRef: uploadTask,
		OnStart: func() {
			guiSeqLock.Wait()
		},
		OnExit:	func(err error) {
			log.Println("fileCopyGUIRoutine Lock Release")
			activeLock.Done()
		},
	}



	uploadTask.FileCopyGUIRoutineSignal = fileCopyGUIRoutineSignal
	uploadTask.FileStatisticGUIRoutineSignal = fileStatisticGUIRoutineSignal
	uploadTask.FileCopyTaskSignal = copyTaskSignal
	uploadTask.FileStatisticTaskSignal = statisticTaskSignal


	statisticTask.Start(true)
	fileCopyTask.Start(true)
	fileCopyGUIRoutine.Start(true)
	fileStatisticGUIRoutine.Start(true)


	// 4 must run after 3 (3 -> 4)
	// 5. mission finished (executed after 4 routines are done)
	go PostMissionRoutine(activeLock, uploadTask)
}


func (uploadTask *UploadTask) Suspend() {
	if uploadTask.Status != UploadStatusRunning {
		log.Println("suspend requires running state, current ", uploadTask.Status)
		return
	} else {
		uploadTask.FileCopyTaskSignal.SendSuspendSignal()       // for copy routine
		uploadTask.FileCopyGUIRoutineSignal.SendSuspendSignal() // for copy gui routine

		uploadTask.Status = UploadStatusSuspend
		uploadTask.GUIPanel.GroupBox.Synchronize(func() {
			percent := uploadTask.BytesCopied * 100 / uploadTask.BytesTotal
			_ = uploadTask.GUIPanel.StatusTextLabel.SetText(fmt.Sprintf("已暂停 - %v%%", percent))
			uploadTask.GUIPanel.ContinueOrSuspendButton.SetEnabled(true)
			uploadTask.GUIPanel.CancelButton.SetEnabled(true)
		})
	}
}

func (uploadTask *UploadTask) Resume() {
	if uploadTask.Status != UploadStatusSuspend {
		log.Println("suspend requires suspend state, current ", uploadTask.Status)
		return
	} else {
		uploadTask.FileCopyTaskSignal.SendResumeSignal()       // for copy routine
		uploadTask.FileCopyGUIRoutineSignal.SendResumeSignal() // for copy gui routine

		uploadTask.Status = UploadStatusRunning
		uploadTask.GUIPanel.GroupBox.Synchronize(func() {
			uploadTask.GUIPanel.ContinueOrSuspendButton.SetEnabled(true)
			uploadTask.GUIPanel.CancelButton.SetEnabled(true)
		})
	}
}

func (uploadTask *UploadTask) Abort() {
	log.Println("current status ", uploadTask.Status)

	if uploadTask.Status == UploadStatusPending {
		uploadTask.FileCopyTaskSignal.SendAbortSignal() // for statistic routine
		uploadTask.FileStatisticTaskSignal.SendAbortSignal() // for statistic gui routine

	} else if uploadTask.Status == UploadStatusRunning {
		uploadTask.FileCopyTaskSignal.SendAbortSignal() // for copy routine

	} else if uploadTask.Status == UploadStatusSuspend {
		uploadTask.FileCopyTaskSignal.SendAbortSignal()
	}

	uploadTask.Status = UploadStatusFailed
}


