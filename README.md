# Go Ascii
## Overview
Library which renders images using ascii characters

## Install
```
go get github.com/fiwippi/go-ascii
```

## Usage
### Default
```go
// Read in an image...

// Generate the ascii image
conf := ascii.DefaultConfig()
asciiImg, err := ascii.Convert(img, conf, nil)
if err != nil {
    log.Fatal(err)
}
```

### With Interpolation
```go
// Given a slice of images

// Generate the interpolated images
mem := &ascii.Memory{}
conf := ascii.DefaultConfig()
for _, img := range images {
    asciiImg, err := ascii.Convert(img, conf, mem)
    if err != nil {
        log.Fatalln(err)
    }
}
```

### With Custom Font
⚠️ - `go-ascii` expects monospace fonts!
```go
// Read in the font file
data, err := os.ReadFile("font_file.ttf")
if err != nil {
    log.Fatal(err)
}

// Parse the font
font, err := opentype.Parse(data)
if err != nil {
    log.Fatal(err)
}

// Perform the conversion
conf := ascii.DefaultConfig()
conf.Font = font
asciiImg, err := ascii.Convert(img, conf, nil)
if err != nil {
    log.Fatal(err)
}
```

## Examples
![example 1](assets/1.jpeg)

![example 2](assets/2.jpeg)

To convert videos check out the example at [examples/video](examples/video)

![example 3](examples/video/assets/explosion.gif)

## License
`BSD-3-Clause`