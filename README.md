# Go-Ascii

## Overview

Library which renders images using ascii characters

## Install

```
go get github.com/fiwippi/go-ascii
```

## Usage

Errors ignored for brevity

### Default

```go
// Read in an image...
img := ...

// Generate the ascii version
asciiImg, _ := ascii.Convert(img)
```

### With Interpolation

```go
// Given a slice of images
images := ...

// Generate the interpolated images
mem := &ascii.Memory{}
for _, img := range images {
    asciiImg, _ := ascii.ConvertWithOpts(img, ascii.Interpolate(Mem))
}
```

### With Custom Font

> **Warning**
> `go-ascii` expects monospace fonts!

```go
// Read in the font file
data, _ := os.ReadFile("font_file.ttf")

// Parse the font
font, _ := opentype.Parse(data)

// Perform the conversion
asciiImg, _ := ascii.ConvertWithOpts(img, ascii.Font(font))
```

## Examples

![example 1](assets/1.jpeg)

![example 2](assets/2.jpeg)

To convert videos check out the example at [examples/video](examples/video)

![example 3](examples/video/assets/explosion.gif)

## License

`BSD-3-Clause`