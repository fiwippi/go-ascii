//go:generate pkger

package main

import (
	"ascii_image/convertor"
	"bufio"
	"flag"
	"fmt"
	"github.com/fiwippi/go-ascii/pkg/ascii"
	"github.com/fiwippi/go-ascii/pkg/images"
	"github.com/gabriel-vasile/mimetype"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// General variables
	var err error

	// Variables for the runtime flags
	var inputPath, outputPath, charset, fontPath, optimise string
	var fontSize, intpf float64
	var overwrite, interpolate, transparency bool

	// Parse in runtime flags
	flag.StringVar(&inputPath, "input", "", "Path to the image/video you want to make ascii. Images must be jpeg or png")
	flag.StringVar(&outputPath, "output", "", "Name of the output image/video you want to make e.g. 'test.jpg'. Images may be jpeg or png")
	flag.StringVar(&charset, "charset", "limited", "Type of charset you want to use, 'limited' or 'extended' or 'block'")
	flag.StringVar(&fontPath, "font", "", "Path to a .ttf font file which the characters will be rendered as. If empty, 'Cascadia Code Mono Bold' is used")
	flag.StringVar(&optimise, "optimise", "balanced", "Whether to optimise 'speed', 'balanced' or 'size' for the file encoding")
	flag.Float64Var(&fontSize, "fontsize", 14, "Font size in points (NOT pixels)")
	flag.Float64Var(&intpf, "intpf", 0.6, "Interpolation factor to use, between 1 (none) to 0 (max) interpolation")
	flag.BoolVar(&overwrite, "overwrite", false, "Whether to automatically overwrite the output file if one already exists without prompting")
	flag.BoolVar(&convertor.Debug, "verbose", false, "Prints verbose information")
	flag.BoolVar(&interpolate, "interpolate", true, "Whether to use interpolation for video rendering")
	flag.BoolVar(&transparency, "transparency", false, "Whether to use a transparent background for the ascii images")
	flag.StringVar(&convertor.FFmpegPath, "ffmpeg", "ffmpeg", "Path to the ffmpeg binary")
	flag.StringVar(&convertor.FFprobePath, "ffprobe", "ffprobe", "Path to the ffprobe binary")
	flag.Parse()

	// Verifies the flags have been filled in
	if !convertor.FileExists(inputPath) {
		log.Fatal("Input does not exist")
	}
	if outputPath == "" {
		log.Fatal("No output path specified")
	}
	if convertor.FileExists(outputPath) && !overwrite {
		// Ask user if they want to overwrite the file which already exists
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Would you like to overwrite the file? (y/N): ")
		scanner.Scan()
		if strings.TrimSpace(scanner.Text()) != "y" {
			return
		}
	}

	// Parses the charset
	var cs int = ascii.CHAR_SET_LIMITED
	if strings.ToLower(charset) == "extended" {
		cs = ascii.CHAR_SET_EXTENDED
	} else if strings.ToLower(charset) == "block" {
		cs = ascii.CHAR_SET_BLOCK
	}

	// Get the font file as bytes and reading its data
	var fontBytes []byte
	if fontPath == "" {
		f, err := pkger.Open("/assets/CascadiaMono-Bold.ttf")
		if err != nil {
			log.Fatal("Error loading font: ", err)
		}
		fontBytes, err = ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Error reading font data: ", err)
		}
	} else {
		fmt.Println(os.Getwd())
		fontBytes, err = ioutil.ReadFile(fontPath)
		if err != nil {
			log.Fatal("Error loading/reading font data: ", err)
		}
	}

	// Set up the ascii config
	ac := ascii.NewAsciiConfig()
	ac.CharSet = cs
	ac.FontSize = fontSize
	ac.FontBytes = fontBytes
	ac.InterpWeight = intpf
	ac.Interpolate = interpolate
	ac.Transparency = transparency

	// Sets the encoding
	var encoding images.CompressionLevel
	if optimise == "balanced" {
		encoding = images.DefaultCompression
	} else if optimise == "speed" {
		encoding = images.BestSpeed
	} else if optimise == "size" {
		encoding = images.BestCompression
	}

	// Detects whether an image or video conversion is needed
	mime, _ := mimetype.DetectFile(inputPath)
	if strings.HasPrefix(mime.String(), "image") {
		err = convertor.ConvertImage(inputPath, outputPath, ac, encoding)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasPrefix(mime.String(), "video") {
		err = convertor.ConvertVideo(inputPath, outputPath, ac, encoding)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Conversion Done")
}
