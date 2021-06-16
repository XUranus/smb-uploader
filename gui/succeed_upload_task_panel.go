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

// Upload Task Panel Control Block
type SucceedUploadTaskPanel struct {
	GroupBox						*walk.GroupBox
	DeleteButton					*walk.PushButton
}

func InitSucceedTasksPanels(succeedTaskList []*db.UploadTaskRecord) {
	for i:= 0; i < len(succeedTaskList); i++ {
		_, _ = AddSucceedUploadTaskPanel(succeedTaskList[i])
	}
}


func RemoveSucceedTaskPanel(taskId string) {
	children := mmw.SucceedTaskScrollView.Children()
	for i := 0; i < children.Len(); i++ {
		if taskId == children.At(i).Name() {
			widget := children.At(i)
			_ = widget.SetParent(nil)
			widget.Dispose()
			fmt.Println(children.Len())
			break
		}
	}

	CheckScrollViewEmpty(mmw.SucceedTaskScrollView, mmw.SucceedTaskScrollEmpty)
}


// Async Function
func AddSucceedUploadTaskPanel(uploadTask *db.UploadTaskRecord) (uploadTaskPanel *SucceedUploadTaskPanel, err error) {
	container := mmw.SucceedTaskScrollView
	uploadTaskPanel = &SucceedUploadTaskPanel{}

	err = GroupBox{
		Name: uploadTask.TaskId,
		AssignTo: &uploadTaskPanel.GroupBox,
		Background: SystemColorBrush{Color: walk.SysColorWindow},
		Layout:     Grid{Columns: 2},
		Visible: true,
		Children: []Widget{
			ImageView{
				Image: GetFileCategoryImage(uploadTask.LocalPath, uploadTask.IsDir),
				MaxSize: Size{Width: 70,Height: 70},
				MinSize: Size{Width: 70,Height: 70},
				Mode: ImageViewModeShrink,
			},
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					Composite{
						Layout: Grid{Columns: 3, Margins: Margins{Left: 5, Right: 0, Top: 0, Bottom: 0}},
						Children: []Widget{
							TextLabel{
								Text: fmt.Sprintf(`已完成`),
								Font: Font{
									Bold:      false,
									PointSize: 15,
									Family:    "Microsoft YaHei",
								},
							},
							TextLabel{
								Text: "            ",
							},
							PushButton{
								AssignTo: &uploadTaskPanel.DeleteButton,
								ImageAboveText: true,
								Text: "删除",
								Image: ImageResourcePath("delete.png"),
								MinSize: Size{Height: 25, Width: 25},
								MaxSize: Size{Height: 25, Width: 25},
								OnClicked: func() {
									container.Synchronize(func() {
										RemoveSucceedTaskPanel(uploadTask.TaskId)
										_ = db.DeleteTaskRecord(uploadTask.TaskId)
									})
								},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							TextLabel{
								Text: fmt.Sprintf("本地路径: %v", util.StringOmit(uploadTask.LocalPath, 50)),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("远程路径: %v", util.StringOmit(uploadTask.TargetPath, 50)),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 2},
						Children: []Widget{
							TextLabel{
								Text: fmt.Sprintf("总大小: %v", util.FileSizeFromBytes(uploadTask.BytesTotal)),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("总项目: %v",uploadTask.ItemsTotal),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("开始时间: %v",time.Unix(uploadTask.StartTime, 0).Format("2006-01-02 15:04:05")),
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							},
							TextLabel{
								Text: fmt.Sprintf("结束时间: %v",time.Unix(uploadTask.FinishTime, 0).Format("2006-01-02 15:04:05")),
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
		mmw.SucceedTaskScrollEmpty.SetVisible(false)
	}

	return uploadTaskPanel, err
}