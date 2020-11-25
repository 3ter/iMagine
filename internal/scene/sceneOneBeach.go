// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"

	"github.com/3ter/iMagine/internal/controltext"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	// BeachScene holds the texts and music in its Scene struct.
	BeachScene Scene
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func getBeachBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}

// TypeBeachTitle prints the text to the respective text elements.
// TODO: This should use pixelgl's typed like in the demo and do something on
// an answer.
// TODO: may want to add a writetotextletterbyletter in the demo scene so all features
// are in one place.
func (s *Scene) TypeBeachTitle() {

	s.title.Clear()
	s.title.Color = colornames.Darkgoldenrod
	titleString := "Welcome to the START. Here is nothing... (yet)!\n"
	writingDoneChannel := make(chan int)
	go controltext.WriteToTextLetterByLetter(s.title, titleString, 60, writingDoneChannel)
	writingDoneChannel <- 1 // init writing the first line
	titleString = "Press Ctrl + Q to quit or Escape for main menu.\n"
	go controltext.WriteToTextLetterByLetter(s.title, titleString, 10, writingDoneChannel)
	titleString = "Press Enter to go to the next area!"
	go controltext.WriteToTextLetterByLetter(s.title, titleString, 10, writingDoneChannel)
}

// DrawBeachScene draws background and text to the window.
func (s *Scene) DrawBeachScene(win *pixelgl.Window) {
	s.bgColor = getBeachBackgroundColor()
	win.Clear(s.bgColor)
	s.title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.title.Bounds().Center())).Moved(pixel.V(0, 300)))

	player.drawTextInBox(win)
	narrator.drawTextInBox(win)
}

// HandleBeachSceneInput listens and processes player input.
func (s *Scene) HandleBeachSceneInput(win *pixelgl.Window, gameState string) string {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		s.handlePlayerCommand()
	}

	if len(win.Typed()) > 0 {
		player.addText(win.Typed())
	}

	return gameState
}

// handlePlayerCommand sets both player and narrator text to be drawn afterwards.
func (s *Scene) handlePlayerCommand() {

	var playerText string
	var narratorText string

	switch s.sceneProgress {
	case "beginning":
		narratorText = `You open your eyes.
You find yourself at a beach. You hear the waves come and go, the red sunset reflects on the water's surface.
As the sunlight falls, a shiny reflection catches your eye.`
		if player.currentTextString == `inspect reflection` {
			s.sceneProgress = `compass 1`
			s.handlePlayerCommand()
			return
		}
	case `compass 1`:
		narratorText = `You walk closer to whatever it is that caught your eye. It was glass that reflected sunlight into your eyes. Glass that belonged to a little device. A compass.`
	}

	narrator.setText(narratorText)
	player.setText(playerText)
}
