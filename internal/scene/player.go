// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"

	"github.com/3ter/iMagine/internal/fileio"
	"github.com/faiface/pixel/text"
)

// Player is defined by its text and contains the game progression in
// the form of its inventory.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Player struct {
	wordInventory []string

	Texter
}

// SetDefaultAttributes initializes the Player struct
func (p *Player) setDefaultAttributes() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	p.fontFace = face
	// pixel.ZV is the zero vector representing the orig(in) (i.e. beginning of the line)
	p.currentTextObject = text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	p.setTextColor(colornames.Blueviolet)

	p.textBox = new(TextBox)
	// TODO: Find a good way to know the window dimensions here...
	// I've used a potential hack, now we only have to use 'window' here to get relative bounds
	p.textBox.dimensions = pixel.V(900, 230)
	p.textBox.topLeftCorner = pixel.V(1024/2-p.textBox.dimensions.X/2, 768-500)
	p.textBox.thickness = 5
}
