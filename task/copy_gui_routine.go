package task

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
	"uploader/util"
)

/**
	FileCopyGUIRoutine monitor FileCopyTask and synchronize copy state to GUI
*/

type FileCopyGUIRoutine struct {
	Signal 				*RoutineSignal

	UploadTaskRef 		*UploadTask

	// events
	OnStart				func()
	OnExit				func(error)
}


func (routine *FileCopyGUIRoutine) Start(async bool) {
	if async {
		go func() {
			routine.StartBlock()
		}()
	} else {
		routine.StartBlock()
	}
}


func (routine *FileCopyGUIRoutine) StartBlock() {
	routine.OnStart()
	log.Println("FileCopyGUIRoutine start")

	var err error = nil
	secondSlot := int64(1)
	GUIRefreshSlot := time.Second * time.Duration(secondSlot)
	uploadTask := routine.UploadTaskRef
	panel := uploadTask.GUIPanel
	loopCounter := 0
	secPassed := time.Now().Unix() - uploadTask.StartTime + 1

	defer routine.OnExit(err)

	_ = panel.ProgressBar.SetMarqueeMode(false)
	panel.ContinueOrSuspendButton.SetVisible(true)
	_ = panel.SrcAndTargetLinkLabel.SetText(fmt.Sprintf(`正在将 %v 个项目从 <a href="%v">%v</a> 复制到 <a href="%v">%v</a>`,
		routine.UploadTaskRef.ItemsTotal,  util.DirPath(routine.UploadTaskRef.LocalPath), util.DirName(routine.UploadTaskRef.LocalPath),
		util.DirPath(routine.UploadTaskRef.TargetPath), util.DirName(routine.UploadTaskRef.TargetPath)))

	for {

		if abort := routine.Signal.CheckSignal(); abort {
			log.Println("FileCopyGUIRoutine received exit signal, return")
			err = AbortError
			return
		}

		itemsLeft := uploadTask.ItemsTotal - uploadTask.ItemsCopied
		bytesLeft := uploadTask.BytesTotal - uploadTask.BytesCopied

		percent := 0
		if uploadTask.BytesTotal > 0 {
			percent = int(uploadTask.BytesCopied * 100 / uploadTask.BytesTotal)
		} else {
			percent = 100
		}

		fmt.Println(int64(GUIRefreshSlot.Seconds()))
		secPassed += int64(GUIRefreshSlot.Seconds())
		bytesSpeed := uploadTask.BytesCopied / secPassed
		speed := util.FileSizeFromBytes(bytesSpeed)

		var secLeft int64
		if bytesSpeed > 0 {
			secLeft = bytesLeft / bytesSpeed
		} else {
			secLeft = 365 * 24 * 60 * 60 // 1 year
		}

		panel.GroupBox.Synchronize(func() {
			_ = panel.StatusTextLabel.SetText(fmt.Sprintf("正在上传 - %v%%", percent))
			panel.ProgressBar.SetValue(percent)
			_ = panel.SpeedTextLabel.SetText(fmt.Sprintf("当前速度: %v/s", speed))
			_ = panel.ItemLeftTextLabel.SetText(fmt.Sprintf("剩余项目: %v(%v)", itemsLeft, util.FileSizeFromBytes(bytesLeft)))
			_ = panel.CurrentCopyNameTextLabel.SetText(fmt.Sprintf("正在上传: %v", util.StringOmit(filepath.Base(uploadTask.CurrentCopyItemPath),20)))
			_ = panel.TimeLeftTextLabel.SetText(fmt.Sprintf("剩余时间: %v", util.SecondToTime(secLeft)))
		})

		loopCounter ++
		time.Sleep(GUIRefreshSlot)
	}
}

