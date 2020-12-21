//go:generate pkger

package main

import (
	"flag"
	"fmt"
	"github.com/markbates/pkger"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/images"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// TODO add readme for both lib and this example
// TODO env sample e.g. /path/to/ffmpeg and ignore

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	// Parse in runtime flags
	var inputFlag = flag.String("input", "", "Path to the image you want to make ascii. Must be jpeg or png")
	var outputFlag = flag.String("output", "", "Name of the output image you want to make e.g. 'test.jpg'. Can be jpeg or png")
	var charset = flag.String("charset", "limited", "Type of charset you want to use, 'limited' or 'extended'")
	var fontPath = flag.String("font", "", "Path to a .ttf font file which the characters will be rendered as. If empty, 'inconsolata bold' is used")
	var fontSize = flag.Float64("fontsize", 14, "Font size in points (NOT pixels)")
	flag.Parse()

	// Verifies the flags have been filled in
	var imgPath, outputPath string
	if imgPath = *inputFlag; !fileExists(imgPath) {
		log.Fatal("Input video does not exist")
	}
	if outputPath = *outputFlag; outputPath == "" {
		log.Fatal("No output path specified")
	}

	// Parses the charset
	var cs int = ascii.CHAR_SET_LIMITED
	if strings.ToLower(*charset) == "extended" {
		cs = ascii.CHAR_SET_EXTENDED
	}

	// Get the font file as bytes and reading its data
	var fontBytes []byte
	var err error
	if *fontPath == "" {
		f, err := pkger.Open("/assets/Inconsolata-Bold.ttf")
		if err != nil {
			log.Fatal("Error loading font: ", err)
		}
		fontBytes, err = ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Error reading font data: ", err)
		}
	} else {
		fmt.Println(os.Getwd())
		fontBytes, err = ioutil.ReadFile(*fontPath)
		if err != nil {
			log.Fatal("Error loading/reading font data: ", err)
		}
	}

	// Set up the ascii config
	ac := &ascii.AsciiConfig{
		CharSet:   cs,
		FontBytes: fontBytes,
		FontSize:  *fontSize,
	}

	// Processes the video frames into ascii frames
	img, err := images.ReadImage(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	ascii_img, err := ascii.GenerateAsciiImage(img, ac)
	if err != nil {
		log.Fatal(err)
	}

	err = images.SaveImage(outputPath, ascii_img)
	if err != nil {
		log.Fatal(err)
	}

	// TODO output settings of conversion
	fmt.Println("Conversion Done")
}
