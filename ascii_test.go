package ascii

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

var testImg image.Image
var conf = DefaultConfig()

func saveImg(img image.Image, name string) error {
	newF, err := os.Create(name)
	if err != nil {
		return err
	}
	defer newF.Close()

	err = jpeg.Encode(newF, img, nil)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	f, err := os.Open("test/fish.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}

	testImg = img
	os.Exit(m.Run())
}

func TestSingle(t *testing.T) {
	conf.FontSize = 20
	ascii, err := Convert(testImg, conf, nil)
	if err != nil {
		t.Error(err)
	}
	err = saveImg(ascii, "single.jpg")
	if err != nil {
		t.Error(err)
	}
}

func TestMultiple(t *testing.T) {
	// Create the test images
	bounds := testImg.Bounds()
	white := image.NewRGBA(bounds)
	black := image.NewRGBA(bounds)
	draw.Draw(white, bounds.Bounds(), image.White, image.Point{}, draw.Over)
	draw.Draw(black, bounds.Bounds(), image.Black, image.Point{}, draw.Over)

	// Perform the conversion
	conf.FontSize = 20
	mem := &Memory{}

	// What white will look like un-interpolated
	ascii, err := Convert(white, conf, nil)
	if err != nil {
		t.Error(err)
	}
	err = saveImg(ascii, "multiple-a.jpg")
	if err != nil {
		t.Error(err)
	}

	// What white will look like interpolated from black to white
	ascii, err = Convert(black, conf, mem)
	if err != nil {
		t.Error(err)
	}
	ascii, err = Convert(white, conf, mem)
	if err != nil {
		t.Error(err)
	}
	err = saveImg(ascii, "multiple-b.jpg")
	if err != nil {
		t.Error(err)
	}
}

func TestInterpolation(t *testing.T) {
	invalidConf := DefaultConfig()
	invalidConf.InterpolateWeight = -1
	_, err := Convert(testImg, invalidConf, &Memory{})
	if err == nil {
		t.Error(fmt.Errorf("error not returned for invalid interpolation"))
	}
	invalidConf.InterpolateWeight = 1.1
	_, err = Convert(testImg, invalidConf, &Memory{})
	if err == nil {
		t.Error(fmt.Errorf("error not returned for invalid interpolation"))
	}
}
