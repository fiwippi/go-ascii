package ascii

import (
	_ "embed"

	"golang.org/x/image/font/opentype"
)

//go:embed CascadiaMono-Bold.ttf
var fontBytes []byte

var DefaultFont = func() *opentype.Font {
	font, err := opentype.Parse(fontBytes)
	if err != nil {
		panic("could not parse default embedded font")
	}
	return font
}()

// Config holds data used to render ascii images
type Config struct {
	// Which charset to use in order to render the image
	CharSet CharSet
	// The fontsize to draw in points (NOT pixels)
	FontSize float64
	// Font used to draw the characters on the image, the
	// font is expected to be monospace
	Font *opentype.Font
	// How strong should the interpolation be between 0-1:
	//   - 0 means no interpolation
	//   - 1 means the last character will always be used
	InterpolateWeight float64
	// Whether to draw on a transparent background, as opposed to a black one
	Transparency bool
}

// DefaultConfig creates a config with the default settings
func DefaultConfig() Config {
	return Config{
		CharSet:           CharSetExtended,
		FontSize:          14,
		Font:              DefaultFont,
		InterpolateWeight: 0.4,
		Transparency:      false,
	}
}
