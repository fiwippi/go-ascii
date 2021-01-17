## ascii-image-creator

library which takes an image as input and returns the image made out of coloured ascii characters

### installation

```
go get github.com/nadav-rahimi/ascii-image-creator
```

### usage
```go
// Get the font file as bytes and reading its data
fontBytes, err := ioutil.ReadFile("font_file.ttf")
if err != nil {
    log.Fatal("Error reading font data: ", err)
}

// Set up the ascii config
ac := &ascii.AsciiConfig{
    CharSet:   ascii.CHAR_SET_LIMITED,
    FontBytes: fontBytes,
    FontSize:  14,
}

// Reads in the image
img, err := images.ReadImage("image.png")
if err != nil {
    log.Fatal(err)
}

generate := func(x, y int) ascii.RGB {
    r, g, b, a := img.At(x, y).RGBA()
    r, g, b, a = r>>8, g>>8, b>>8, a>>8 // Colours must be 8 bit
    return ascii.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

width, height := img.Bounds().Max.X, img.Bounds().Max.Y
ascii_img, err := ac.GenerateAsciiImage(width, height, generate)
if err != nil {
    log.Fatal(err)
}

err = images.SaveImage(outputPath, ascii_img)
if err != nil {
    log.Fatal(err)
}
```

### todo 

- [x] video interpolation
- [x] alpha channel support + ability to draw on a transparent background
- [x] render fonts less blurry