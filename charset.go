package ascii

type CharSet int

const (
	CharSetLimited  CharSet = iota // Use the limited character set: " .:-=+*#%@"
	CharSetExtended                // Use the extended character set: ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	CharSetBlock                   // Use the block character set: "█"
)

// parseBrightness returns an appropriate ascii string based on
// the brightness of a value in the range 0-254
func (cs CharSet) parseBrightness(b uint8) string {
	var ascii string
	switch cs {
	case CharSetLimited:
		ascii = " .:-=+*#%@"
	case CharSetExtended:
		ascii = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	case CharSetBlock:
		return "█"
	}

	index := min(int(float64(b)/254*float64(len(ascii))), len(ascii)-1)
	return string(ascii[index])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
