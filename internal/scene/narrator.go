// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// Narrator is defined by its text.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Narrator struct {
	face              font.Face
	currentTextObject *text.Text
	currentTextString string
}

func (p *Narrator) setTextFontFace(face font.Face) {
	p.currentTextObject = text.New(p.currentTextObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	p.currentTextObject.WriteString(p.currentTextString)
}

func (p *Narrator) setTextColor(col color.RGBA) {
	p.currentTextObject.Color = col
}

func (p *Narrator) setText(str string) {
	p.currentTextObject.Clear()
	p.currentTextObject.WriteString(str)
}

func (p *Narrator) addText(str string) {
	p.currentTextObject.WriteString(str)
}
