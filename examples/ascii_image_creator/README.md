## ascii-image-creator

Takes an image as input and returns the image made out of coloured ascii characters

### options

```
  -charset string
        Type of charset you want to use, 'limited' or 'extended' (default "limited")
  -font string
        Path to a .ttf font file which the characters will be rendered as. If empty, 'inconsolata bold' is used
  -fontsize float
        Font size in points (NOT pixels) (default 14)
  -input string
        Path to the image you want to make ascii. Must be jpeg or png
  -output string
        Name of the output image you want to make e.g. 'test.jpg'. Can be jpeg or png
```

### build
1. get dependencies using
    - `go mod download`
    - `go get go get github.com/markbates/pkger/cmd/pkger`
2. run the `build.bat` script