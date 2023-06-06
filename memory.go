package ascii

const interpolationWeight float64 = 0.4

type coord struct {
	X, Y int
}

// Memory stores the brightnesses of pixels at
// specific coordinates.
//
// If passed in to ConvertWithOpts using the Interpolate
// option for successive calls, go-ascii interpolates the
// character used to represent the pixel.
//
// This is useful for converting multiple frames in a video
// where you might want the gradual change between characters
// to be less pronounced
type Memory struct {
	data map[coord]float64
}

// Reset clears any data the Memory may hold so that if used
// again with Interpolate, no interpolation would occur for
// the first call
func (m *Memory) Reset() {
	m.data = make(map[coord]float64)
}

func (m *Memory) interpolate(b float64, x, y int) float64 {
	c := coord{x, y}

	// If interpolation memory exists for this pixel then interpolate the brightness
	if oldBrightness, found := m.data[c]; found {
		b = (b * (1 - interpolationWeight)) + (oldBrightness * interpolationWeight)
	}

	// Store the new brightness value in memory
	m.data[c] = b

	// Return this new brightness
	return b
}
