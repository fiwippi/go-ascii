package ascii

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
