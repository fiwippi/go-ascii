package convertor

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/ascii"
	"github.com/nadav-rahimi/ascii-image-creator/pkg/images"
	"io/ioutil"
)

func CreateAscii(inputDir, outputDir string, ac *ascii.AsciiConfig) error {
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

			ascii_img, err := ascii.GenerateAsciiImage(img, ac)
			if err != nil {
				return err
			}

			err = images.SaveImage(outputDir + f.Name(), ascii_img)
			if err != nil {
				return err
			}
			bar.Increment()
		}
	}
	bar.Finish()
	return nil
}

func VideoToFrames(vidPath, inputPath string) error {
	command := []string{"-i", vidPath, "-vsync", "0", fmt.Sprintf("%sout_%%04d.png", inputPath)}
	err := runCmd(command, "Duration")
	if err != nil {
		return err
	}
	return nil
}

func MergeFrames(outputPath, vidPath, vidName string) error {
	framerate, err := getFramerate(vidPath)
	if err != nil {
		return err
	}
	stringFramerate := fmt.Sprintf("%v", framerate)
	command := []string{"-y", "-framerate", stringFramerate, "-i", fmt.Sprintf("%sout_%%04d.png", outputPath), "-i", vidPath, "-c:v", "libx264", "-c:a", "copy", "-vf", "eq=brightness=0.06:saturation=2", "-map", "0:v:0", "-map", "1:a:0", "-r", stringFramerate, "-pix_fmt", "yuv420p", vidName}
	err = runCmd(command, "Duration")
	if err != nil {
		return err
	}
	return nil
}
