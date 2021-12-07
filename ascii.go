package ascii

import (
	"errors"
	"image"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const (
	CharSetLimited  = iota // Use the limited character set: " .:-=+*#%@"
	CharSetExtended        // Use the extended character set: ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	CharSetBlock           // Use the block character set: "█"
	CharSetCustom          // Use a custom character set supplied through AsciiConfig.CustomChatSet
)

// TODO remove need for FreeType

// Coordinates of a pixel. Used in AsciiConfig to remember the colour
// of pixels at old coordinates so that their colour can be interpolated
type coord struct {
	X, Y int
}

// AsciiConfig holds data used to render ascii images
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
	InterpMemory map[coord]float64
	// Whether to draw on a transparent background, as opposed to the default black
	Transparency bool
	// A string of the custom charset, the left (0 index) of the
	// string is shown when the brightness is low and vice versa
	CustomCharSet string
}

// NewAsciiConfig creates a new ascii config struct with default settings
func NewAsciiConfig() *AsciiConfig {
	return &AsciiConfig{
		CharSet:      CharSetLimited,
		FontSize:     14,
		FontBytes:    nil,
		Interpolate:  true,
		InterpWeight: 0.4,
		InterpMemory: make(map[coord]float64),
	}
}

// drawAsciiChar draws a string character on a specific position at the image
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

// brightnessToAscii returns an appropriate ascii string based on the brightness of a pixel
func (ac *AsciiConfig) brightnessToAscii(b uint8) string {
	if ac.CharSet == CharSetBlock {
		return "█"
	}

	var ascii string
	if ac.CharSet == CharSetExtended {
		ascii = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	} else if ac.CharSet == CharSetLimited {
		ascii = " .:-=+*#%@"
	} else if ac.CharSet == CharSetCustom {
		if len(ac.CustomCharSet) < 1 {
			return ""
		}
		ascii = ac.CustomCharSet
	}

	index := min(int(float64(b)/254*float64(len(ascii))), len(ascii)-1)
	return string(ascii[index])
}

// ConvertImage converts a given image into ascii based on the configured AsciiConfig
func (ac *AsciiConfig) ConvertImage(img image.Image) (image.Image, error) {
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	asciiImg, err := ac.GenerateAsciiImage(width, height, ImgColours(img))
	if err != nil {
		return nil, err
	}
	return asciiImg, nil
}

// GenerateAsciiImage generates an ascii image based on the configured AsciiConfig. This function uses getColour() to identify
// what colour should each pixel be and then draws the ascii characters in that colour. This function is more low-leve than
// ConverTImage and it can be used to draw other simulations like state machines, check ./examples/fluids for more details
func (ac *AsciiConfig) GenerateAsciiImage(width, height int, getColour func(x, y int) RGBA) (image.Image, error) {
	// Ensure the interpolation memory exists
	if ac.Interpolate && ac.InterpMemory == nil {
		return nil, errors.New("no interpolation memory is available, either create this memory or turn interpolation off")
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
	glyphBounds, _, found := face.GlyphBounds('█')
	if !found {
		return nil, errors.New("failed getting font face width")
	}
	fontWidthPixel := glyphBounds.Max.X.Ceil() + face.Metrics().Descent.Ceil()

	// Create a new image to hold the ascii characters
	asciiImg := image.NewRGBA(bounds)
	var background = image.Black
	if ac.Transparency {
		background = image.Transparent
	}
	draw.Draw(asciiImg, asciiImg.Bounds(), background, image.Point{}, draw.Over)

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
				if oldBrightness, found := ac.InterpMemory[coord{x, y}]; found {
					interpolatedBrightness = (float64(brightness) * (1 - ac.InterpWeight)) + (float64(oldBrightness) * ac.InterpWeight)
				}

				// Store the brightness value in memory
				ac.InterpMemory[coord{x, y}] = interpolatedBrightness
			}

			// Get the ascii string for the corresponding brightness value
			err = drawAsciiChar(asciiImg, x, y, ac.brightnessToAscii(uint8(interpolatedBrightness)), c, ac.FontSize, clr)
			if err != nil {
				return nil, err
			}
		}
	}

	return asciiImg, nil
}
