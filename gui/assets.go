package gui

import (
	"path/filepath"
	"strings"
)

var assetsPath = ".\\img"

func GetFileCategoryImage(path string, isDir bool) string {
	if isDir {
		return filepath.Join(assetsPath, "folder.png")
	}
	suffixes := strings.Split(path, ".")
	suffix := strings.ToLower(suffixes[len(suffixes) - 1])
	if suffix == "mp4" || suffix == "rmvb" || suffix == "avi" || suffix == "flv" {
		return filepath.Join(assetsPath, "media.png")
	} else if suffix == "exe" {
		return filepath.Join(assetsPath, "exe.png")
	} else if suffix == "zip" || suffix == "7z" || suffix == "iso" || suffix == "rar" || suffix == "gz" {
		return filepath.Join(assetsPath, "zip.png")
	} else if suffix == "jpg" || suffix == "jpeg" || suffix == "png" || suffix == "gif" || suffix == "ico" || suffix == "bmp"{
		return filepath.Join(assetsPath, "image.png")
	} else {
		return filepath.Join(assetsPath, "file.png")
	}
}


func ImageResourcePath(filename string) string {
	return filepath.Join(assetsPath, filename)
}