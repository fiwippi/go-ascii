 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate
BINARY_NAME=ascii_converter

build:
		# Get dependencies
		$(GOMOD) download
		$(GOGET) github.com/markbates/pkger/cmd/pkger

		# Build the pkger file for embedding the font
		$(GOGENERATE)

		# Build for linux and windows
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)_linux
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)_windows
clean:
		$(GOCLEAN)
		rm -R bin
