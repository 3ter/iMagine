// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"

	"github.com/3ter/iMagine/controltext"
	"golang.org/x/image/colornames"
)

var (
	// BeachScene holds the texts and music in its Scene struct.
	BeachScene Scene
	// To make a continuous press of backspace possible.
	backspaceCounter int
)

// TODO: The following two functions are not used anymore. Remove them and find out what role these scene
// files should play.

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func getBeachBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}

// TypeBeachTitle prints the text to the respective text elements.
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
