package fluids

import (
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Merges all the frames found in frameDir and combines them
// into a video with name vidName
func mergeFrames(frameDir, vidName string) error {
	// Create the ffmpeg command
	command := []string{"-y", "-framerate", "100", "-i", fmt.Sprintf("%s%%d.png", frameDir), "-c:v", "libx264", "-c:a", "copy", "-vf", "eq=brightness=0.06:saturation=2", "-r", "100", "-pix_fmt", "yuv420p", vidName}
	cmd := exec.Command("ffmpeg", command...)

	// Start process
	err := cmd.Start()
	if err != nil {
		return errors.New("Cannot start the process: " + err.Error())
	}

	// Wait for the process to finish
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Notify the user the spinner is being created
	spinner.Describe("Creating video...")

	// Run the spinner until the video is created and then exit
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-done:
			_ = spinner.Finish()
			return nil
		case _ = <-ticker.C:
			spinner.Add(1)
		}
	}
}
