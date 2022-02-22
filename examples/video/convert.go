package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"os/exec"
	"strings"

	"github.com/fiwippi/go-ascii"
)

func Convert(ctx context.Context, conf ascii.Config, src, dst string, args ...string) error {
	imgD, errD := decode(ctx, src)
	imgE, errE := encode(ctx, conf, dst, args...)

	for img := range imgD {
		imgE <- img
	}
	close(imgE)

	err := <-errD
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	err = <-errE
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

func encode(ctx context.Context, conf ascii.Config, path string, args ...string) (chan<- image.Image, <-chan error) {
	// Make the channels
	errC := make(chan error, 1)
	imgC := make(chan image.Image)

	// Process the extra args
	var cmdArgs []string
	cmdArgs = append(cmdArgs,
		"-hide_banner", "-loglevel", "error",
		"-f", "image2pipe", "-c:v", "png", "-i", "-",
		"-y",
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

		mem := &ascii.Memory{}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}
		defer stdin.Close()

		err = cmd.Start()
		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}

		for img := range imgC {
			asciiImg, err := ascii.Convert(img, conf, mem)
			if err != nil {
				errC <- fmtCmdErr(err, e)
				return
			}

			err = png.Encode(stdin, asciiImg)
			if err != nil {
				errC <- fmtCmdErr(err, e)
				return
			}
		}

		stdin.Close()
		err = cmd.Wait()

		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}
		errC <- nil
	}()

	return imgC, errC
}

func decode(ctx context.Context, path string) (<-chan image.Image, <-chan error) {
	// Make the channels
	errC := make(chan error, 1)
	imgC := make(chan image.Image)

	// Create the command
	cmd := exec.CommandContext(ctx,
		"ffmpeg", "-i", path,
		"-hide_banner", "-loglevel", "error",
		"-vcodec", "png", "-f", "image2pipe", "-",
	)
	var e bytes.Buffer
	cmd.Stderr = &e

	go func() {
		defer close(errC)
		defer close(imgC)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}
		defer stdout.Close()

		err = cmd.Start()
		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}

		for {
			img, err := png.Decode(stdout)
			if err != nil {
				// Treat EOFs as end of the pipe so nil error
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				errC <- err
				return
			}
			imgC <- img
		}

		err = cmd.Wait()
		if err != nil {
			errC <- fmtCmdErr(err, e)
			return
		}
	}()

	return imgC, errC
}

func fmtCmdErr(err error, b bytes.Buffer) error {
	return fmt.Errorf("%w: %s", err, strings.TrimRight(b.String(), "\n"))
}
