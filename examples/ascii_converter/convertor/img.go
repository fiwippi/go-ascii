package convertor

import (
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/images"
)

func ConvertImage(imgPath, outputPath string, ac *ascii.AsciiConfig) error {
	// Processes the video frames into ascii frames
	img, err := images.ReadImage(imgPath)
	if err != nil {
		return err
	}

	generate := func(x, y int) ascii.RGB {
		r, g, b, _ := img.At(x, y).RGBA()
		r, g, b = r>>8, g>>8, b>>8 // Colours
		return ascii.RGB{uint8(r), uint8(g), uint8(b)}
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	ascii_img, err := ac.GenerateAsciiImage(width, height, generate)
	if err != nil {
		return err
	}

	err = images.SaveImage(outputPath, ascii_img)
	if err != nil {
		return err
	}

	return nil
}
