// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"

	"github.com/faiface/pixel"

	"github.com/3ter/iMagine/internal/fileio"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// Player is defined by its text and contains the game progression in
// the form of its inventory.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Player struct {
	fontFace font.Face
	// The color in this object is settable, for changing the font a new object
	// is needed.
	currentTextObject *text.Text
	currentTextString string

	textBox TextBox

	wordInventory []string
}

func (p *Player) setDefaultAttributes() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	p.fontFace = face
	// pixel.ZV is the zero vector representing the orig(in) (i.e. beginning of the line)
	p.currentTextObject = text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	p.setTextColor(color.White)
}

func (p *Player) setTextFontFace(face font.Face) {
	p.currentTextObject = text.New(p.currentTextObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	p.currentTextObject.WriteString(p.currentTextString)
}

func (p *Player) setTextColor(col color.Color) {
	p.currentTextObject.Color = col
}

func (p *Player) setText(str string) {
	p.currentTextObject.Clear()
	p.currentTextObject.WriteString(str)
}

func (p *Player) addText(str string) {
	p.currentTextObject.WriteString(str)
}

func (p *Player) drawTextInBox() {
	p.textBox.drawTextBox()

}
