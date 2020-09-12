// taken from the tutorial here: https://github.com/faiface/pixel/wiki/Typing-text-on-the-screen

package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

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

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	for !win.Closed() {
		win.Clear(colornames.Black)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)

		for i := 0; i < 6; i++ {
			basicTxt.Orig.X += 10
			basicTxt.Orig.Y += 9000
			fmt.Fprintln(basicTxt, "Orig:", basicTxt.Orig, "Dot:", basicTxt.Dot)
			// After (!) each line written text.Dot takes the X coordinate of text.Orig
		}
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 1))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
