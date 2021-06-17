windres -o resource.syso resource.rc
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags="-H windowsgui"