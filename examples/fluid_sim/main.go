//go:generate pkger

package main

import (
	"bufio"
	"flag"
	fluids "fluids/pkg"
	"fmt"
	"github.com/fiwippi/ascii-image-creator/pkg/ascii"
	"github.com/fiwippi/ascii-image-creator/pkg/images"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Variables for the runtime flags
	var speed int
	var fontSize float64
	var overwrite, randClrs bool
	var inputPath, outputPath string

	// Parse in runtime flags
	flag.StringVar(&inputPath, "input", "", "Path to the image you want to make fluid (all non-black pixels are treated as fluid)")
	flag.StringVar(&outputPath, "output", "", "Name of the output video you want to make e.g. 'test.mp4'")
	flag.BoolVar(&overwrite, "overwrite", false, "Whether to automatically overwrite the output file if one already exists without prompting")
	flag.BoolVar(&randClrs, "random-colours", true, "Whether each fluid should be a random colour, if not a predetermined pattern is used instead")
	flag.IntVar(&speed, "speed", 4, "How quickly the fluid should fall, in the range 1-9,  1 is slowest, 9 is fastest")
	flag.Float64Var(&fontSize, "fontsize", 10, "Font size in points (NOT pixels)")
	flag.Parse()

	// Verifies speed is valid
	if speed < 1 || speed > 9 {
		log.Fatal("Speed must be in range 1-9")
	}

	// Verifies the flags have been filled in
	if !fluids.FileExists(inputPath) {
		log.Fatal("Input does not exist")
	}
	if outputPath == "" {
		log.Fatal("No output path specified")
	}
	if fluids.FileExists(outputPath) && !overwrite {
		// Ask user if they want to overwrite the file which already exists
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Would you like to overwrite the file? (y/N): ")
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			return
		}
	}

	// Load the default font
	var fontBytes []byte
	f, err := pkger.Open("/assets/CascadiaMono-Bold.ttf")
	if err != nil {
		log.Fatal("Error loading font: ", err)
	}
	fontBytes, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Error reading font data: ", err)
	}

	// Create the ascii config
	ac := ascii.NewAsciiConfig()
	ac.CharSet = ascii.CHAR_SET_EXTENDED
	ac.CustomCharSet = "MWKXY#"
	ac.FontSize = fontSize
	ac.FontBytes = fontBytes

	// Read in the input image
	img, err := images.ReadImage(inputPath)
	if err != nil {
		log.Fatal(img, err)
	}

	// Convert the image to the fluid grid
	g := fluids.CreateGrid(img)

	// Simulate the fluid
	err = g.Simulate(ac, outputPath, speed, randClrs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DONE!")
}
