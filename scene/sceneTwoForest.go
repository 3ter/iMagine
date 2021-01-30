// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	// ForestScene ...
	ForestScene Scene
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func getForestBackgroundColor() color.RGBA {
	return colornames.Forestgreen
}

// TypeForestTitle ...
func (s *Scene) TypeForestTitle() {
	if s.title.Dot != s.title.Orig {
		s.title.Clear()
		s.title.Color = colornames.Darkgoldenrod
	}
	titleString := "Welcome to the FOREST. There are some trees here!\n"
	titleString += "Press Ctrl + Q to quit or Escape for main menu."
	s.title.WriteString(titleString)
}

// DrawForestScene ...
func (s *Scene) DrawForestScene(win *pixelgl.Window) {
	s.bgColor = getForestBackgroundColor()
	win.Clear(s.bgColor)
	s.title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.title.Bounds().Center())).Moved(pixel.V(0, 300)))
}

// HandleForestSceneInput ...
func (s *Scene) HandleForestSceneInput(win *pixelgl.Window, gameState string) string {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
	return gameState
}
