package task

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"
	"uploader/db"
	"uploader/gui"
	"uploader/util"
)

func FileStatisticGUIRoutine(lock *sync.WaitGroup,  statisticTask *FileStatisticTask, panel *gui.ActiveUploadTaskPanel, uploadTask *UploadTask) {
	GUIRefreshSlot := time.Second * 1
	defer lock.Done()
	log.Println("FileStatisticGUIRoutine Start")

	for {

		select {
		case _ = <- uploadTask.RuntimeChannel.StatisticGUIRoutineShouldExit:
			log.Println("FileStatisticGUIRoutine Receive Exit Signal, Ready To Return")
			return
		default:
		}

		// check action
		//select {
		//case _ = <- uploadTask.RuntimeChannel.AbortChan:
		//	log.Println("FileStatisticGUIRoutine Receive Abort Signal, Ready To Return")
		//	return
		//default:
		//}

		panel.GroupBox.Synchronize(func() {
			_ = panel.StatusTextLabel.SetText(fmt.Sprintf("已发现 %v 个项目(%v)", statisticTask.ItemsFound, util.FileSizeFromBytes(statisticTask.BytesCount)))
			_ = panel.SrcAndTargetLinkLabel.SetText(fmt.Sprintf(`正在将 %v 个项目从 <a href="%v">%v</a> 复制到 <a href="%v">%v</a>`,
				statisticTask.ItemsFound,  util.DirPath(uploadTask.LocalPath), util.DirName(uploadTask.LocalPath),
				util.DirPath(uploadTask.TargetPath), util.DirName(uploadTask.TargetPath)))
		})

		time.Sleep(GUIRefreshSlot)
	}

}


func FileCopyGUIRoutine(lock *sync.WaitGroup, panel *gui.ActiveUploadTaskPanel, uploadTask *UploadTask) {
	lock.Wait()
	log.Println("enter FileCopyGUIRoutine")
	secondSlot := int64(1)
	GUIRefreshSlot := time.Second * time.Duration(secondSlot)
	lastCopiedSize := uploadTask.BytesCopied
	cnt := 0

	defer panel.GroupBox.Synchronize(func() {
		_ = panel.StatusTextLabel.SetText(fmt.Sprintf("正在上传 - %v%%", 100))
		panel.ProgressBar.SetValue(100)
		_ = panel.ItemLeftTextLabel.SetText(fmt.Sprintf("剩余项目: %v(%v)", 0, util.FileSizeFromBytes(0)))
	})

	for {

		select {
		case _ = <- uploadTask.RuntimeChannel.FileCopyGUIRoutineShouldExit:
			return
		default:
		}

		// check action
		select {
		case _ = <- uploadTask.RuntimeChannel.AbortChan:
			log.Println("FileCopyGUIRoutine Receive Abort Signal, Ready To Return")
			return
		case _ = <- uploadTask.RuntimeChannel.SuspendChan:
			log.Println("suspend, block FileCopyGUIRoutine")
			select {
			case _ = <- uploadTask.RuntimeChannel.ResumeChan:
			case _ = <- uploadTask.RuntimeChannel.AbortChan:
				log.Println("FileCopyGUIRoutine Receive Abort Signal, Ready To Return")
				return
			}
			log.Println("resume, recover FileCopyGUIRoutine")
		default:
		}


		itemsLeft := uploadTask.ItemsTotal - uploadTask.ItemsCopied
		bytesLeft := uploadTask.BytesTotal - uploadTask.BytesCopied
		percent := uploadTask.BytesCopied * 100 / uploadTask.BytesTotal
		bytesSpeed := (uploadTask.BytesCopied - lastCopiedSize) / secondSlot
		if bytesSpeed == 0 {
			bytesSpeed = 1
		}
		speed := util.FileSizeFromBytes(bytesSpeed)
		secLeft := bytesLeft / bytesSpeed
		lastCopiedSize = uploadTask.BytesCopied
		panel.GroupBox.Synchronize(func() {
			_ = panel.StatusTextLabel.SetText(fmt.Sprintf("正在上传 - %v%%", percent))
			panel.ProgressBar.SetValue(int(percent))
			_ = panel.SpeedTextLabel.SetText(fmt.Sprintf("当前速度: %v/s", speed))
			_ = panel.ItemLeftTextLabel.SetText(fmt.Sprintf("剩余项目: %v(%v)", itemsLeft, util.FileSizeFromBytes(bytesLeft)))
			_ = panel.CurrentCopyNameTextLabel.SetText(fmt.Sprintf("正在上传: %v", filepath.Base(uploadTask.CurrentCopyItemPath)))
			if cnt % 3 == 0 {
				_ = panel.TimeLeftTextLabel.SetText(fmt.Sprintf("剩余时间: %v", util.SecondToTime(secLeft)))
			}
		})

		cnt++
		time.Sleep(GUIRefreshSlot)
	}
}


func PostMissionExitRoutine(lock *sync.WaitGroup, uploadTask *UploadTask) {
	lock.Wait()

	log.Println("enter PostMissionExitRoutine, ",uploadTask)

	// judge current mission status
	// mission may be:
	// 1. Canceled
	// 2. Succeed
	// 3. Failed
	if uploadTask.Status != UploadStatusFailed {
		uploadTask.Status = UploadStatusSucceed
	}

	uploadTask.FinishTime = time.Now().Unix()

	// sync to db
	uploadTaskRecord := UploadTaskToUploadTaskRecord(uploadTask)
	_ = db.UpdateTaskRecord(uploadTaskRecord)

	// remove for gui active list, and add for inactive list
	gui.GetMyMainWindow().Synchronize(func() {
		gui.RemoveActiveTaskPanel(uploadTask.TaskId)
		if uploadTask.Status == UploadStatusSucceed {
			_, _ = gui.AddSucceedUploadTaskPanel(&uploadTaskRecord)
		} else {
			_, _ = gui.AddFailedUploadTaskPanel(&uploadTaskRecord)
		}
	})

	//TODO:: Add and fix bugs
	//close(uploadTask.RuntimeChannel.FileCopyGUIRoutineShouldExit)
	//close(uploadTask.RuntimeChannel.StatisticGUIRoutineShouldExit)
	//close(uploadTask.RuntimeChannel.AbortChan)
	//close(uploadTask.RuntimeChannel.ResumeChan)
	//close(uploadTask.RuntimeChannel.SuspendChan)
}