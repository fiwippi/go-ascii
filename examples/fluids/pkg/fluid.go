package fluids

const fluidMax = 12    // The max amount of fluid that can fill a cell
const initialFluid = 8 // The initial amount of fluid every cell starts with

type fluid int

// Max comparison for fluids
func max(a, b fluid) fluid {
	if a > b {
		return a
	}
	return b
}

// Min comparison for fluids
func min(a, b fluid) fluid {
	if a < b {
		return a
	}
	return b
}
