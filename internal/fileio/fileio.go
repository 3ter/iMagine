// Package fileio implements additional functions to load game specific files
// like fonts or music
package fileio

import (
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// LoadFileToString loads the contents of a file into a string or dies
func LoadFileToString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// LoadTTF has been taken from the pixel Wiki
func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

// TtfFromBytesMust has been taken from pixel-examples "Typewriter"
// https://github.com/faiface/pixel-examples/tree/master/typewriter
func TtfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

// GetStreamer had initially been taken from the pixel Wiki
func GetStreamer(filePath string) *effects.Volume {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	return volume
}

// LoadPicture has been copied from
// https://github.com/faiface/pixel/wiki/Drawing-a-Sprite
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
