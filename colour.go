package ascii

import "image"

// Colour struct used to draw the characters
type RGBA struct {
	R, G, B, A uint8
}

// Implements so that image.NewUniform can be used
func (rgba RGBA) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(rgba.R)
	r |= r << 8
	g := uint32(rgba.G)
	g |= g << 8
	b := uint32(rgba.B)
	b |= b << 8
	a := uint32(rgba.A)
	a |= a << 8

	return r, g, b, a
}

// Wrapper for the generate parameter to get the colours from an image
func ImgColours(img image.Image) func(x, y int) RGBA {
	return func(x, y int) RGBA {
		r, g, b, a := img.At(x, y).RGBA()
		r, g, b, a = r>>8, g>>8, b>>8, a>>8 // Colours
		return RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
}
