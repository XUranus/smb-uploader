package util

import (
	"fmt"
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

func FileSizeFromBytes(nBytes int64) string {
	if nBytes < 1024 {
		return fmt.Sprintf("%vB", nBytes)
	} else if nBytes < 1024*1024 {
		return fmt.Sprintf("%vKB", nBytes / 1024)
	} else if nBytes < 1024*1024*1024 {
		return fmt.Sprintf("%vMB", nBytes / (1024 * 1024))
	} else if nBytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%vGB", nBytes / (1024 * 1024 * 1024))
	} else {
		return fmt.Sprintf("%vTB", nBytes / (1024 * 1024 * 1024 * 1024))
	}
}
