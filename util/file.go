package util

import (
	"os"
	"path/filepath"
	"strings"
)

func getFileName(filepath string) string {
	fileNameSplit := strings.Split(filepath, "\\")
	return fileNameSplit[len(fileNameSplit) - 1]
}


func getFolderName(filepath string) string {
	folderNameSplit := strings.Split(filepath, "\\")
	return folderNameSplit[len(folderNameSplit) - 1]
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DirName(path string) string {
	stat, _ := os.Stat(path)
	ret := ""
	if stat.IsDir() {
		ret = stat.Name()
	} else {
		ret = DirName(filepath.Dir(path))
	}
	if ret == "\\" {
		ret = filepath.VolumeName(path)
	}
	return ret
}

func DirPath(path string) string {
	stat, _ := os.Stat(path)
	if stat.IsDir() {
		return path
	} else {
		return filepath.Dir(path)
	}
}

