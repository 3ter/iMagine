// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"github.com/3ter/iMagine/internal/fileio"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// TextBox has dimensions and coordinates where it is drawn
type TextBox struct {
	// (x, y) providing (width, height) of the box
	dimensions    pixel.Vec
	topLeftCorner pixel.Vec
	thickness     float64
}

func (box *TextBox) drawTextBox(win *pixelgl.Window) {
	texture, err := fileio.LoadPicture("../assets/lavaTexture.jpg")
	if err != nil {
		panic(err)
	}
	imd := imdraw.New(texture)
	imd.Intensity = 1.0
	imd.Picture = box.topLeftCorner
	imd.Push(box.topLeftCorner)
	imd.Picture = box.topLeftCorner.Add(box.dimensions)
	imd.Push(box.topLeftCorner.Add(pixel.V(box.dimensions.X, box.dimensions.Y)))
	imd.Rectangle(box.thickness)
	imd.Draw(win)
}
