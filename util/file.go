package util

import (
	"fmt"
	"os"
	"path/filepath"
)

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
