// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import(

	"image/color"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var(
	BeachScene Scene
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func getBeachBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}


func (s *Scene) TypeBeachTitle() {
	if s.title.Dot != s.title.Orig {
		s.title.Clear()
		s.title.Color = colornames.Darkgoldenrod
	}
	titleString := "Welcome to the START. Here is nothing... (yet)!\n"
	titleString += "Press Ctrl + Q to quit or Escape for main menu."
	titleString += "Press Enter to go to the next area!"
	s.title.WriteString(titleString)
}

func (s *Scene) DrawBeachScene(win *pixelgl.Window) {
	s.bgColor = getBeachBackgroundColor()
	win.Clear(s.bgColor)
	s.title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.title.Bounds().Center())).Moved(pixel.V(0, 300)))
}

func (s *Scene) HandleBeachSceneInput(win *pixelgl.Window, gameState string) string{
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		gameState = "Forest"
	}

	return gameState
}
