package gui

import (
	"path/filepath"
	"strings"
)

func GetFileCategoryImage(path string, isDir bool) string {
	if isDir {
		return filepath.Join(resourceDir, "folder.png")
	}
	suffixes := strings.Split(path, ".")
	suffix := strings.ToLower(suffixes[len(suffixes) - 1])
	if suffix == "mp4" || suffix == "rmvb" || suffix == "avi" || suffix == "flv" {
		return filepath.Join(resourceDir, "media.png")
	} else if suffix == "exe" {
		return filepath.Join(resourceDir, "exe.png")
	} else if suffix == "zip" || suffix == "7z" || suffix == "iso" || suffix == "rar" || suffix == "gz" {
		return filepath.Join(resourceDir, "zip.png")
	} else if suffix == "jpg" || suffix == "jpeg" || suffix == "png" || suffix == "gif" || suffix == "ico" || suffix == "bmp"{
		return filepath.Join(resourceDir, "image.png")
	} else {
		return filepath.Join(resourceDir, "file.png")
	}
}


func ImageResourcePath(filename string) string {
	return filepath.Join(resourceDir, filename)
}