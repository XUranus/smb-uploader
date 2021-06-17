package gui

import "github.com/lxn/walk"

func PopMessageBox(title string, message string) {
	walk.MsgBox(nil, title, message, walk.MsgBoxIconInformation)
}