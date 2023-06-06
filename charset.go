package ascii

// Charset is the set of characters which go-ascii
// uses to convert a pixel into an ascii character.
//
// The leftmost character is used to render the darkest
// pixel (black) and the rightmost is used for the
// lightest pixel (white). As the brightness the index
// of the character used increases to the right
type Charset string

const (
	CharsetLimited  Charset = " .:-=+*#%@"
	CharsetExtended         = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	CharsetBlock            = "█"
)
