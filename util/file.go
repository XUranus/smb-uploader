package util

import (
	"fmt"
	"os"
	"path/filepath"
	"uploader/logger"
)

/**
	return name of dir of the path
 */
func DirName(path string) string {
	stat, err := os.Stat(path)
	if err != nil {
		logger.CommonLogger.Error("DirName", err)
	}
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

/**
	return path of the dir of the path
 */
func DirPath(path string) string {
	stat, err := os.Stat(path)
	if err != nil {
		logger.CommonLogger.Error("DirPath", err)
	}
	if stat.IsDir() {
		return path
	} else {
		return filepath.Dir(path)
	}
}

/**
	convert bytes num to file size string
	input: 1024 , return "1B"
 */
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

/**
	check is smb target path is available, for some smb url need authentication
 */
func CheckSMBTargetAvailability(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}