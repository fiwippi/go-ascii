package convertor

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/images"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ConvertVideo(vidPath, outputPath string, ac *ascii.AsciiConfig) error {
	// Creating the input and output temp directories
	var framesDir, asciiFramesDir string
	var err error

	if framesDir, err = ioutil.TempDir("", "frames"); err != nil {
		log.Fatal("Cannot create frames dir: ", err)
	}
	framesDir = filepath.ToSlash(framesDir) + "/"
	if Debug {
		log.Println("Frames dir:", framesDir)
	}
	if asciiFramesDir, err = ioutil.TempDir("", "frames_ascii"); err != nil {
		log.Fatal("Cannot create ascii frames dir: ", err)
	}
	asciiFramesDir = filepath.ToSlash(asciiFramesDir) + "/"
	if Debug {
		log.Println("Ascii Frames dir:", asciiFramesDir)
	}
	defer os.RemoveAll(framesDir)
	defer os.RemoveAll(asciiFramesDir)

	// Exports the input video into its separate frames in the temporary frames folder
	fmt.Println("Exporting video to frames:")
	err = videoToFrames(vidPath, framesDir)
	if err != nil {
		log.Fatal("Error encountered while exporting the video into its frames: ", err)
	}

	// Processes the video frames into ascii frames
	fmt.Println("Processing frames:")
	err = createAscii(framesDir, asciiFramesDir, ac)
	if err != nil {
		log.Fatal("Error encountered while processing the frames into ascii frames: ", err)
	}

	// Merges the ascii frames into a final video output
	fmt.Println("Merging frames into final video:")
	err = mergeFrames(asciiFramesDir, vidPath, outputPath)
	if err != nil {
		log.Fatal("Error encountered while merging the ascii frames into final video output: ", err)
	}

	return nil
}

func videoToFrames(vidPath, inputPath string) error {
	command := []string{"-i", vidPath, "-vsync", "0", fmt.Sprintf("%sout_%%04d.png", inputPath)}
	err := runCmd(command, "Duration")
	if err != nil {
		return err
	}
	return nil
}

func createAscii(inputDir, outputDir string, ac *ascii.AsciiConfig) error {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	bar := pb.StartNew(len(files))
	for _, f := range files {
		if f.Name() != ".gitkeep" {
			fp := inputDir + f.Name()

			img, err := images.ReadImage(fp)
			if err != nil {
				return err
			}

			generate := func(x, y int) ascii.RGB {
				r, g, b, _ := img.At(x, y).RGBA()
				r, g, b = r>>8, g>>8, b>>8 // Colours
				return ascii.RGB{uint8(r), uint8(g), uint8(b)}
			}

			width, height := img.Bounds().Max.X, img.Bounds().Max.Y
			ascii_img, err := ac.GenerateAsciiImage(width, height, generate)
			if err != nil {
				log.Fatal(err)
			}

			err = images.SaveImage(outputDir+f.Name(), ascii_img)
			if err != nil {
				return err
			}
			bar.Increment()
		}
	}
	bar.Finish()
	return nil
}

func mergeFrames(outputPath, vidPath, vidName string) error {
	if Debug {
		log.Println(outputPath, vidPath, vidName)
	}
	framerate, err := getFramerate(vidPath)
	if err != nil {
		return err
	}
	stringFramerate := fmt.Sprintf("%v", framerate)
	command := []string{"-y", "-framerate", stringFramerate, "-i", fmt.Sprintf("%sout_%%04d.png", outputPath), "-i", vidPath, "-c:v", "libx264", "-c:a", "copy", "-vf", "eq=brightness=0.06:saturation=2", "-map", "0:v:0", "-map", "1:a:0?", "-r", stringFramerate, "-pix_fmt", "yuv420p", vidName}
	err = runCmd(command, "Duration")
	if err != nil {
		return err
	}
	return nil
}
