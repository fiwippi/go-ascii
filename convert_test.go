package ascii

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

var testImg image.Image

func saveImg(img image.Image, name string) error {
	newF, err := os.Create("test/" + name)
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
	f, err := os.Open("test/jellyfish.jpg")
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

func TestConvert(t *testing.T) {
	ascii, err := Convert(testImg)
	require.Nil(t, err)
	assert.Nil(t, saveImg(ascii, "convert-default.jpg"))
}

func TestConvertWithOpts(t *testing.T) {
	t.Run("Pts", func(t *testing.T) {
		ascii, err := ConvertWithOpts(testImg, FontPts(40))
		require.Nil(t, err)
		assert.Nil(t, saveImg(ascii, "convert-pts.jpg"))
	})

	t.Run("Charset", func(t *testing.T) {
		ascii, err := ConvertWithOpts(testImg, CSet(CharsetExtended))
		require.Nil(t, err)
		assert.Nil(t, saveImg(ascii, "convert-cset_extended.jpg"))

		ascii, err = ConvertWithOpts(testImg, CSet(CharsetLimited))
		require.Nil(t, err)
		assert.Nil(t, saveImg(ascii, "convert-cset_limited.jpg"))

		ascii, err = ConvertWithOpts(testImg, CSet(CharsetBlock))
		require.Nil(t, err)
		assert.Nil(t, saveImg(ascii, "convert-cset_block.jpg"))
	})
}

func TestValidOptions(t *testing.T) {
	t.Run("Pts", func(t *testing.T) {
		_, err := ConvertWithOpts(testImg, FontPts(1))
		assert.Nil(t, err)
	})

	t.Run("Font", func(t *testing.T) {
		_, err := ConvertWithOpts(testImg, Font(defaultFont))
		assert.Nil(t, err)
	})

	t.Run("Memory", func(t *testing.T) {
		mem := &Memory{}
		_, err := ConvertWithOpts(testImg, Interpolate(mem))
		assert.Nil(t, err)
	})
}

func TestInvalidOptions(t *testing.T) {
	t.Run("Pts", func(t *testing.T) {
		_, err := ConvertWithOpts(testImg, FontPts(0))
		assert.NotNil(t, err)

		_, err = ConvertWithOpts(testImg, FontPts(-1))
		assert.NotNil(t, err)
	})

	t.Run("Font", func(t *testing.T) {
		_, err := ConvertWithOpts(testImg, Font(nil))
		assert.NotNil(t, err)
	})

	t.Run("Memory", func(t *testing.T) {
		_, err := ConvertWithOpts(testImg, Interpolate(nil))
		assert.NotNil(t, err)
	})
}
