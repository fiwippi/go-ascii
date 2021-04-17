build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/fluids_linux
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/fluids_windows.exe
clean:
		go clean
		rm -R bin
