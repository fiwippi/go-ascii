@echo off
go generate

set GOOS=linux
set GOARCH=amd64
go build -o builds/converter_linux_amd64.exe

set GOOS=windows
set GOARCH=amd64
go build -o builds/converter_windows_amd64.exe

echo "Built