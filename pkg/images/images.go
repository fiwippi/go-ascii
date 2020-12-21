package images

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func ReadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func SaveImage(path string, img image.Image) error {
	var encodeMethod int = 0
	if strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".jpg") {
		encodeMethod = 1
	} else if strings.HasSuffix(path, ".png") {
		encodeMethod = 2
	} else {
		return errors.New("File must be .jpeg/.jpg or .png")
	}

	toimg, err := os.Create(path)
	if err != nil {
		return err
	}
	defer toimg.Close()

	switch encodeMethod {
	case 1:
		if err = jpeg.Encode(toimg, img, nil); err != nil {
			return err
		}
	case 2:
		if err = png.Encode(toimg, img); err != nil {
			return err
		}
	}

	return nil
}
