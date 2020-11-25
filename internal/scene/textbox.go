// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// TextBox has dimensions and coordinates where it is drawn
type TextBox struct {
	// (x, y) providing (width, height) of the box
	dimensions    pixel.Vec
	topLeftCorner pixel.Vec
	thickness     float64
	// margin of the text to the edges of the text box in pixels
	margin float64
}

func (box *TextBox) drawTextBox(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	imd.Push(box.topLeftCorner.Add(pixel.V(0, -box.dimensions.Y)))
	imd.Push(box.topLeftCorner.Add(pixel.V(box.dimensions.X, 0)))
	imd.Rectangle(box.thickness)
	imd.Draw(win)
}
