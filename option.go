package ascii

import (
	"fmt"

	"golang.org/x/image/font/opentype"
)

type options struct {
	font    *opentype.Font
	charset Charset
	fontPts float64
	mem     *Memory
}

// Option is a function which is supplied to
// ConvertWithOpts and which mutates the
// settings which the convertor uses.
//
// Options include:
//   * CSet -> Character set
//   * FontPts -> Font size in pts
//   * Font -> Font
//   * Interpolate -> Interpolation of characters
type Option func(args *options) error

// CSet changes the character set that the convertor uses
func CSet(c Charset) Option {
	return func(args *options) error {
		args.charset = c
		return nil
	}
}

// FontPts changes the font size the convertor renders
// the characters in which is specified in points (pts)
// as opposed to pixels
func FontPts(pts float64) Option {
	return func(args *options) error {
		if pts <= 0 {
			return fmt.Errorf("font size cannot be smaller than 0")
		}
		args.fontPts = pts
		return nil
	}
}

// Font changes which font the convertor uses,
// fonts are expected to be monospace
func Font(f *opentype.Font) Option {
	return func(args *options) error {
		if f == nil {
			return fmt.Errorf("font cannot be nil")
		}
		args.font = f
		return nil
	}
}

// Interpolate is able to interpolate the character used
// if multiple ConvertWithOpts calls are used with this
// option and if a valid Memory struct is provided.
//
// This makes the change in characters less pronounced
// between successive images, this is useful if you are
// converting successive frames of a video.
func Interpolate(mem *Memory) Option {
	return func(args *options) error {
		if mem == nil {
			return fmt.Errorf("memory supplied is nil")
		}
		if mem.data == nil {
			mem.Reset()
		}
		args.mem = mem
		return nil
	}
}
