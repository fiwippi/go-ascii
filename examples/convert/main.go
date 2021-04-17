package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"github.com/fiwippi/go-ascii"
	convert "github.com/fiwippi/go-ascii/examples/convert/pkg"
	"github.com/gabriel-vasile/mimetype"
	"log"
	"os"
	"strings"
)

//go:embed assets/CascadiaMono-Bold.ttf
var fontBytes []byte

func main() {
	// General variables
	var err error

	// Variables for the runtime flags
	var inputPath, outputPath, charset string
	var fontSize, intpf float64
	var overwrite, interpolate, transparency bool

	// Parse in runtime flags
	flag.StringVar(&inputPath, "input", "", "Path to the image/video you want to make ascii. Images must be jpeg or png")
	flag.StringVar(&outputPath, "output", "", "Name of the output image/video you want to make e.g. 'test.jpg'. Images may be jpeg or png")
	flag.StringVar(&charset, "charset", "limited", "Type of charset you want to use, 'limited' or 'extended' or 'block'")
	flag.Float64Var(&fontSize, "fontsize", 14, "Font size in points (NOT pixels)")
	flag.Float64Var(&intpf, "intpf", 0.6, "Interpolation factor to use, between 1 (none) to 0 (max) interpolation")
	flag.BoolVar(&overwrite, "overwrite", false, "Whether to automatically overwrite the output file if one already exists without prompting")
	flag.BoolVar(&convert.Debug, "verbose", false, "Prints verbose information")
	flag.BoolVar(&interpolate, "interpolate", true, "Whether to use interpolation for video rendering")
	flag.BoolVar(&transparency, "transparency", false, "Whether to use a transparent background for the ascii images")
	flag.StringVar(&convert.FFmpegPath, "ffmpeg", "ffmpeg", "Path to the ffmpeg binary")
	flag.StringVar(&convert.FFprobePath, "ffprobe", "ffprobe", "Path to the ffprobe binary")
	flag.Parse()

	// Verifies the flags have been filled in
	if !convert.FileExists(inputPath) {
		log.Fatal("Input does not exist")
	}
	if outputPath == "" {
		log.Fatal("No output path specified")
	}
	if convert.FileExists(outputPath) && !overwrite {
		// Ask user if they want to overwrite the file which already exists
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Would you like to overwrite the file? (y/N): ")
		scanner.Scan()
		if strings.TrimSpace(scanner.Text()) != "y" {
			return
		}
	}

	// Parses the charset
	var cs = ascii.CharSetLimited
	if strings.ToLower(charset) == "extended" {
		cs = ascii.CharSetExtended
	} else if strings.ToLower(charset) == "block" {
		cs = ascii.CharSetBlock
	}

	// Set up the ascii config
	ac := ascii.NewAsciiConfig()
	ac.CharSet = cs
	ac.FontSize = fontSize
	ac.FontBytes = fontBytes
	ac.InterpWeight = intpf
	ac.Interpolate = interpolate
	ac.Transparency = transparency

	// Detects whether an image or video conversion is needed
	mime, _ := mimetype.DetectFile(inputPath)
	if strings.HasPrefix(mime.String(), "image") {
		err = convert.ConvertImage(inputPath, outputPath, ac)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasPrefix(mime.String(), "video") {
		err = convert.ConvertVideo(inputPath, outputPath, ac)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Conversion Done")
}
