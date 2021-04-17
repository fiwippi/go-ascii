build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ascii_linux
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/ascii_windows
clean:
		go clean
		rm -R bin
