package gui

import "github.com/lxn/walk"



func ShowCustomNotify(title string, message string) {
	icon, _ := walk.Resources.Icon(ImageResourcePath("upload.ico"))
	_ = mmw.NotifyIcon.ShowCustom(title, message, icon)
}

func RefreshMainWindow() {
	_ = mmw.MainWindow.SetSize(mmw.Size())
}

func PopMessageBox(title string, message string) {
	walk.MsgBox(nil, title, message, walk.MsgBoxIconInformation)
}

func CheckScrollViewEmpty(scrollView *walk.ScrollView, emptyPanel *walk.Composite)  {
	if scrollView.Children().Len() == 0 {
		emptyPanel.SetVisible(true)
	}
}

func GetMyMainWindow() *MyMainWindow {
	return mmw
}

//func LoadCompleted() bool {
//	if mmw == nil {
//		return false
//	}
//	if mmw.MainWindow == nil {
//		return false
//	}
//	if mmw.NotifyIcon == nil || mmw.FailedTaskScrollView == nil || mmw.SucceedTaskScrollView == nil || mmw.ActiveTaskScrollView == nil {
//		return false
//	}
//	return true
//}