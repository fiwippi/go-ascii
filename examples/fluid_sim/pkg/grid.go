package fluids

import (
	"errors"
	"fmt"
	"github.com/fiwippi/go-ascii/pkg/ascii"
	"github.com/fiwippi/go-ascii/pkg/images"
	"github.com/schollz/progressbar/v3"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
)

// Grid composed of cells of fluid
type Grid struct {
	fluids        [][]fluid // 2D slice representing the cell data
	width, height int       // Width and height of the grid in number of cells,
	// each cell is one pixel when reading creating the
	// grid from an image
}

// Returns a grid from an image
func CreateGrid(img image.Image) Grid {
	// Get the image bounds to make the Grid object
	bounds := img.Bounds()
	g := Grid{width: bounds.Max.X, height: bounds.Max.Y}

	// Make the empty Grid
	g.fluids = make([][]fluid, g.width)
	for i := range g.fluids {
		g.fluids[i] = make([]fluid, g.height)
	}

	// Populate the Grid with the image values
	for y := bounds.Min.Y; y < g.height; y++ {
		for x := bounds.Min.X; x < g.width; x++ {
			red, green, blue, _ := img.At(x, y).RGBA()
			red, green, blue = red>>8, green>>8, blue>>8

			//  If colour is present in the image then fill the grid cell with fluid
			if red+blue+green > 0 {
				g.fluids[x][y] = initialFluid
			} else {
				g.fluids[x][y] = 0
			}
		}
	}

	return g
}

// Simulates the fluid falling in the grid and creates a video from all the frames created
func (g Grid) Simulate(ac *ascii.AsciiConfig, outputPath string, speed int, randClr bool) error {
	var err error

	// Determines whether to use random colours or a pattern
	if randClr {
		randomColours = true
	} else {
		randomColours = false
	}

	// Ensures the speed is in the correct range, this affects how fast the fluid "falls"
	s := float64(speed)
	if s < 1 || s > 9 {
		return errors.New("Speed must be in range 1-9")
	}

	// Creates the temp directory to save each frame to
	var framesDir string
	if framesDir, err = ioutil.TempDir("", "fluidSimFrames"); err != nil {
		return errors.New("Cannot create frames dir")
	}
	framesDir = filepath.ToSlash(framesDir) + "/"
	// Defers a function to ensure all frames and the directory are deleted once the function ends
	defer func() {
		err := os.RemoveAll(framesDir)
		if err != nil {
			log.Fatal("Error when deleting frame dir", err)
		}
	}()

	// Creates the spinner used to output user progress
	progressbar.OptionSpinnerType(43)(spinner)
	spinner.Describe("Running simulation...")
	counter := 1

	// Begins simulating each cell
	for {
		var changed = false
		spinner.Add(1) // Updates the spinner

		// Iterate from the bottom up (so no need to read from copy buffer),
		// this is because the cells below are updated before the cells above
		// so dont have to remember a cell's previous state to see if fluid
		// can flow down into it
		for y := g.height - 2; y >= 0; y-- {
			for x := g.width - 2; x >= 0; x-- {
				// Skip pixels with no fluid
				if g.fluids[x][y] == 0 {
					continue
				}

				// How much fluid can travel to the cell below and how much can that cell receive
				availableFluid := fluid(math.Ceil(float64(g.fluids[x][y]) * s / 10))
				if g.fluids[x][y] == 1 {
					availableFluid = 1
				}
				spaceToFill := fluidMax - g.fluids[x][y+1]

				// Calculate how much fluid can be put into the new cell and how much should be leftover
				leftoverFluid := max(availableFluid-spaceToFill, 0)
				filledFluid := availableFluid - leftoverFluid

				// Perform fluid operations for the cell below
				g.fluids[x][y+1] += filledFluid
				g.fluids[x][y] -= filledFluid

				// Notify the Grid state has chained so do not exit the simulation
				if filledFluid > 0 {
					changed = true
				}

				// If there is any leftover fluid, check if we can fill pixels to the bottom left or right
				if leftoverFluid > 0 {
					// How much space do the cells to the bottom left and right have?
					var spaceToFillLeft, spaceToFillRight fluid = 0, 0
					var rightExists, leftExists bool // Used to avoid checking for left and right  later
					//  on if they dont exist i.e. you're at a wall
					if x < g.width-1 {
						spaceToFillRight = fluidMax - g.fluids[x+1][y+1]
						rightExists = true
					}
					if x > 0 {
						spaceToFillLeft = fluidMax - g.fluids[x-1][y+1]
						leftExists = true
					}

					// If we can fill any of the cells
					if spaceToFillLeft > 0 || spaceToFillRight > 0 {
						// If the height of both cells is unequal, fill the cell
						// with the lowest value so that is is equal to the other one
						if spaceToFillLeft != spaceToFillRight {
							// Check if the left side needs to be equalised
							if spaceToFillLeft > spaceToFillRight && spaceToFillLeft > 0 {
								var fillNeeded fluid
								if !rightExists {
									fillNeeded = spaceToFillLeft
								} else {
									fillNeeded = g.fluids[x+1][y+1] - g.fluids[x-1][y+1]
								}
								filled := min(fillNeeded, leftoverFluid)

								// Add fluid to the bottom left cell
								g.fluids[x-1][y+1] += filled

								// Remove fluid from the original cell
								leftoverFluid -= filled
								g.fluids[x][y] -= filled
							}

							// Check if the right side needs to be equalised
							if spaceToFillRight > spaceToFillLeft && spaceToFillRight > 0 {
								var fillNeeded fluid
								if !leftExists {
									fillNeeded = spaceToFillRight
								} else {
									fillNeeded = g.fluids[x-1][y+1] - g.fluids[x+1][y+1]
								}
								filled := min(fillNeeded, leftoverFluid)

								// Add fluid to the bottom left cell
								g.fluids[x+1][y+1] += filled

								// Remove fluid from the original cell
								leftoverFluid -= filled
								g.fluids[x][y] -= filled
							}
						}

						// If there is any leftover fluid then split it in two and assign
						// it to each side as much as possible
						toLeft := leftoverFluid / 2
						toRight := leftoverFluid - toLeft

						// Assign as much as possible to the left side
						if leftExists {
							leftSpace := fluidMax - g.fluids[x-1][y+1]
							leftFill := min(leftSpace, toLeft)
							g.fluids[x-1][y+1] += leftFill
							g.fluids[x][y] -= leftFill
						}

						// Assign as much as possible to the right side
						if rightExists {
							rightSpace := fluidMax - g.fluids[x+1][y+1]
							rightFill := min(rightSpace, toRight)
							g.fluids[x+1][y+1] += rightFill
							g.fluids[x][y] -= rightFill
						}
					}
				}
			}
		}

		// Once the frame is done processing then save it to the frame dir
		name := fmt.Sprintf("%s%v.png", framesDir, counter)
		err = g.saveGridImage(ac, name)
		if err != nil {
			return errors.New("Cannot save the grid image: " + err.Error())
		}
		counter += 1

		// If no frame updating happened then the simulation is done
		if !changed {
			// Simulation done so now we convert all the frames to a video and save them to the file
			err = mergeFrames(framesDir, outputPath)
			if err != nil {
				return errors.New("Error merging frames: " + err.Error())
			}
			return nil
		}
	}
}

// Function to wrap creating and saving the grid image into the frames directory
func (g Grid) saveGridImage(ac *ascii.AsciiConfig, name string) error {
	// Create the image
	ascii_img, err := ac.GenerateAsciiImage(g.width, g.height, g.getColour)
	if err != nil {
		return err
	}

	// Save the image
	err = images.SaveImage(name, ascii_img, images.BestSpeed)
	if err != nil {
		return err
	}
	return nil
}
