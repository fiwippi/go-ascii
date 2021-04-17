package fluids

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
)

// Returns a seed to get a consistent random colour from a set of coordinates
func getSeed(x, y int) int64 {
	// Convert the coordinates to uint32 used for
	// both the pattern and the random colour seeds
	a := uint64(x)
	b := uint64(y)

	// If pattern wanted, returns a seed which generates a diagonal pattern
	if !randomColours {
		return int64((a & 0xffff) | ((b & 0xffff) << 16))
	}

	// Create a byte array of the bits from a and b
	aBytes := make([]byte, 8)
	bBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(aBytes, a)
	binary.LittleEndian.PutUint64(bBytes, b)

	// Create the hash
	h := sha256.New()
	h.Write(aBytes)
	h.Write(bBytes)

	// Process the hash and return the seed
	var seed int64
	sum := h.Sum(nil)
	buffer := bytes.NewBuffer(sum)
	binary.Read(buffer, binary.LittleEndian, &seed)

	return seed
}
