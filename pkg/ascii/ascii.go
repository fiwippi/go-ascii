package ascii

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	_ "github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	_ "image/draw"
	_ "io/ioutil"
)

// TODO sample the colour instead of using the first input

const (
	CHAR_SET_LIMITED = iota
	CHAR_SET_EXTENDED
)

type coord struct {
	x, y int
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

	index := int(float32(b)/255 * float32(len(ascii)-1))
	return string(ascii[index])
}

func pntToDPI(pnt float64) float64 {
	return (6.756756756756757 * pnt) -7.432432432432421
}

func drawAsciiChar(img *image.RGBA, x, y int, char string, c *freetype.Context, clr color.Color) error {
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))

	pt := freetype.Pt(x, y)
	if _, err := c.DrawString(char, pt); err != nil {
		return err
	}
	return nil
}

func GenerateAsciiImage(img image.Image, ac *AsciiConfig) (image.Image, error) {
	// Get the bounds for the image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Read the font data
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

	// Get the ascii strings for each pixel in the image
	asciiMap := make(map[coord]string)
	colourMap := make(map[coord]color.Color)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			if x % fontWidthPixel == 0 && y % fontHeightPixel == 0 {
				r, g, b, _ := img.At(x, y).RGBA()
				colourMap[coord{x, y}] = img.At(x, y)

				// Convert rgb values to be in range 0-255 so 8 bits for grayscale
				r, g, b = r >> 8, g >> 8, b >> 8

				// Get a brightness value for the image
				brightness := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

				// Get the ascii string for the corresponding brightness value
				asciiMap[coord{x, y}] = brightnessToAscii(brightness, ac)
			}
		}
	}

	// Create a new image to hold the ascii characters
	ascii_img := image.NewRGBA(img.Bounds())
	draw.Draw(ascii_img, ascii_img.Bounds(), image.Black, image.Point{}, draw.Over)

	// Write the ascii characters to the new image
	for k := range asciiMap {
		char := asciiMap[k]
		r, g, b, _ := colourMap[k].RGBA()
		r, g, b = r >> 8, g >> 8, b >> 8
		bi := 1.1 // Brightness increase
		r_f, g_f, b_f := float64(r) * bi, float64(g) * bi, float64(b) * bi
		if r_f > 255 {r = 255}
		if g_f > 255 {g = 255}
		if b_f > 255 {b = 255}
		clr := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

		//drawBasicAsciiChar(ascii_img, k.x, k.y, char, asciiFont, clr)
		err = drawAsciiChar(ascii_img, k.x, k.y, char, c, clr)
		if err != nil {
			return nil, err
		}
	}

	return ascii_img, nil
}
