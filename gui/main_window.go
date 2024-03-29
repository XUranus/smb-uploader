package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"uploader/db"
	"uploader/logger"
)

const resourceDir = "img"
var mmw	*MyMainWindow

type MyMainWindow struct {
	*walk.MainWindow

	NotifyIcon				*walk.NotifyIcon

	ActiveTaskScrollView 	*walk.ScrollView
	SucceedTaskScrollView 	*walk.ScrollView
	FailedTaskScrollView  	*walk.ScrollView

	ActiveTaskScrollEmpty 	*walk.Composite
	SucceedTaskScrollEmpty 	*walk.Composite
	FailedTaskScrollEmpty  	*walk.Composite
}

var (
	AbortTaskIDChan   chan string
	ResumeTaskIDChan  chan string
	SuspendTaskIDChan chan string
)


func InitMainWindow() {
	err := MainWindow{
		Icon: ImageResourcePath("upload.ico"),
		AssignTo: 	&mmw.MainWindow,
		Title:    	"上传任务列表",
		MinSize:  	Size{Width: 600, Height: 600},
		MaxSize: 	Size{Width: 600, Height: 1000},
		Visible: 	true,
		Layout: 	VBox{},
		Children: []Widget{
			TabWidget{
				Font: Font{Bold: false, PointSize: 12, Family:  "Microsoft YaHei"},
				Pages: []TabPage{
					{
						Title: "  进行中  ",
						Content: Composite{
							Layout: Grid{Columns: 1},
							Children: []Widget{
								Composite{
									AssignTo: &mmw.ActiveTaskScrollEmpty,
									Visible: true,
									Layout: VBox{},
									Children: []Widget{
										ImageView{
											MaxSize: Size{Height: 100, Width: 100},
											Mode: ImageViewModeShrink,
											Alignment: AlignHCenterVCenter,
											Image: ImageResourcePath("empty.png"),
										},
									},
								},
								ScrollView {
									Layout: VBox{MarginsZero: true},
									VerticalFixed: false,
									AssignTo: &mmw.ActiveTaskScrollView,
									Children: []Widget{},
								},
							},
						},
					}, {
						Title: "  已完成  ",
						Content: Composite{
							Layout: Grid{Columns: 1},
							Children: []Widget{
								Composite{
									AssignTo: &mmw.SucceedTaskScrollEmpty,
									Visible: true,
									Layout: VBox{},
									Children: []Widget{
										ImageView{
											MaxSize: Size{Height: 100, Width: 100},
											Mode: ImageViewModeShrink,
											Alignment: AlignHCenterVCenter,
											Image: ImageResourcePath("empty.png"),
										},
									},
								},
								ScrollView {
									Layout: VBox{MarginsZero: true},
									VerticalFixed: false,
									AssignTo: &mmw.SucceedTaskScrollView,
									Children: []Widget{},
								},
							},
						},
					}, {
						Title: "  已失败  ",
						Content: Composite{
							Layout: Grid{Columns: 1},
							Children: []Widget{
								Composite{
									AssignTo: &mmw.FailedTaskScrollEmpty,
									Visible: true,
									MaxSize: Size{Height: 150},
									Layout: VBox{},
									Children: []Widget{
										ImageView{
											MaxSize: Size{Height: 100, Width: 100},
											Mode: ImageViewModeShrink,
											Alignment: AlignHCenterVCenter,
											Image: ImageResourcePath("empty.png"),
										},
									},
								},
								ScrollView {
									Layout: VBox{MarginsZero: true},
									VerticalFixed: false,
									AssignTo: &mmw.FailedTaskScrollView,
									Children: []Widget{},
								},
							},
						},
					},
				},
			},
		},
	}.Create()
	if err != nil {
		logger.CommonLogger.Error("InitMainWindow", err)
	}

	// prevent from exit
	mmw.MainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason){
		mmw.MainWindow.SetVisible(false)
		*canceled = true
	})
}

func InitNotifyIcon() {
	// We load our icon from a file.
	icon, err := walk.Resources.Icon(ImageResourcePath("upload.ico"))
	if err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}

	// Create the notify icon and make sure we clean it up on exit.
	mmw.NotifyIcon, err = walk.NewNotifyIcon(mmw.MainWindow)
	if err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}
	//defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := mmw.NotifyIcon.SetIcon(icon); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}
	if err := mmw.NotifyIcon.SetToolTip("本地上传器"); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	mmw.NotifyIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		mmw.MainWindow.Show()
		_ = mmw.MainWindow.SetFocus()
	})

	// put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := mmw.NotifyIcon.ContextMenu().Actions().Add(exitAction); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}

	// put an about action into the context menu.
	aboutAction := walk.NewAction()
	if err := aboutAction.SetText("关于"); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}
	aboutAction.Triggered().Attach(func() {
		PopMessageBox("关于", "SMB本地上传器 https://github.com/XUranus/smb-uploader")
	})
	if err := mmw.NotifyIcon.ContextMenu().Actions().Add(aboutAction); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}


	// The notify icon is hidden initially, so we have to make it visible.
	if err := mmw.NotifyIcon.SetVisible(true); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := mmw.NotifyIcon.ShowInfo("本地上传器已启动", "点击查看所有任务列表"); err != nil {
		logger.CommonLogger.Error("InitNotifyIcon", err)
	}
}

func InitWindow()  {
	mmw = new(MyMainWindow)
	InitMainWindow()
	InitNotifyIcon()
}

func InitResourcePath(homePath string) {
	_ = walk.Resources.SetRootDirPath(homePath)
}


func InitInactiveTasksPanels(succeedTaskList []*db.UploadTaskRecord, failedTaskLists []*db.UploadTaskRecord) {
	InitSucceedTasksPanels(succeedTaskList)
	InitFailedTasksPanels(failedTaskLists)
}


func StartMainWindowBlock(suspendChan chan string, resumeChan chan string, abortChan chan string) {
	SuspendTaskIDChan = suspendChan
	ResumeTaskIDChan = resumeChan
	AbortTaskIDChan = abortChan
	mmw.Run()
}

