package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Create Flags
	src := flag.String("i", "", "path to video to convert to ascii")
	stringArgs := flag.String("args", "", "specify extra args for ffmpeg")
	fontsize := flag.Float64("fontsize", 14, "fontsize of the ascii characters")
	overwrite := flag.Bool("y", false, "automatically overwrites the output file if it exists")

	// Parse flags
	flag.Usage = func() {
		fmt.Printf("Usage: ./video -i in.mp4 out.mp4\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var output string
	if flag.Parse(); len(flag.Args()) > 0 {
		output = flag.Args()[0]
	} else {
		flag.Usage()
		os.Exit(1)
	}

	if *src == "" {
		fmt.Println("Input file not specified!")
		os.Exit(1)
	}

	// Check if overwrite
	if exists(*src) && !*overwrite {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Would you like to overwrite the file? (y/N): ")
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			fmt.Println("File already exists!")
			os.Exit(1)
		}
	}

	// Perform the conversion
	args := strings.Split(*stringArgs, " ")
	err := Convert(context.Background(), *src, output, *fontsize, args...)
	if err != nil {
		log.Fatalln(err)
	}
}

func exists(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return false
	}
	return true
}
