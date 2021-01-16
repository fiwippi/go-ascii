package convertor

import (
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/images"
)

func ConvertImage(imgPath, outputPath string, ac *ascii.AsciiConfig, cl images.CompressionLevel) error {
	// Processes the video frames into ascii frames
	img, err := images.ReadImage(imgPath)
	if err != nil {
		return err
	}

	generate := func(x, y int) ascii.RGBA {
		r, g, b, a := img.At(x, y).RGBA()
		r, g, b, a = r>>8, g>>8, b>>8, a>>8 // Colours
		return ascii.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	ascii_img, err := ac.GenerateAsciiImage(width, height, generate)
	if err != nil {
		return err
	}

	err = images.SaveImage(outputPath, ascii_img, cl)
	if err != nil {
		return err
	}

	return nil
}
