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
)

type RGB struct {
	R, G, B uint8
}

// Implements so that image.NewUniform can be used
func (rgb RGB) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(rgb.R) << 8, uint32(rgb.G) << 8, uint32(rgb.B) << 8, 255 << 8
}

type AsciiConfig struct {
	CharSet   int
	FontSize  float64
	FontBytes []byte
}

func brightnessToAscii(b uint8, ac *AsciiConfig) string {
	var ascii string
	if ac.CharSet == CHAR_SET_EXTENDED {
		ascii = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	} else if ac.CharSet == CHAR_SET_LIMITED {
		ascii = " .:-=+*#%@"
	}

	index := int(float32(b) / 255 * float32(len(ascii)-1))
	return string(ascii[index])
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

	fontWidthPixelFixed, _ := face.GlyphAdvance(rune('B'))
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
				brightness := uint8(0.299*float64(clr.R) + 0.587*float64(clr.G) + 0.114*float64(clr.B))

				// Get the ascii string for the corresponding brightness value
				err = drawAsciiChar(ascii_img, x, y, brightnessToAscii(brightness, ac), c, clr)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return ascii_img, nil
}
