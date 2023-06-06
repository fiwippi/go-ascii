package ascii

import (
	"fmt"
	"image"
	"image/draw"
)

// Convert renders the given image using ascii characters
// using default options. If you want to select specific
// configuration parameters such as font size you should
// use ConvertWithOpts
func Convert(img image.Image) (image.Image, error) {
	return ConvertWithOpts(img)
}

// ConvertWithOpts renders the given image using ascii characters.
//
// You can pass in Option(s) to configure the settings which the
// renderer users.
func ConvertWithOpts(img image.Image, opts ...Option) (image.Image, error) {
	// Ensure image exists
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	// Create the default options
	defOpts := &options{
		font:    defaultFont,
		charset: CharsetExtended,
		fontPts: 14,
	}

	// Change options according to modifiers
	for _, setter := range opts {
		if setter == nil {
			return nil, fmt.Errorf("option supplied is nil")
		}

		err := setter(defOpts)
		if err != nil {
			return nil, err
		}
	}

	// Perform the conversion
	return convert(img, defOpts)
}

func convert(img image.Image, opts *options) (image.Image, error) {
	// Calculate the font face metrics
	pf, err := parseFont(opts.font, opts.fontPts)
	if err != nil {
		return nil, err
	}

	// Create the new image
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds.Bounds(), image.Black, image.Point{}, draw.Over)

	// Convert the charset to its runes, the conversion to
	//runes is done so that unicode characters can be indexed
	// appropriately instead of individual code points
	rs := []rune(opts.charset)

	// Loop over the new image's coordinates
	for y := bounds.Min.Y; y < bounds.Max.Y; y += pf.height {
		for x := bounds.Min.X; x < bounds.Max.X; x += pf.width {
			// Get the colour
			clr := img.At(x, y)

			// Scale the values from 0-65535 to 0-255
			r, g, b, _ := clr.RGBA()
			r, g, b = r>>8, g>>8, b>>8

			// Get a brightness value of the colour from
			// here: https://www.w3.org/TR/AERT/#color-contrast
			//
			// This method isn't super-accurate since the
			// standards used are dated but this is negligible
			// in terms of the final image produced
			bright := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			// Interpolate if memory is specified
			if opts.mem != nil {
				bright = opts.mem.interpolate(bright, x, y)
			}

			// Scale brightness in range 0-1
			bright /= 255
			if bright > 1.0 {
				bright = 1.0
			}

			// Use that as a percentage of the charset's length
			// to get the index of the respective rune
			index := int(bright * float64(len(rs)-1))
			if index > len(rs)-1 {
				index = len(rs) - 1
			}

			// Draw the rune
			char := string(rs[index])
			pf.drawString(char, clr, newImg, x, y)
		}
	}

	return newImg, nil
}
