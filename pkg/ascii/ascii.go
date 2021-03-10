package ascii

import (
	"errors"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
)

const (
	CHAR_SET_LIMITED  = iota // Use the limited character set: " .:-=+*#%@"
	CHAR_SET_EXTENDED        // Use the extended character set: ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	CHAR_SET_BLOCK           // Use the block character set: "█"
	CHAR_SET_CUSTOM          // Use a custom character set supplied through AsciiConfig.CustomChatSet
)

// Coordinates of a pixel. Used in AsciiConfig to remember the colour
// of pixels at old coordinates so that their colour can be interpolated
type Coord struct {
	X, Y int
}

// Holds data used to render ascii images
type AsciiConfig struct {
	// Which charset to use in order to render the image
	CharSet int
	// The fontsize to draw in points (NOT pixels)
	FontSize float64
	// The bytes of a .ttf font which will be used to draw characters
	// on the image, example on loading a font is in the README.md
	FontBytes []byte
	// Whether to interpolate the colours, this makes the image smoother
	// but may introduce a delay in some cases
	Interpolate bool
	// How strong should the interpolation be 0-1 where 0 = no interpolation,
	// 1 will draw a constant colour so use any colours smaller than one
	InterpWeight float64
	// Slice holding the old brightness values for interpolation
	InterpMemory map[Coord]float64
	// Whether to draw on a transparent background, as opposed to the default black
	Transparency bool
	// A string of the custom charset, the left (0 index) of the
	// string is shown when the brightness is low and vice versa
	CustomCharSet string
}

// Creates a new ascii config struct with default settings
func NewAsciiConfig() *AsciiConfig {
	return &AsciiConfig{
		CharSet:      CHAR_SET_LIMITED,
		FontSize:     14,
		FontBytes:    nil,
		Interpolate:  true,
		InterpWeight: 0.4,
		InterpMemory: make(map[Coord]float64),
	}
}

// Draws a string character on a specific position at the image
func drawAsciiChar(img *image.RGBA, x, y int, char string, c *freetype.Context, fontsize float64, clr RGBA) error {
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))

	// Converts the coordinates to the fixed.Int26_6 coordinates that freetype uses
	pt := freetype.Pt(x, y+int(c.PointToFixed(fontsize)>>6))
	if _, err := c.DrawString(char, pt); err != nil {
		return err
	}
	return nil
}

// Returns an appropriate ascii string based on the brightness of a pixel
func (ac *AsciiConfig) brightnessToAscii(b uint8) string {
	if ac.CharSet == CHAR_SET_BLOCK {
		return "█"
	}

	var ascii string
	if ac.CharSet == CHAR_SET_EXTENDED {
		ascii = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	} else if ac.CharSet == CHAR_SET_LIMITED {
		ascii = " .:-=+*#%@"
	} else if ac.CharSet == CHAR_SET_CUSTOM {
		if len(ac.CustomCharSet) < 1 {
			return ""
		}
		ascii = ac.CustomCharSet
	}

	index := min(int(float64(b)/254*float64(len(ascii))), len(ascii)-1)
	return string(ascii[index])
}

// Generates an ascii image based on the configured AsciiConfig. This function uses getColour() to identify
// what colour should each pixel be and then draws the ascii characters in that colour. Check the README if
// you want to use getColour to draw a new ascii image from the existing one, if you want to use it to draw
// other simulations like state machines then check examples/fluid to find out how to do that
func (ac *AsciiConfig) GenerateAsciiImage(width, height int, getColour func(x, y int) RGBA) (image.Image, error) {
	// Ensure the interpolation memory exists
	if ac.Interpolate && ac.InterpMemory == nil {
		return nil, errors.New("No interpolation memory is available, either create this memory or turn interpolation off")
	}

	// Parse the initial image data
	bounds := image.Rect(0, 0, width, height)

	// Read the font data
	f, err := freetype.ParseFont(ac.FontBytes)
	if err != nil {
		return nil, err
	}

	// Create the font context
	c := freetype.NewContext()
	c.SetDPI(96)
	c.SetFont(f)
	c.SetFontSize(ac.FontSize)
	c.SetClip(bounds.Bounds())

	// Get the pixel width and height of the font
	opts := truetype.Options{Size: ac.FontSize}
	face := truetype.NewFace(f, &opts)

	// Height
	fontHeightPixel := face.Metrics().Height.Ceil() + face.Metrics().Descent.Ceil()

	// Width
	glyphBounds, _, found := face.GlyphBounds(rune('█'))
	if !found {
		return nil, errors.New("Failed getting font face width")
	}
	fontWidthPixel := glyphBounds.Max.X.Ceil() + face.Metrics().Descent.Ceil()

	// Create a new image to hold the ascii characters
	ascii_img := image.NewRGBA(bounds)
	var background = image.Black
	if ac.Transparency {
		background = image.Transparent
	}
	draw.Draw(ascii_img, ascii_img.Bounds(), background, image.Point{}, draw.Over)

	// Draw the new image
	for y := bounds.Min.Y; y < height; y += fontHeightPixel {
		for x := bounds.Min.X; x < width; x += fontWidthPixel {
			// Get the colour
			clr := getColour(x, y)

			// Get a brightness value for the image
			brightness := 0.299*float64(clr.R) + 0.587*float64(clr.G) + 0.114*float64(clr.B)

			// Interpolate the value if interpolation is turned on
			var interpolatedBrightness = brightness
			if ac.Interpolate {
				// If interpolation memory exists for this pixel then interpolate the brightness
				if oldBrightness, found := ac.InterpMemory[Coord{x, y}]; found {
					interpolatedBrightness = (float64(brightness) * (1 - ac.InterpWeight)) + (float64(oldBrightness) * ac.InterpWeight)
				}

				// Store the brightness value in memory
				ac.InterpMemory[Coord{x, y}] = interpolatedBrightness
			}

			// Get the ascii string for the corresponding brightness value
			err = drawAsciiChar(ascii_img, x, y, ac.brightnessToAscii(uint8(interpolatedBrightness)), c, ac.FontSize, clr)
			if err != nil {
				return nil, err
			}
		}
	}

	return ascii_img, nil
}
