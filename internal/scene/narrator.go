// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"

	"github.com/3ter/iMagine/internal/fileio"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

// Narrator is defined by its text.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Narrator struct {
	fontFace          font.Face
	currentTextObject *text.Text
	currentTextString string

	textBox *TextBox
}

// TODO: I've repeated myself: Those functions are mostly copied from 'player.go'

// SetDefaultAttributes initializes the Player struct
func (n *Narrator) setDefaultAttributes() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	n.fontFace = face
	// pixel.ZV is the zero vector representing the orig(in) (i.e. beginning of the line)
	n.currentTextObject = text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	n.setTextColor(colornames.Blueviolet)

	n.textBox = new(TextBox)
	// TODO: Find a good way to know the window dimensions here...
	// I've used a potential hack, now we only have to use 'window' here to get relative bounds
	n.textBox.dimensions = pixel.V(900, 230)
	n.textBox.topLeftCorner = pixel.V(1024/2-n.textBox.dimensions.X/2, 768-100)
	n.textBox.thickness = 5
}

func (n *Narrator) setTextFontFace(face font.Face) {
	n.currentTextObject = text.New(n.currentTextObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	n.currentTextObject.WriteString(n.currentTextString)
}

func (n *Narrator) setTextColor(col color.RGBA) {
	n.currentTextObject.Color = col
}

func (n *Narrator) setText(str string) {
	n.currentTextObject.Clear()
	n.currentTextObject.WriteString(str)
}

func (n *Narrator) addText(str string) {
	n.currentTextObject.WriteString(str)
}

func (n *Narrator) drawTextInBox(win *pixelgl.Window) {
	n.textBox.drawTextBox(win)

	// margin to the text box in pixels
	margin := 20.0
	// TODO: The y coordinate is guesswork and dependend on the font face used!
	marginVec := pixel.V(margin, margin-55)
	n.currentTextObject.Draw(win, pixel.IM.Moved(n.textBox.topLeftCorner.Add(marginVec)))
}
