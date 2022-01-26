package ascii

// coord represents the coordinates of a pixel
type coord struct {
	X, Y int
}

// Memory stores the brightness of pixels and can be successively passed
// to the Convert function so that the ascii character shape is interpolated,
// this is useful for example if you are converting successive frames in
// a video since it makes the change in characters look more natural
type Memory struct {
	data map[coord]float64
}

// Reset clears any data the memory holds so it can be used again as new
func (m *Memory) Reset() {
	m.data = make(map[coord]float64)
}
