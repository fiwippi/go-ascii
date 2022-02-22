package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/theckman/yacspin"

	"github.com/fiwippi/go-ascii"
	"video/internal/parse"
)

func Convert(ctx context.Context, conf ascii.Config, src, dst string, args ...string) error {
	imgD, ffDuration, ffProgress, errD := decode(ctx, src)
	imgE, errE := encode(ctx, dst, args...)

	// Handle the output of user progress
	s, err := createSpinner()
	if err != nil {
		return err
	}
	defer s.Stop()

	go func() {
		err = s.Start()
		if err != nil {
			log.Fatalln(err)
		}

		d := <-ffDuration
		for p := range ffProgress {
			s.Message(fmt.Sprintf("%02.2f%%", (float32(p)/float32(d))*100))
		}
	}()

	// Handle the decoding/encoding
	mem := &ascii.Memory{}
	for img := range imgD {
		asciiImg, err := ascii.Convert(img, conf, mem)
		if err != nil {
			close(imgE)
			return err
		}

		imgE <- asciiImg
		err = <-errE
		if err != nil {
			close(imgE)
			return err
		}
	}
	close(imgE)

	select {
	case err := <-errD:
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
	case err := <-errE:
		if err != nil {
			return fmt.Errorf("encode: %w", err)
		}
	}
	return nil
}

func encode(ctx context.Context, path string, args ...string) (chan<- image.Image, <-chan error) {
	// Make the channels
	errC := make(chan error, 1)
	imgC := make(chan image.Image)

	// Process the extra args
	var cmdArgs []string
	cmdArgs = append(cmdArgs,
		"-hide_banner", "-loglevel", "info",
		"-f", "image2pipe", "-c:v", "png", "-i", "-",
		"-y", "-an", "-pix_fmt", "yuv420p",
	)
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, path)

	// Filter out empty args
	n := 0
	for _, val := range cmdArgs {
		if val != "" {
			cmdArgs[n] = val
			n++
		}
	}
	cmdArgs = cmdArgs[:n]

	// Create the command
	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	var e bytes.Buffer
	cmd.Stderr = &e

	go func() {
		defer close(errC)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			errC <- fmtCmdErr(err, strings.TrimRight(e.String(), "\n"))
			return
		}
		defer stdin.Close()

		err = cmd.Start()
		if err != nil {
			errC <- fmtCmdErr(err, strings.TrimRight(e.String(), "\n"))
			return
		}

		for img := range imgC {
			err = png.Encode(stdin, img)
			if err != nil {
				errC <- fmtCmdErr(err, strings.TrimRight(e.String(), "\n"))
				return
			}
			errC <- nil // Send an acknowledgement the image encoded successfully
		}
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			errC <- fmtCmdErr(err, strings.TrimRight(e.String(), "\n"))
			return
		}
		errC <- nil
	}()

	return imgC, errC
}

func decode(ctx context.Context, path string) (<-chan image.Image, <-chan time.Duration, <-chan time.Duration, <-chan error) {
	// Make the channels
	errC := make(chan error, 1)
	imgC := make(chan image.Image)
	timeDur := make(chan time.Duration, 1)
	timeProg := make(chan time.Duration, 1)

	// Create the command
	cmd := exec.CommandContext(ctx,
		"ffmpeg", "-i", path,
		"-hide_banner", "-loglevel", "info",
		"-vcodec", "png", "-f", "image2pipe", "-",
	)

	go func() {
		defer close(errC)
		defer close(imgC)
		defer close(timeDur)
		defer close(timeProg)

		// Parse stderr for the duration of the track
		stderr, err := cmd.StderrPipe()
		if err != nil {
			errC <- err
			return
		}
		defer stderr.Close()

		var lastErr string
		go func() {
			sc := parse.FFScanner(stderr)

			for sc.Scan() {
				if !strings.Contains(sc.Text(), "time=") {
					lastErr = sc.Text()

					t := parse.ScanDuration(sc.Text())
					if t != 0 {
						timeDur <- t
					}
					continue
				}

				t := parse.ScanTime(sc.Text())
				if t != 0 {
					timeProg <- t
				}
			}
		}()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errC <- fmtCmdErr(err, lastErr)
			return
		}
		defer stdout.Close()

		err = cmd.Start()
		if err != nil {
			errC <- fmtCmdErr(err, lastErr)
			return
		}

		for {
			img, err := png.Decode(stdout)
			if err != nil {
				// Treat EOFs as end of the pipe so just break
				// and finish the command
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				errC <- fmtCmdErr(err, lastErr)
				return
			}

			imgC <- img
		}

		err = cmd.Wait()
		if err != nil {
			errC <- fmtCmdErr(err, lastErr)
			return
		}
	}()

	return imgC, timeDur, timeProg, errC
}

func fmtCmdErr(err error, s string) error {
	return fmt.Errorf("%w: %s", err, strings.TrimRight(s, "\n"))
}

func createSpinner() (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " Processing",
		SuffixAutoColon: true,
		Message:         "0%",
		StopMessage:     "Done",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	return yacspin.New(cfg)
}
