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

type CompressionLevel int

const (
	JPEG = iota
	PNG

	BestCompression = iota
	DefaultCompression
	BestSpeed
)

func PNGCompressionLevel(level CompressionLevel) (png.CompressionLevel, error) {
	switch level {
	case BestCompression:
		return png.BestCompression, nil
	case BestSpeed:
		return png.BestSpeed, nil
	case DefaultCompression:
		return png.DefaultCompression, nil
	}

	return 0, errors.New("Cannot convert compression level")
}

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

func SaveImage(path string, img image.Image, cls ...CompressionLevel) error {
	pathLower := strings.ToLower(path)

	var cl CompressionLevel
	if len(cls) > 0 {
		cl = cls[0]
	}

	var encodeMethod int
	if strings.HasSuffix(pathLower, ".jpeg") || strings.HasSuffix(pathLower, ".jpg") {
		encodeMethod = JPEG
	} else if strings.HasSuffix(pathLower, ".png") {
		encodeMethod = PNG
	} else {
		return errors.New("File must be .jpeg/.jpg or .png")
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
		pngCL, err := PNGCompressionLevel(cl)
		if err != nil {
			return err
		}

		var pngEnc = &png.Encoder{
			CompressionLevel: pngCL,
		}

		if err = pngEnc.Encode(toimg, img); err != nil {
			return err
		}
	}

	return nil
}
