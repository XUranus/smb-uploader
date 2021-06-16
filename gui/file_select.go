package gui

import (
	"github.com/lxn/walk"
)

func SelectedPath(isDir bool, fileFilter string) (path string, ok bool, err error){

	if isDir {
		dlg := walk.FileDialog{
			Title: "选择上传的目录",
			Filter: "*.*",
		}
		if ok, err := dlg.ShowBrowseFolder(mmw); err != nil {
			return "", false, err
		} else if !ok {
			return "", false, nil
		}
		return dlg.FilePath, true, nil

	} else {
		if fileFilter == "" {
			fileFilter = "所有文件|*"
		}
		dlg := walk.FileDialog{
			Title: "选择上传的文件",
			Filter: fileFilter, //"压缩文件(7z/rar/zip)|*.7z;*.zip;*.rar",
		}
		if ok, err := dlg.ShowOpen(mmw); err != nil {
			return "", false, err
		} else if !ok {
			return "", false, nil
		}
		return dlg.FilePath, true, nil
	}

	//dlg.FilePath = initialDirPath // TODO::RestorePath not support in walk yet, temporally meaning less
}
