package task

import (
	"errors"
	"os"
	"path/filepath"
	"uploader/logger"
)

/**
	FileCopyTask do file copy
 */
type FileCopyTask struct {

	Signal 						*RoutineSignal

	// static
	SourcePath					string
	TargetPath					string
	IsDir						bool
	UploadTaskRef    			*UploadTask

	// dynamic fields
	ItemsCopied					int64
	BytesCopied					int64
	CurrentCopyItemPath			string

	// events
	OnItemsCopiedChanged		func(int64)
	OnBytesCopiedChanged		func(int64)
	OnCurrentCopyItemChanged	func(string)
	OnExit						func(error)
}

func (copyTask *FileCopyTask) Start(async bool) {
	if async {
		go func() {
			copyTask.BlockStart()
		}()
	} else {
		copyTask.BlockStart()
	}
}


func (copyTask *FileCopyTask) BlockStart() {
	var err error = nil

	if copyTask.IsDir {
		err = copyTask.CopyFolder(copyTask.SourcePath, copyTask.TargetPath)

	} else {
		_, filename := filepath.Split(copyTask.SourcePath)
		finalTargetPath := filepath.Join(copyTask.TargetPath, filename)

		copyTask.CurrentCopyItemPath = finalTargetPath
		copyTask.OnCurrentCopyItemChanged(finalTargetPath)

		_, err = copyTask.CopyFile(copyTask.SourcePath, finalTargetPath)

		if err == nil {
			copyTask.ItemsCopied = 1
			copyTask.OnItemsCopiedChanged(1)
		}
	}

	copyTask.OnExit(err)
}


func (copyTask *FileCopyTask) CopyFolder(source string, dest string) (err error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourceFilePointer := filepath.Join(source, obj.Name())
		destinationFilePointer := filepath.Join(dest, obj.Name())
		if obj.IsDir() {
			err = copyTask.CopyFolder(sourceFilePointer, destinationFilePointer)
			if err != nil {
				return err
			}
		} else {
			copyTask.CurrentCopyItemPath = sourceFilePointer
			copyTask.OnCurrentCopyItemChanged(copyTask.CurrentCopyItemPath)

			_ , err = copyTask.CopyFile(sourceFilePointer, destinationFilePointer)

			if err != nil {
				return err
			} else {
				copyTask.ItemsCopied ++
				copyTask.OnItemsCopiedChanged(copyTask.ItemsCopied)
			}
		}
	}
	return
}

func (copyTask *FileCopyTask) CopyFile(source string, dest string) (written int64, err error) {
	copyTask.CurrentCopyItemPath = source
	written = 0
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	written, err = copyTask.copyBuffer(destFile, sourceFile)
	if err == nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			if sourceInfo != nil {
				err = os.Chmod(dest, sourceInfo.Mode())
			}
		}
	}


	// TODO:: err May Be "EOF", Cause Of This Error Hasn't Been Detected Yet!
	if err != nil && err.Error() == "EOF"  {
		err = nil
	}


	return written, err
}




type Writer interface {
	Write(p []byte) (n int, err error)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

type LimitedReader struct {
	Reader // underlying reader
	N int64  // max bytes remaining
}

var errInvalidWrite = errors.New("invalid write result")

var ErrShortWrite = errors.New("short write")

var EOF = errors.New("EOF")

func (copyTask *FileCopyTask) copyBuffer(dst Writer, src Reader) (written int64, err error) {
	var buf []byte
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	size := 32 * 1024
	if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf = make([]byte, size)
	for {

		if abort := copyTask.Signal.CheckSignal(); abort {
			logger.CommonLogger.Info("copyBuffer", "FileCopyTask received exit signal, return")
			return 0, AbortError
		}

		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)

			copyTask.BytesCopied += int64(nw)
			copyTask.OnBytesCopiedChanged(copyTask.BytesCopied)

			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

