// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"github.com/3ter/iMagine/fileio"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

type mainMenuItem struct {
	Text      string
	sceneName string
	State     string
}

var globalMenuItems = []*mainMenuItem{
	{"Demo", "Demo", "selected"},
	{"Start", "Beach", "unselected"},
	{"Quit", "Quit", "unselected"},
}

func (s *Scene) initMainMenu() {
	s.bgColor = colornames.Black
	s.textColor = colornames.White
}

func returnMenuTexts(atlasRegular, atlasBold *text.Atlas) []*text.Text {
	menuTexts := make([]*text.Text, len(globalMenuItems))
	for i, menuItem := range globalMenuItems {
		txt := text.New(pixel.ZV, atlasRegular)
		if menuItem.State == "selected" {
			txt = text.New(pixel.ZV, atlasBold)
		}
		txt.WriteString(menuItem.Text)
		menuTexts[i] = txt
	}

	return menuTexts
}

func (s *Scene) drawMainMenu(win *pixelgl.Window) {
	win.Clear(s.bgColor)

	menuTextVerticalOffset := 50 // pixels

	regularFace := fileio.TtfFromBytesMust(goregular.TTF, 20)
	boldFace := fileio.TtfFromBytesMust(gobold.TTF, 20)
	atlasRegular := text.NewAtlas(regularFace, text.ASCII)
	atlasBold := text.NewAtlas(boldFace, text.ASCII)

	menuTexts := returnMenuTexts(atlasRegular, atlasBold)
	for i, menuText := range menuTexts {
		centerTextMatrix := pixel.IM.Moved(win.Bounds().Center().Sub(menuText.Bounds().Center()))
		verticalAdjustVector := pixel.V(0, float64(-menuTextVerticalOffset*i))
		menuText.Draw(win, centerTextMatrix.Moved(verticalAdjustVector))
	}
}

//HandleMainMenuAndReturnState handles menu items and the associated states
func (s *Scene) onUpdateMainMenu(win *pixelgl.Window) {

	if win.JustPressed(pixelgl.KeyDown) {
		for i, menuItem := range globalMenuItems {
			if menuItem.State == "selected" && i < len(globalMenuItems)-1 {
				menuItem.State = "unselected"
				globalMenuItems[i+1].State = "selected"
				break
			}
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		for i, menuItem := range globalMenuItems {
			if menuItem.State == "selected" && i > 0 {
				menuItem.State = "unselected"
				globalMenuItems[i-1].State = "selected"
				break
			}
		}
	}

	if win.JustPressed(pixelgl.KeyEnter) {
		for _, menuItem := range globalMenuItems {
			if menuItem.State == "selected" {
				GlobalCurrentScene = menuItem.sceneName
			}
		}
	}
}
