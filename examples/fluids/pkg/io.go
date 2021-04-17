package fluids

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

const (
	// Image type
	JPEG = iota
	PNG
)

// Returns an image from a file
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

// Saves a jpeg or png image to a file with a given
// compression leve, these levels only work for png images
func SaveImage(path string, img image.Image) error {
	pathLower := strings.ToLower(path)

	var encodeMethod int
	if strings.HasSuffix(pathLower, ".jpeg") || strings.HasSuffix(pathLower, ".jpg") {
		encodeMethod = JPEG
	} else if strings.HasSuffix(pathLower, ".png") {
		encodeMethod = PNG
	} else {
		return errors.New("file must be .jpeg/.jpg or .png")
	}

	toimg, err := os.Create(path)
	if err != nil {
		return err
	}
	defer toimg.Close()

	switch encodeMethod {
	case JPEG:
		if err = jpeg.Encode(toimg, img, nil); err != nil {
			return err
		}
	case PNG:
		if err = png.Encode(toimg, img); err != nil {
			return err
		}
	}

	return nil
}
