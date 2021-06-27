package task

import (
	"fmt"
	"time"
	"uploader/gui"
	"uploader/logger"
	"uploader/util"
)



/**
	FileStatisticGUIRoutine monitor FileStatisticTask and synchronize statistic state to GUI
 */
type FileStatisticGUIRoutine struct {
	Signal 				*RoutineSignal

	// fields
	StatisticTask 		*FileStatisticTask
	LocalPath			string
	TargetPath			string
	Panel 				*gui.ActiveUploadTaskPanel

	// events
	OnExit				func(error)
}

func (routine *FileStatisticGUIRoutine) Start(async bool) {
	if async {
		go func() {
			routine.StartBlock()
		}()
	} else {
		routine.StartBlock()
	}
}

func (routine *FileStatisticGUIRoutine) StartBlock() {
	logger.CommonLogger.Info("StartBlock", "FileStatisticGUIRoutine Start")
	GUIRefreshSlot := time.Second * 1
	var err error = nil

	defer routine.OnExit(err)

	for {

		if abort := routine.Signal.CheckSignal(); abort {
			logger.CommonLogger.Info("StartBlock", "FileStatisticGUIRoutine received exit signal, return")
			err = AbortError
			return
		}

		routine.Panel.GroupBox.Synchronize(func() {
			_ = routine.Panel.StatusTextLabel.SetText(fmt.Sprintf("已发现 %v 个项目(%v)", util.NumberWithComma(routine.StatisticTask.ItemsFound), util.FileSizeFromBytes(routine.StatisticTask.BytesCount)))
			_ = routine.Panel.SrcAndTargetLinkLabel.SetText(fmt.Sprintf(`正在将 %v 个项目从 <a href="%v">%v</a> 复制到 <a href="%v">%v</a>`,
				util.NumberWithComma(routine.StatisticTask.ItemsFound),  util.DirPath(routine.LocalPath), util.DirName(routine.LocalPath),
				util.DirPath(routine.TargetPath), util.DirName(routine.TargetPath)))
		})

		time.Sleep(GUIRefreshSlot)
	}

}

