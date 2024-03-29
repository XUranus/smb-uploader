# smb-uploader
a local uploader program, used to be invoked by Web to upload big file and folder to smb server

![](screenshot.png)

## Features
 - [x] GUI support
 - [x] Suspend running task
 - [x] Resume running task
 - [x] Abort running task
 - [x] Http interface

## Build(X86 32)
Requires:
 - Windows10
 - Go > 1.6
 - GCCGO
```
windres -o resource.syso resource.rc
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags="-H windowsgui"
```
or directly run `build.bat`


After `uploader.exe` build accomplished, run `install.bat` to install files to `C:\Program Files(X86)\SmbUploader`


## Getting Start
visit `smbuploader://` in browser, or directly run `C:\Programs Files(X86)\uploader.exe`, then make a `POST` request to `http://127.0.0.1:8888`:
```
{
	"targetPath": "\\\\192.168.3.5\\Users\\A\\Desktop",
	"isDir": true,
}
```
`targetPath` is your target smb url (or you can test with a local path like `C:\Users\A\Desktop` without smb server)


## TODO
 - [x] GUI support
 - [x] Suspend running task
 - [x] Resume running task
 - [x] Abort running task
 - [x] Http interface
 - [x] Single process
 - [x] Windows registry script 
 - [x] Log modules
 - [ ] Recover task from error logs

![](smb.png)
