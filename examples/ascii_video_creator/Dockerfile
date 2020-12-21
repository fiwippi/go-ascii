FROM golang:1.14

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p exes/

RUN go generate

RUN GOOS=linux GOARCH=amd64 go build -o exes/ascii_video_linux_x64
RUN GOOS=windows GOARCH=amd64 go build -o exes/ascii_video_windows_x64
