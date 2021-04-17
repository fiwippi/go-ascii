// Some code based off of code by floostack/transcoder as follows:
//
//MIT License
//
//Copyright (c) 2020 FlooStack
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package convert

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	FFmpegPath  string = "ffmpeg"
	FFprobePath string = "ffmpeg"
	Debug              = false
)

// Converts a sexagesimal time (hh:mm:ss.ms) to seconds
func durToSec(t string) (float64, error) {
	times := strings.Split(t, ":")

	h, m := times[0], times[1]
	s_ms := strings.Split(times[2], ".")
	s, ms := s_ms[0], s_ms[1][:2]

	t_s, err := time.ParseDuration(fmt.Sprintf("%vh%vm%vs%vms", h, m, s, ms))
	if err != nil {
		return -1, err
	}

	return t_s.Seconds(), nil
}

// Strips a carriage return or newline from a string
func stripReturn(s string) string {
	if strings.HasSuffix(s, "\r\n") {
		return strings.Trim(s, "\r\n")
	}
	return strings.Trim(s, "\n")
}

// Gets the duration of a video
func getDuration(vidPath string) (float64, error) {
	var outb bytes.Buffer

	args := []string{"-i", vidPath, "-show_entries", "format=duration", "-v", "quiet", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal"}

	cmd := exec.Command(FFprobePath, args...)
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		return -1, err
	}

	t, err := durToSec(stripReturn(outb.String()))
	if err != nil {
		return -1, err
	}

	return t, nil
}

// Gets the framerate of a video
func getFramerate(vidPath string) (float64, error) {
	var outb bytes.Buffer

	args := []string{"-i", vidPath, "-v", "error", "-select_streams", "v", "-of", "default=noprint_wrappers=1:nokey=1", "-show_entries", "stream=r_frame_rate"}

	cmd := exec.Command(FFprobePath, args...)
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		return -1, err
	}

	nums := strings.Split(stripReturn(outb.String()), "/")
	var n1, n2 float64
	if n1, err = strconv.ParseFloat(nums[0], 64); err != nil {
		return -1, err
	}
	if n2, err = strconv.ParseFloat(nums[1], 64); err != nil {
		return -1, err
	}
	return n1 / n2, nil
}

// Sends the process of an ffmpeg process to a channel
func progress(pType string, total float64, out chan int64, stderrIn io.ReadCloser) {
	defer stderrIn.Close()

	scanner := bufio.NewScanner(stderrIn)

	split := func(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		if i := bytes.IndexByte(data, '\r'); i >= 0 {
			// We have a cr terminated line
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}
	scanner.Split(split)

	for scanner.Scan() {
		line := scanner.Text()
		if Debug {
			log.Println(line)
		}
		var re = regexp.MustCompile(`=\s+`)
		st := re.ReplaceAllString(line, `=`)
		if pType == "Duration" && strings.Contains(st, "time=") {
			currentTime := strings.Split(strings.Split(st, "time=")[1], " ")[0]
			timesec, err := durToSec(currentTime)
			if err != nil {
				log.Fatal(err)
			}
			counter := (timesec * 100) / total
			out <- int64(counter)
		}
		if pType == "Frames" && strings.Contains(st, "frame=") {
			currentFrame := strings.Split(strings.Split(st, "frame=")[1], " ")[0]
			if currentFrame == "" {
				out <- -1
				return
			}
			frameInt, err := strconv.Atoi(currentFrame)
			if err != nil {
				log.Fatal(err)
			}
			counter := (float64(frameInt) * 100) / total
			out <- int64(counter)
		}
	}
	out <- -1
}

// Runs an ffmpeg command
func runCmd(command []string, pType string) error {
	cmd := exec.Command(FFmpegPath, command...)

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start process
	err = cmd.Start()
	if err != nil {
		return err
	}

	var total float64
	var max int64
	var ffmpegInput string
	out := make(chan int64)
	if pType == "Duration" {
		inputs := strings.Split(strings.Join(command, " "), "-i ")
		if len(inputs) > 2 {
			ffmpegInput = strings.Split(inputs[2], " ")[0]
		} else {
			ffmpegInput = strings.Split(inputs[1], " ")[0]
		}
		total, err = getDuration(ffmpegInput)
		max = 100
		if err != nil {
			return err
		}
	} else if pType == "Frames" {
		outputDir := strings.Split(strings.Split(strings.Join(command, " "), "-i ")[1], "/")[0]
		files, err := ioutil.ReadDir(outputDir)
		if err != nil {
			return err
		}
		total = float64(len(files))
		max = int64(total)
	}
	go progress(pType, total, out, stderrIn)

	bar := pb.Start64(max)
	func() {
		for {
			select {
			case percent := <-out:
				if percent == -1 {
					bar.SetCurrent(max)
					close(out)
					return
				} else {
					bar.SetCurrent(percent)
				}
			}
		}
	}()
	bar.Finish()

	return nil
}
