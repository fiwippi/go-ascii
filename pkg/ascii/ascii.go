package ascii

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	_ "github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"image"
	_ "image/draw"
	_ "io/ioutil"
)

// TODO sample the colour instead of using the first input

const (
	CHAR_SET_LIMITED = iota
	CHAR_SET_EXTENDED
	CHAR_SET_BLOCK
)

type RGB struct {
	R, G, B uint8
}

// Implements so that image.NewUniform can be used
func (rgb RGB) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(rgb.R) << 8, uint32(rgb.G) << 8, uint32(rgb.B) << 8, 255 << 8
}

type Coord struct {
	X, Y int
}

type AsciiConfig struct {
	CharSet      int
	FontSize     float64
	FontBytes    []byte
	Interpolate  bool
	InterpWeight float64
	InterpMemory map[Coord]float64
}

func NewAsciiConfig() *AsciiConfig {
	return &AsciiConfig{
		CharSet:      CHAR_SET_LIMITED,
		FontSize:     14,
		FontBytes:    nil,
		Interpolate:  true,
		InterpWeight: 0.6,
		InterpMemory: make(map[Coord]float64),
	}
}

func pntToDPI(pnt float64) float64 {
	return (6.756756756756757 * pnt) - 7.432432432432421
}

func drawAsciiChar(img *image.RGBA, x, y int, char string, c *freetype.Context, clr RGB) error {
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))

	pt := freetype.Pt(x, y)
	if _, err := c.DrawString(char, pt); err != nil {
		return err
	}
	return nil
}

func (ac *AsciiConfig) brightnessToAscii(b uint8) string {
	if ac.CharSet == CHAR_SET_BLOCK {
		return "█"
	}

	var ascii string
	if ac.CharSet == CHAR_SET_EXTENDED {
		ascii = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	} else if ac.CharSet == CHAR_SET_LIMITED {
		ascii = " .:-=+*#%@"
	}

	index := int(float32(b) / 255 * float32(len(ascii)-1))
	return string(ascii[index])
}

func (ac *AsciiConfig) GenerateAsciiImage(width, height int, getColour func(x, y int) RGB) (image.Image, error) {
	// Read the font data
	bounds := image.Rect(0, 0, width, height)

	f, err := truetype.Parse(ac.FontBytes)
	if err != nil {
		return nil, err
	}

	// Create the font context
	c := freetype.NewContext()
	c.SetDPI(pntToDPI(ac.FontSize))
	c.SetFont(f)
	c.SetClip(bounds.Bounds())

	// Get the pixel width and height of the font
	opts := truetype.Options{Size: ac.FontSize}
	face := truetype.NewFace(f, &opts)

	fontWidthPixelFixed, _ := face.GlyphAdvance(rune('█'))
	fontWidthPixel := fontWidthPixelFixed.Ceil()
	fontHeightPixel := face.Metrics().Ascent.Round()

	// Create a new image to hold the ascii characters
	ascii_img := image.NewRGBA(bounds)
	draw.Draw(ascii_img, ascii_img.Bounds(), image.Black, image.Point{}, draw.Over)

	// Draw the new image
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			if x%fontWidthPixel == 0 && y%fontHeightPixel == 0 {
				// Get the colour
				clr := getColour(x, y)

				// Get a brightness value for the image
				brightness := 0.299*float64(clr.R) + 0.587*float64(clr.G) + 0.114*float64(clr.B)

				// Interpolate the value if interpolation is turned on
				var interpolatedBrightness = brightness
				if ac.Interpolate {
					// If interpolation memory exists for this pixel then interpolate the brightness
					if oldBrightness, found := ac.InterpMemory[Coord{x, y}]; found {
						interpolatedBrightness = (float64(brightness) * ac.InterpWeight) + (float64(oldBrightness) * (1 - ac.InterpWeight))
					}

					// Store the brightness value in memory
					ac.InterpMemory[Coord{x, y}] = interpolatedBrightness
				}

				// Get the ascii string for the corresponding brightness value
				err = drawAsciiChar(ascii_img, x, y, ac.brightnessToAscii(uint8(interpolatedBrightness)), c, clr)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return ascii_img, nil
}
