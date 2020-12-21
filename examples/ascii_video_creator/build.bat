@echo off
go generate

set GOOS=linux
set GOARCH=amd64
go build -o builds/video_converter_linux_amd64.exe

set GOOS=windows
set GOARCH=amd64
go build -o builds/video_converter_windows_amd64.exe

echo Built Video Convertor