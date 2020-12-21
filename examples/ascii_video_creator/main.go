//go:generate pkger

package main

import (
	"ascii_video/convertor"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/markbates/pkger"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO build all executables in docker

// go build -o converter.exe && converter.exe --input davidguettafloyd.mp4 --output test.mkv
// go build -o converter.exe && converter.exe --input merged_output.mkv --output akira_v2.mkv
// go build -o converter.exe && converter.exe --input volttackle.mkv --output pokemon_inconsolata.mkv

// go generate && go build -o converter.exe && converter.exe --input volttackle-short.mkv --output pokemon_inconsolata_shhort.mkv --charset extended --fontsize 40
// go build -o converter.exe && converter.exe --input merged_output.mkv --output akira_big.mkv --charset extended --fontsize 30
func main() {
	// Parse in runtime flags
	var inputFlag = flag.String("input", "", "Path to the video you want to make ascii")
	var outputFlag = flag.String("output", "", "Name of the output video you want to make e.g. 'test.mkv'")
	var charset = flag.String("charset", "limited", "Type of charset you want to use, 'limited' or 'extended'. Default is 'limited'")
	var fontPath = flag.String("font", "", "Path to a .ttf font file which the characters will be rendered as")
	var fontSize = flag.Float64("fontsize", 14, "Font size in points (NOT pixels). Default is 14pt")
	flag.Parse()

	// Verifies the flags have been filled in
	var vidPath, outputPath string
	if vidPath = *inputFlag; !convertor.FileExists(vidPath) {
		log.Fatal("Input video does not exist")
	}
	if outputPath = *outputFlag; outputPath == "" {
		log.Fatal("No output path specified")
	}

	// Gets the absolute path to the executable to load the env files
	execPath, err := os.Executable()
	execPathSegments := strings.Split(filepath.ToSlash(execPath), "/")
	execDirPath := strings.Join(execPathSegments[:len(execPathSegments)-1], "/") + "/"

	// Load in the ENV values
	err = godotenv.Load(execDirPath + ".env")
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
	}

	convertor.FFmpegPath = os.Getenv("FFMPEG_PATH")
	convertor.FFprobePath = os.Getenv("FFPROBE_PATH")

	// Parses the charset
	var cs int = ascii.CHAR_SET_LIMITED
	if strings.ToLower(*charset) == "extended" {
		cs = ascii.CHAR_SET_EXTENDED
	}

	// Get the font file as bytes and reading its data
	var fontBytes []byte
	if *fontPath == "" {
		f, err := pkger.Open("/convertor/Inconsolata-Bold.ttf")
		if err != nil {
			log.Fatal("Error loading font: ",err)
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
		CharSet:  cs,
		FontBytes: fontBytes,
		FontSize: *fontSize,
	}

	// Creating the input and output temp directories
	var framesDir, asciiFramesDir string
	if framesDir, err = ioutil.TempDir("", "frames"); err != nil {
		log.Fatal("Cannot create frames dir: ", err)
	}
	framesDir = filepath.ToSlash(framesDir) + "/"
	if asciiFramesDir, err = ioutil.TempDir("", "frames_ascii"); err != nil {
		log.Fatal("Cannot create ascii frames dir: ", err)
	}
	asciiFramesDir = filepath.ToSlash(asciiFramesDir) + "/"
	defer os.RemoveAll(framesDir); defer os.RemoveAll(asciiFramesDir)

	// Exports the input video into its seperate frames in the temporary frames folder
	fmt.Println("Exporting video to frames:")
	err = convertor.VideoToFrames(vidPath, framesDir)
	if err != nil {
		log.Fatal("Error encountered while exporting the video into its frames: ", err)
	}

	// Processes the video frames into ascii frames
	fmt.Println("Processing frames:")
	err = convertor.CreateAscii(framesDir, asciiFramesDir, ac)
	if err != nil {
		log.Fatal("Error encountered while processing the frames into ascii frames: ", err)
	}

	// Merges the ascii frames into a final video output
	fmt.Println("Merging frames into final video:")
	err = convertor.MergeFrames(asciiFramesDir, vidPath, outputPath)
	if err != nil {
		log.Fatal("Error encountered while merging the ascii frames into final video output: ", err)
	}
}