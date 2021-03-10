package fluids

import (
	"github.com/fiwippi/ascii-image-creator/pkg/ascii"
	"github.com/lucasb-eyer/go-colorful"
	"math/rand"
)

// Whether to use random colours for the grid
var randomColours bool

// Clamps a float to be in the range [0, 1]
func clampRGB(a float64) float64 {
	if a < 0 {
		return 0
	} else if a > 1 {
		return 1
	}
	return a
}

// Returns a color for a cell in the grid based on its coordinates
func (g Grid) getColour(x, y int) ascii.RGBA {
	// If no fluid then black
	if g.fluids[x][y] == 0 {
		return ascii.RGBA{0, 0, 0, 255}
	}

	// Get the brightness of the colour
	lum := float64(g.fluids[x][y]) / float64(fluidMax)

	// Seed the random generator using the coordinates of the point
	rand.Seed(getSeed(x, y))

	// Colour calculation method 1 - Background all merges to one colour
	// Here the background is majority blue, to make it more red, green etc.
	// then you need to change the values given to redf, greenf and bluef.
	// For example, for a red background use, (120, 20, 20) instead of the
	// current (20, 20, 120)
	redf := clampRGB(float64(20)/255 + (rand.Float64() * lum))
	greenf := clampRGB(float64(20)/255 + (rand.Float64() * lum))
	bluef := clampRGB(float64(120)/255 + (rand.Float64() * lum))
	clr := colorful.Color{redf, greenf, bluef}
	red, green, blue := clr.RGB255()

	return ascii.RGBA{R: red, G: green, B: blue, A: 255}

	//// Colour calculation method 2 - Background is darker version of fluid
	//// Calculate luminosity and the threshold for full hue and saturation
	//thresh := float64((fluidMax)-1) / float64(fluidMax)
	//
	//// Calculate the hue
	//rand.Seed(int64(x) + int64(y))
	//hue := rand.Float64() * 360
	//
	//// If the cell is full
	//if lum > thresh {
	//	red, green, blue := colorful.Hcl(hue, 0.75, 0.8).Clamped().RGB255()
	//	return ascii.RGBA{R: red, G: green, B: blue, A: 255}
	//}
	//// If the cell isn't full then make it darker depending on how empty it is
	//h, s, v := colorful.Hcl(hue, 0.75, 0.8).Clamped().Hsv()
	//red, green, blue := colorful.Hsv(h, s, 0.05 + (lum * (v - 0.05))).Clamped().RGB255()
	//
	//return ascii.RGBA{R: red, G: green, B: blue, A: 255}
}
