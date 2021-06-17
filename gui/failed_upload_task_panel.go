package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"time"
	"uploader/db"
	"uploader/util"
)

type FailedUploadTaskPanel struct {
	GroupBox						*walk.GroupBox
	RecoverButton					*walk.PushButton
	DeleteButton					*walk.PushButton
}

func InitFailedTasksPanels(failedTaskList []*db.UploadTaskRecord) {
	for i:= 0; i < len(failedTaskList); i++ {
		_, _ = AddFailedUploadTaskPanel(failedTaskList[i])
	}
}

func RemoveFailedTaskPanel(taskId string) {
	children := mmw.FailedTaskScrollView.Children()
	for i := 0; i < children.Len(); i++ {
		if taskId == children.At(i).Name() {
			widget := children.At(i)
			_ = widget.SetParent(nil)
			widget.Dispose()
			break
		}
	}
	CheckScrollViewEmpty(mmw.FailedTaskScrollView, mmw.FailedTaskScrollEmpty)
}


// Async Function
func AddFailedUploadTaskPanel(uploadTask *db.UploadTaskRecord) (uploadTaskPanel *FailedUploadTaskPanel, err error) {
	container := mmw.FailedTaskScrollView
	uploadTaskPanel = &FailedUploadTaskPanel{}

	err = GroupBox{
		Name:       uploadTask.TaskId,
		AssignTo:   &uploadTaskPanel.GroupBox,
		Background: SystemColorBrush{Color: walk.SysColorWindow},
		Layout:     Grid{Columns: 2},
		Visible:    true,
		Children: []Widget{
			ImageView{
				Image:   GetFileCategoryImage(uploadTask.LocalPath, uploadTask.IsDir),
				MaxSize: Size{Width: 70, Height: 70},
				MinSize: Size{Width: 70, Height: 70},
				Mode:    ImageViewModeShrink,
			},
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					Composite{
						Layout: Grid{Columns: 4, Margins: Margins{Left: 5, Right: 0, Top: 0, Bottom: 0}},
						Children: []Widget{
							TextLabel{
								Text: fmt.Sprintf(`任务失败: %v`, uploadTask.ErrorMessage),
								Font: Font{
									Bold:      false,
									PointSize: 15,
									Family:    "Microsoft YaHei",
								},
							},
							TextLabel{
								Text: util.MakeSpace(12),
							},
							PushButton{
								AssignTo:       &uploadTaskPanel.RecoverButton,
								ImageAboveText: true,
								Text:           "恢复",
								Image:          ImageResourcePath("recover.png"),
								MinSize:        Size{Height: 25, Width: 25},
								MaxSize:        Size{Height: 25, Width: 25},
								OnClicked: func() {
									//TODO::
									walk.MsgBox(mmw, "提示", "暂不支持重传", walk.MsgBoxIconInformation)
								},
							},
							PushButton{
								AssignTo:       &uploadTaskPanel.DeleteButton,
								ImageAboveText: true,
								Text:           "删除",
								Image:          ImageResourcePath("delete.png"),
								MinSize:        Size{Height: 25, Width: 25},
								MaxSize:        Size{Height: 25, Width: 25},
								OnClicked: func() {
									RemoveFailedTaskPanel(uploadTask.TaskId)
									_ = db.DeleteTaskRecord(uploadTask.TaskId)
								},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							TextLabel{
								Text: fmt.Sprintf("本地路径: %v", util.StringOmit(uploadTask.LocalPath, 50)),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("远程路径: %v", util.StringOmit(uploadTask.TargetPath, 50)),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("开始时间: %v", time.Unix(uploadTask.StartTime, 0).Format("2006-01-02 15:04:05")),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
						},
					},
				},
			},
		},
	}.Create(NewBuilder(container))


	if err != nil {
		log.Println(err)
	} else {
		mmw.FailedTaskScrollEmpty.SetVisible(false)
	}

	return uploadTaskPanel, err
}