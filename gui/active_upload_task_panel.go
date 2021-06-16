package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

// Upload Task Panel Control Block
type ActiveUploadTaskPanel struct {
	GroupBox						*walk.GroupBox
	SrcAndTargetLinkLabel 			*walk.LinkLabel
	StatusTextLabel 				*walk.TextLabel
	ContinueOrSuspendButton 		*walk.ToolButton
	CancelButton					*walk.ToolButton
	ProgressBar 					*walk.ProgressBar
	CurrentCopyNameTextLabel		*walk.TextLabel
	TimeLeftTextLabel				*walk.TextLabel
	ItemLeftTextLabel				*walk.TextLabel
	SpeedTextLabel					*walk.TextLabel
}


func NewActiveUploadTaskPanel(container walk.Container) (uploadTaskPanel *ActiveUploadTaskPanel, err error) {

	mw.MainWindow.Synchronize(func() {
		uploadTaskPanel = &ActiveUploadTaskPanel{}
		// bind group box
		uploadTaskPanel.GroupBox, err = walk.NewGroupBox(container)

		err = GroupBox{
			AssignTo: &uploadTaskPanel.GroupBox,
			Background: SystemColorBrush{Color: walk.SysColorWindow},
			Layout:     Grid{Columns: 1},
			Visible: false,
			Children: []Widget{
				LinkLabel{
					AssignTo: &uploadTaskPanel.SrcAndTargetLinkLabel,
					Font:     Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
					Text:     "",//`正在将 1 个项目从 <a href="">XXX</a> 复制到 <a href="">XXX</a>`,
					OnLinkActivated: func(link *walk.LinkLabelLink) {
						log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
					},
				},
				Composite{
					Layout: Grid{Columns: 4},
					Children: []Widget{
						TextLabel{
							AssignTo: &uploadTaskPanel.StatusTextLabel,
							Text:     "",//已暂停 - 已完成40%
							Font: Font{
								Bold:      false,
								PointSize: 12,
								Family:    "Microsoft YaHei",
							},
						},
						TextLabel{
							Text: "            ",
						},
						ToolButton{
							Text: "►",
							AssignTo: &uploadTaskPanel.ContinueOrSuspendButton,
							Background: SystemColorBrush{Color: walk.SysColorWindow},
						},
						ToolButton{
							Text: "✖",
							AssignTo: &uploadTaskPanel.CancelButton,
							Background: SystemColorBrush{Color: walk.SysColorBtnHighlight},
						},
					},
				},
				ProgressBar{
					MinValue: 0,
					MaxValue: 100,
					Value: 0,
					MarqueeMode: false,
					AssignTo: &uploadTaskPanel.ProgressBar,
				},
				Composite{
					Layout: Grid{Columns: 2},
					Children: []Widget{
						TextLabel{
							Text: "",//名称：QLM牛逼.7z
							Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							AssignTo: &uploadTaskPanel.CurrentCopyNameTextLabel,
						},
						TextLabel{
							Text: "", //名称：QLM牛逼.7z
							Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							AssignTo: &uploadTaskPanel.SpeedTextLabel,
						},
						TextLabel{
							Text: "", //剩余时间：11:45:14
							Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							AssignTo: &uploadTaskPanel.TimeLeftTextLabel,
						},
						TextLabel{
							Text: "", //剩余项目：1（1145.14MB）
							Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
							AssignTo: &uploadTaskPanel.ItemLeftTextLabel,
						},
					},
				},
			},
		}.Create(NewBuilder(container))
	})

	if err!= nil {
		fmt.Println(err)
	}

	return uploadTaskPanel, err
}
