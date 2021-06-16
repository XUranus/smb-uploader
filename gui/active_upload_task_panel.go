package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os/exec"
	"syscall"
	"uploader/util"
)

// Upload Task Panel Control Block
type ActiveUploadTaskPanel struct {
	// Parent Container
	GroupBox						*walk.GroupBox

	// Sub Widgets
	SrcAndTargetLinkLabel 			*walk.LinkLabel
	StatusTextLabel 				*walk.TextLabel
	ContinueOrSuspendButton 		*walk.PushButton
	CancelButton					*walk.PushButton
	ProgressBar 					*walk.ProgressBar
	CurrentCopyNameTextLabel		*walk.TextLabel
	TimeLeftTextLabel				*walk.TextLabel
	ItemLeftTextLabel				*walk.TextLabel
	SpeedTextLabel					*walk.TextLabel
}

func RemoveActiveTaskPanel(taskId string) {
	children := mmw.ActiveTaskScrollView.Children()
	for i := 0; i < children.Len(); i++ {
		if taskId == children.At(i).Name() {
			widget := children.At(i)
			_ = widget.SetParent(nil)
			widget.Dispose()
			fmt.Println(children.Len())
			break
		}
	}
	CheckScrollViewEmpty(mmw.ActiveTaskScrollView, mmw.ActiveTaskScrollEmpty)
}


// Synchronize UploadTask to view, add a task view as new task view to panel
func AddActiveUploadTaskPanel(taskId string, localPath string, targetPath string, isDir bool) (uploadTaskPanel *ActiveUploadTaskPanel, err error) {

	container := mmw.ActiveTaskScrollView
	uploadTaskPanel = &ActiveUploadTaskPanel{}

	err = GroupBox{
		Name: taskId,
		AssignTo: &uploadTaskPanel.GroupBox,
		Background: SystemColorBrush{Color: walk.SysColorWindow},
		Layout:     Grid{Columns: 1},
		Visible: true,
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					Composite{
						Layout: Grid{Columns: 4, Margins: Margins{Left: 5, Right: 0, Top: 0, Bottom: 0}},
						Children: []Widget{
							LinkLabel{
								AssignTo: &uploadTaskPanel.SrcAndTargetLinkLabel,
								Font:     Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
								Text:     "初始化上传任务中",//`正在将 1 个项目从 <a href="">XXX</a> 复制到 <a href="">XXX</a>`,
								OnLinkActivated: func(link *walk.LinkLabelLink) {
									log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
									cmd := exec.Command("cmd", "/c", "start", link.URL())
									cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
									_ = cmd.Start()
								},
							},
							TextLabel{
								Text: util.MakeSpace(12),
							},
							PushButton{
								Name: "SendSuspendSignal",
								AssignTo: &uploadTaskPanel.ContinueOrSuspendButton,
								ImageAboveText: true,
								Text: "暂停/继续",
								Visible: false,
								Image:  ImageResourcePath("suspend.png"),
								MinSize: Size{Height: 25, Width: 25},
								MaxSize: Size{Height: 25, Width: 25},
								OnClicked: func() {
									button := uploadTaskPanel.ContinueOrSuspendButton
									statusTextLabel := uploadTaskPanel.StatusTextLabel

									if button.Name() == "SendSuspendSignal" { // switch to suspend
										_ = statusTextLabel.SetText("正在暂停...")
										button.SetName("Continue")
										continueImage, _ := walk.ImageFrom(ImageResourcePath("continue.png"))
										_ = button.SetImage(continueImage)
										RequestSuspendActiveTask(taskId)
									} else { // switch to continue
										_ = statusTextLabel.SetText("正在恢复...")
										button.SetName("SendSuspendSignal")
										suspendImage, _ := walk.ImageFrom(ImageResourcePath("suspend.png"))
										_ = button.SetImage(suspendImage)
										RequestResumeActiveTask(taskId)
									}
									button.SetEnabled(false)
									uploadTaskPanel.CancelButton.SetEnabled(false)
								},
							},
							PushButton{
								AssignTo: &uploadTaskPanel.CancelButton,
								ImageAboveText: true,
								Text: "取消",
								Image:  ImageResourcePath("cancel.png"),
								MinSize: Size{Height: 25, Width: 25},
								MaxSize: Size{Height: 25, Width: 25},
								OnClicked: func() {
									log.Println("cancel button pushed")
									RequestAbortActiveTask(taskId)
								},
							},
						},
					},
					TextLabel{
						AssignTo: &uploadTaskPanel.StatusTextLabel,
						Text:  "等待中",//已暂停 - 已完成40%
						Font: Font{
							Bold:      false,
							PointSize: 12,
							Family:    "Microsoft YaHei",
						},
					},
					ProgressBar{
						MinValue: 0,
						MaxValue: 100,
						MaxSize: Size{Height: 20},
						Value: 0,
						MarqueeMode: true,
						AssignTo: &uploadTaskPanel.ProgressBar,
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							TextLabel{
								Text: "正在上传: 等待中...",//正在上传：QLM牛逼.7z
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
								AssignTo: &uploadTaskPanel.CurrentCopyNameTextLabel,
							},
							TextLabel{
								Text: "当前速度: 计算中...", //当前速度：3M/s
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
								AssignTo: &uploadTaskPanel.SpeedTextLabel,
							},
							TextLabel{
								Text: "剩余时间: 计算中...", //剩余时间：11:45:14
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
								AssignTo: &uploadTaskPanel.TimeLeftTextLabel,
							},
							TextLabel{
								Text: "剩余项目: 计算中...", //剩余项目：1（1145.14MB）
								Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
								AssignTo: &uploadTaskPanel.ItemLeftTextLabel,
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
		mmw.ActiveTaskScrollEmpty.SetVisible(false)
	}

	return uploadTaskPanel, err
}



func RequestSuspendActiveTask(taskId string) {
	SuspendTaskIDChan <- taskId
}


func RequestResumeActiveTask(taskId string) {
	ResumeTaskIDChan <- taskId
}

func RequestAbortActiveTask(taskId string) {
	var dialog *walk.Dialog
	_, _ = Dialog{
		Icon: ImageResourcePath("upload.ico"),
		AssignTo: &dialog,
		Title:    "取消任务",
		Layout:   Grid{Columns: 1},
		Children: []Widget{
			TextLabel{
				Text: "确定要取消任务吗?",
				Font: Font{Bold: false, PointSize: 10, Family: "Microsoft YaHei"},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						Text: "确定",
						Font: Font{Bold: false, PointSize: 8, Family: "Microsoft YaHei"},
						OnClicked: func() {
							//walk.MsgBox(dialog, "提示", "暂不支持取消任务", walk.MsgBoxIconInformation)
							AbortTaskIDChan <- taskId
							dialog.Close(0)
						},
					},
					PushButton{
						Text: "取消",
						Font: Font{Bold: false, PointSize: 8, Family: "Microsoft YaHei"},
						OnClicked: func() {
							dialog.Close(0)
						},
					},
				},
			},
		},
	}.Run(mmw)
}