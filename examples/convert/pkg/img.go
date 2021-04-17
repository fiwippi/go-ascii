package convert

import (
	"github.com/fiwippi/go-ascii"
)

func ConvertImage(imgPath, outputPath string, ac *ascii.AsciiConfig) error {
	// Processes the video frames into ascii frames
	img, err := ReadImage(imgPath)
	if err != nil {
		return err
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	asciiImg, err := ac.GenerateAsciiImage(width, height, ascii.ImgColours(img))
	if err != nil {
		return err
	}

	err = SaveImage(outputPath, asciiImg)
	if err != nil {
		return err
	}

	return nil
}
