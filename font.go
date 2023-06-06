package ascii

import (
	_ "embed"
	"errors"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed CascadiaMono-Bold.ttf
var fontBytes []byte

var defaultFont *opentype.Font

func init() {
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		panic("could not parse default embedded font")
	}
	defaultFont = f
}

type parsedFont struct {
	face          font.Face
	height, width int
}

func parseFont(f *opentype.Font, pts float64) (parsedFont, error) {
	// Create the font face
	faceOpts := &opentype.FaceOptions{
		Size:    pts,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	face, err := opentype.NewFace(f, faceOpts)
	if err != nil {
		return parsedFont{}, err
	}

	// Process font face metrics
	glyphBounds, _, found := face.GlyphBounds('█')
	if !found {
		return parsedFont{}, errors.New("failed getting font face width")
	}

	return parsedFont{
		face:   face,
		height: face.Metrics().Ascent.Round(),
		width:  glyphBounds.Max.X.Round(),
	}, nil
}

func (pf parsedFont) drawString(s string, clr color.Color, dst *image.RGBA, x, y int) {
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(clr),
		Face: pf.face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}
