## ascii-video-creator

Takes an image as input and returns the image made out of coloured ascii characters

### config

- a `.env` must be supplied with paths pointing to ffmpeg and ffprobe, a sample `.env` file is provided
- the `.env` file must be in the same directory as the executable

### options
```
  -charset string
        Type of charset you want to use, 'limited' or 'extended'. Default is 'limited' (default "limited")
  -font string
        Path to a .ttf font file which the characters will be rendered as
  -fontsize float
        Font size in points (NOT pixels). Default is 14pt (default 14)
  -input string
        Path to the video you want to make ascii
  -output string
        Name of the output video you want to make e.g. 'test.mkv'
```

### build
1. get dependencies using
    - `go mod download`
    - `go get go get github.com/markbates/pkger/cmd/pkger`
2. run the `build.bat` script