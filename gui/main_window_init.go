package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"path/filepath"
)

// main window pointer, used to mainpulate gui widget
var mainWindow *walk.MainWindow = nil
var notifyIcon *walk.NotifyIcon = nil
var dialog *walk.Dialog
var TaskScrollView *walk.ScrollView = nil
//func doProgress() {
//	var dialog *walk.Dialog
//	var progressBar *walk.ProgressBar
//	err := Dialog{
//		AssignTo: &dialog,
//		Title:    "Progress dialog",
//		MinSize:  Size{Width: 300, Height: 200},
//		Layout:   VBox{},
//		Children: []Widget{
//			ProgressBar{AssignTo: &progressBar},
//		},
//	}.Create(mainWindow)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	dialog.Starting().Attach(func() {
//		go progressWorker(dialog, progressBar)
//	})
//	dialog.Run()
//}
//
//func progressWorker(dialog *walk.Dialog, progressBar *walk.ProgressBar) {
//	defer dialog.Synchronize(func() {
//		dialog.Close(0)
//	})
//
//	length := 10
//	dialog.Synchronize(func() {
//		progressBar.SetRange(0, int(length))
//	})
//	workWithCallback(length, func(n int64) {
//		fmt.Println("progress", n)
//		dialog.Synchronize(func() {
//			progressBar.SetValue(int(n))
//		})
//	})
//}
//
//func workWithCallback(length int, callback func(int64)) {
//	for i := 1; i <= length; i++ {
//		time.Sleep(time.Second)
//		callback(int64(i))
//	}
//}

// init main window
func init() {
	fmt.Println("init main window")
	err := MainWindow{
		AssignTo: &mainWindow,
		Title:    "Hidden Main Window",
		MinSize:  Size{Width: 1, Height: 1},
		Visible: false,
	}.Create()
	if err != nil {
		fmt.Println(err)
	}
}

// init notify icon
func init() {
	fmt.Println("init notify icon")

	// We load our icon from a file.
	icon, err := walk.Resources.Icon(filepath.Join("img","stop.ico"))
	if err != nil {
		log.Fatal(err)
	}

	// Create the notify icon and make sure we clean it up on exit.
	notifyIcon, err = walk.NewNotifyIcon(mainWindow)
	if err != nil {
		log.Fatal(err)
	}
	//defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := notifyIcon.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := notifyIcon.SetToolTip("本地上传器"); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	notifyIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		//if err := ni.ShowCustom(
		//	"Walk NotifyIcon Example",
		//	"There are multiple ShowX methods sporting different icons.",
		//	icon); err != nil {
		//
		//	log.Fatal(err)
		//}

		//_, err := walk.NewLabel(dialog)
		//if  err != nil {
		//	fmt.Println(err)
		//}

		fmt.Println("click icon")

		dialog.Show()
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := notifyIcon.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := notifyIcon.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := notifyIcon.ShowInfo("本地上传器已启动", "点击查看所有任务列表"); err != nil {
		log.Fatal(err)
	}
}



func initTaskStatusList() {
	//image, err:= walk.ImageFrom(filepath.Join("img", "cancel.png"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	err := Dialog{
		AssignTo: &dialog,
		Title:    "任务列表",
		MinSize:  Size{Width: 500, Height: 1000},
		MaxSize:  Size{Width: 600, Height: 1200},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			TabWidget{
				Font: Font{Bold:false, PointSize: 12, Family:  "Microsoft YaHei"},
				Pages: []TabPage{
					TabPage {
						Title:"  进行中  ",
						Content: ScrollView{
							Layout: VBox{MarginsZero: true},
							VerticalFixed: false,
							AssignTo: &TaskScrollView,
							Children: []Widget{
								GroupBox{
									Background: SystemColorBrush{Color: walk.SysColorWindow},
									Layout:     Grid{Columns: 1},
									Children: []Widget{
										LinkLabel{
											Font:     Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
											Text:     `正在将 1 个项目从 <a href="">XXX</a> 复制到 <a href="">XXX</a>`,
											OnLinkActivated: func(link *walk.LinkLabelLink) {
												log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
											},
										},
										Composite{
											Layout: Grid{Columns: 4},
											Children: []Widget{
												TextLabel{
													Text:     "已暂停 - 已完成40%",
													Font: Font{
														Bold:      false,
														PointSize: 12,
														Family:    "Microsoft YaHei",
													},
												},
												TextLabel{Text: "            "},
												ToolButton{Text: "►",  Background: SystemColorBrush{Color: walk.SysColorWindow}},
												ToolButton{Text: "✖", Background: SystemColorBrush{Color: walk.SysColorBtnHighlight}},
											},
										},
										ProgressBar{MinValue: 0, MaxValue: 100, Value: 30, Name: "www", MarqueeMode: true},
										Composite{
											Layout: Grid{Columns: 2},
											Children: []Widget{
												TextLabel{Text: "名称：QLM牛逼.7z", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "速度：10MB/S", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余时间：11:45:14", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余项目：1（1145.14MB）", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
											},
										},
									},
								},
							},
						},
					},
					TabPage {
						Title: "  已完成  ",
						Content: ScrollView{
							Layout: VBox{MarginsZero: true},
							VerticalFixed: false,
							AssignTo: &TaskScrollView,
							Children: []Widget{
								GroupBox{
									Background: SystemColorBrush{Color: walk.SysColorWindow},
									Layout:     Grid{Columns: 1},
									Children: []Widget{
										LinkLabel{
											Font:     Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
											Text:     `正在将 1 个项目从 <a href="">XXX</a> 复制到 <a href="">XXX</a>`,
											OnLinkActivated: func(link *walk.LinkLabelLink) {
												log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
											},
										},
										Composite{
											Layout: Grid{Columns: 4},
											Children: []Widget{
												TextLabel{
													Text:     "已暂停 - 已完成40%",
													Font: Font{
														Bold:      false,
														PointSize: 12,
														Family:    "Microsoft YaHei",
													},
												},
												TextLabel{Text: "            "},
												ToolButton{Text: "►",  Background: SystemColorBrush{Color: walk.SysColorWindow}},
												ToolButton{Text: "✖", Background: SystemColorBrush{Color: walk.SysColorBtnHighlight}},
											},
										},
										ProgressBar{MinValue: 0, MaxValue: 100, Value: 30, Name: "www", MarqueeMode: true},
										Composite{
											Layout: Grid{Columns: 2},
											Children: []Widget{
												TextLabel{Text: "名称：QLM牛逼.7z", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "速度：10MB/S", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余时间：11:45:14", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余项目：1（1145.14MB）", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
											},
										},
									},
								},
							},
						},
					},
					TabPage {
						Title: "  已失败  ",
						Content: ScrollView{
							Layout: VBox{MarginsZero: true},
							VerticalFixed: false,
							AssignTo: &TaskScrollView,
							Children: []Widget{
								GroupBox{
									Background: SystemColorBrush{Color: walk.SysColorWindow},
									Layout:     Grid{Columns: 1},
									Children: []Widget{
										LinkLabel{
											Font:     Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"},
											Text:     `正在将 1 个项目从 <a href="">XXX</a> 复制到 <a href="">XXX</a>`,
											OnLinkActivated: func(link *walk.LinkLabelLink) {
												log.Printf("id: '%s', url: '%s'\n", link.Id(), link.URL())
											},
										},
										Composite{
											Layout: Grid{Columns: 4},
											Children: []Widget{
												TextLabel{
													Text:     "已暂停 - 已完成40%",
													Font: Font{
														Bold:      false,
														PointSize: 12,
														Family:    "Microsoft YaHei",
													},
												},
												TextLabel{Text: "            "},
												ToolButton{Text: "►",  Background: SystemColorBrush{Color: walk.SysColorWindow}},
												ToolButton{Text: "✖", Background: SystemColorBrush{Color: walk.SysColorBtnHighlight}},
											},
										},
										ProgressBar{MinValue: 0, MaxValue: 100, Value: 30, Name: "www", MarqueeMode: true},
										Composite{
											Layout: Grid{Columns: 2},
											Children: []Widget{
												TextLabel{Text: "名称：QLM牛逼.7z", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "速度：10MB/S", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余时间：11:45:14", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
												TextLabel{Text: "剩余项目：1（1145.14MB）", Font: Font{Bold: false, PointSize: 9, Family: "Microsoft YaHei"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Visible: true,
	}.Create(mainWindow)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("create dialog")

	// disable close event
	dialog.Closing().Attach(func(canceled *bool, reason walk.CloseReason){
		dialog.SetVisible(false)
		*canceled = true
	})


	//dialog.Run()
}


func StartMainWindow() {
	fmt.Println("start to Run Loop")
	//defer ni.Dispose()


	initTaskStatusList()

	fmt.Println("run main window")
	mainWindow.Run()

	fmt.Println("exit")

}