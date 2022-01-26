package ascii

import (
	"errors"
	"fmt"
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Convert draws a given image using ascii characters based on the parameters
// in conf. If you want to interpolate the characters you can supply the memory.
// This will make the change in characters between frames look more natural.
func Convert(img image.Image, conf Config, memory *Memory) (image.Image, error) {
	// Validate interpolation params
	var interpolate bool
	if memory != nil {
		if conf.InterpolateWeight < 0 || conf.InterpolateWeight > 1 {
			return nil, fmt.Errorf("interpolation weight should be between 0 and 1 (inclusive)")
		}

		interpolate = true
		if memory.data == nil {
			memory.Reset()
		}
	}

	// Ensure the font exists
	if conf.Font == nil {
		return nil, fmt.Errorf("no font specified")
	}

	// Create the font face
	opts := &opentype.FaceOptions{
		Size:    conf.FontSize,
		DPI:     72,
		Hinting: 0,
	}
	face, err := opentype.NewFace(conf.Font, opts)
	if err != nil {
		return nil, err
	}

	// Calculate the font width and height
	height := face.Metrics().Ascent.Round()
	glyphBounds, _, found := face.GlyphBounds('█')
	if !found {
		return nil, errors.New("failed getting font face width")
	}
	width := glyphBounds.Max.X.Round()

	// Create the new image
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	var background = image.Black
	if conf.Transparency {
		background = image.Transparent
	}
	draw.Draw(newImg, bounds.Bounds(), background, image.Point{}, draw.Over)

	// Draw onto the new image
	for y := bounds.Min.Y; y < bounds.Max.Y; y += height {
		for x := bounds.Min.X; x < bounds.Max.X; x += width {
			// Get the colour
			clr := img.At(x, y)
			r, g, b, a := img.At(x, y).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8

			// Get a brightness value for the image
			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			// Interpolate the value if interpolation is turned on
			if interpolate {
				// If interpolation memory exists for this pixel then interpolate the brightness
				if oldBrightness, found := memory.data[coord{x, y}]; found {
					brightness = (brightness * (1 - conf.InterpolateWeight)) + (oldBrightness * conf.InterpolateWeight)
				}

				// Store the brightness value in memory
				memory.data[coord{x, y}] = brightness
			}

			// Draw the string
			(&font.Drawer{
				Dst:  newImg,
				Src:  image.NewUniform(clr),
				Face: face,
				Dot:  fixed.P(x, y),
			}).DrawString(conf.CharSet.parseBrightness(uint8(brightness)))
		}
	}

	return newImg, nil
}
