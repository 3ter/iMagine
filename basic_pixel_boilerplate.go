// taken from the tutorial here: https://github.com/faiface/pixel/wiki/Typing-text-on-the-screen

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

func loadTTF(path string, size float64) (font.Face, error) {
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

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true) // remove potential artifacts

	face, err := loadTTF("intuitive.ttf", 52)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(100, 500), atlas)

	for i := 0; i < 6; i++ {
		txt.Orig.X += 10
		txt.Orig.Y += 9000
		fmt.Fprintln(txt, "Orig:", txt.Orig, "Dot:", txt.Dot)
		// After (!) each line written text.Dot takes the X coordinate of text.Orig
	}

	for !win.Closed() {
		win.Clear(colornames.Black)
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 1))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
