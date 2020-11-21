// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// Texter is defined by its text
type Texter struct {
	fontFace font.Face
	// The color in this object is settable, for changing the font a new object
	// is needed.
	currentTextObject *text.Text
	currentTextString string

	textBox *TextBox
}

func (t *Texter) setTextFontFace(face font.Face) {
	t.currentTextObject = text.New(t.currentTextObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	t.currentTextObject.WriteString(t.currentTextString)
}

func (t *Texter) setTextColor(col color.Color) {
	t.currentTextObject.Color = col
}

func (t *Texter) setText(str string) {
	t.currentTextObject.Clear()
	t.currentTextObject.WriteString(str)
}

func (t *Texter) addText(str string) {
	t.currentTextObject.WriteString(str)
}

func (t *Texter) drawTextInBox(win *pixelgl.Window) {
	t.textBox.drawTextBox(win)

	// margin to the text box in pixels
	margin := 20.0
	// TODO: The y coordinate is guesswork and dependend on the font face used!
	marginVec := pixel.V(margin, margin-55)
	t.currentTextObject.Draw(win, pixel.IM.Moved(t.textBox.topLeftCorner.Add(marginVec)))
}
