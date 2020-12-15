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
	// To make a continuous press of backspace possible.
	backspaceCounter int
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func getBeachBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}

// TypeBeachTitle prints the text to the respective text elements.
// TODO: This should use pixelgl's typed like in the demo and do something on
// an answer.
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
	s.hint.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.hint.Bounds().Center())).Moved(
		pixel.V(0, 2*s.hint.Bounds().H())))

	player.drawTextInBox(win)
	narrator.drawTextInBox(win)
}

func handleBackspace(win *pixelgl.Window, player *Player) {
	if win.JustPressed(pixelgl.KeyBackspace) && len(player.currentTextString) > 0 {
		player.setText(player.currentTextString[:len(player.currentTextString)-1])
		backspaceCounter = int(-120 * 0.5) // Framerate times seconds to wait until continuous backspace kicks in.
	} else if win.Pressed(pixelgl.KeyBackspace) && len(player.currentTextString) > 0 {
		backspaceCounter++
		backspaceDeletionSpeed := int(120 / 40) // Framerate divided by deletions per second.
		if backspaceCounter > 0 && backspaceCounter%backspaceDeletionSpeed == 0 {
			player.setText(player.currentTextString[:len(player.currentTextString)-1])
			backspaceCounter = 0
		}
	}
}

// HandleBeachSceneInput listens and processes player input.
func (s *Scene) HandleBeachSceneInput(win *pixelgl.Window, gameState string) string {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
	handleBackspace(win, &player)
	if win.JustPressed(pixelgl.KeyEnter) {
		if len(s.script.responseQueue) == 0 && len(s.script.keywordResponseMap) == 0 {
			s.parseScriptFile()
		}
		s.executeScriptFromQueue()
	}

	if len(win.Typed()) > 0 {
		player.addText(win.Typed())
	}

	return gameState
}
