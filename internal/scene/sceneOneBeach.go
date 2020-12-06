// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"
	"regexp"
	"strings"

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
		// TODO: Remove old command if fully replaced
		// s.handlePlayerCommand()
		s.enactScriptFile()
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

	switch s.progress {
	case "beginning":
		narratorText = `You open your eyes.
You find yourself at a beach. You hear the waves come and go, the red sunset reflects on the water's surface.
As the sunlight falls, a shiny reflection catches your eye.`
		if player.currentTextString == `inspect reflection` {
			s.progress = `compass 1`
			s.handlePlayerCommand()
			return
		}
	case `compass 1`:
		narratorText = `You walk closer to whatever it is that caught your eye. It was glass that reflected sunlight into your eyes. Glass that belonged to a little device. A compass.`
	}

	narrator.setText(narratorText)
	player.setText(playerText)
}

// getMatchedAmbienceCmd removes the command marker from a string and returns it.
// If no match was found it returns an empty string which has length 0.
func getMatchedAmbienceCmd(line string) string {
	ambienceCmdMarkerRegexp := regexp.MustCompile("^`\\[(.+)\\]`$")
	submatchSlice := ambienceCmdMarkerRegexp.FindStringSubmatch(line)
	if len(submatchSlice) > 0 {
		return submatchSlice[1]
	}
	return ""
}

func (s *Scene) enactScriptFile() {
	var activeScript string

	// Find currently active script part and remove progress line
	hashRegexp := regexp.MustCompile(`(?m:^# )`)
	scriptParts := hashRegexp.Split(s.script.file, -1)
	progressRegexp := regexp.MustCompile(`^` + s.progress)
	for _, scriptPart := range scriptParts {
		if progressRegexp.MatchString(scriptPart) {
			untilFirstLineEndRegexp := regexp.MustCompile(`^\w+\n`)
			activeScript = untilFirstLineEndRegexp.Split(scriptPart, 2)[1]
			break
		}
	}

	// Separate directives (ambience / text / keywords) by a blank line
	blankLineRegexp := regexp.MustCompile(`\n\n`)
	activeScriptSlice := blankLineRegexp.Split(activeScript, -1)
	activeScriptSlice = activeScriptSlice[:len(activeScriptSlice)-1] // to remove last element (empty string)

	var ambienceCmdSlice []string
	var narratorTextSlice []string
	var playerCmdToResponseMap = make(map[string]narratorResponse)

	var narratorResponseSlice []narratorResponse

	var submatchSlice []string
	for lineNumber, scriptLine := range activeScriptSlice {
		ambienceCmd := getMatchedAmbienceCmd(scriptLine)
		if len(ambienceCmd) > 0 {
			ambienceCmdSlice = append(ambienceCmdSlice, ambienceCmd)
		}

		// Empty the previous ambience slice to fill the text (non-command) line with it
		cmdRegexp := regexp.MustCompile("^`")
		if !cmdRegexp.MatchString(scriptLine) {
			response := narratorResponse{
				narratorTextLine: scriptLine,
			}
			response.ambienceCmdSlice = ambienceCmdSlice
			ambienceCmdSlice = nil
			narratorResponseSlice = append(narratorResponseSlice, response)
		}

		// Get response directives
		playerCmdMarkerRegexp := regexp.MustCompile("^`\\((.+)\\)(?: > (.+))?`$")
		submatchSlice = playerCmdMarkerRegexp.FindStringSubmatch(scriptLine)
		var playerCmd string
		if len(submatchSlice) > 0 {
			playerCmd = submatchSlice[1]
		}
		if len(submatchSlice) == 2 {
			// Only one (sub)match means no jump so push all into the queue
			var narratorResponseSlice []string
			for _, narratorTextLine := range activeScriptSlice[lineNumber:] {
				cmdMarkerRegexp := regexp.MustCompile("^`")
				if cmdMarkerRegexp.MatchString(narratorTextLine) {
					break
				} else {
					narratorResponseSlice = append(narratorResponseSlice, narratorTextLine)
				}
			}
			playerCmdToResponseMap[playerCmd] = narratorResponse{narratorTextLine: err}
		} else if len(submatchSlice) == 3 {
			// Two (sub)matches mean no more messages to come but a jump to a new progress state
			playerCmdToResponseMap[playerCmd] = narratorResponse{progressUpdate: submatchSlice[2]}
		}
	}

	s.script.responseQueue = narratorResponseSlice

	// Execute ambience directives
	// TODO: Implement

	// Check if progress change
	for keyword, response := range playerCmdToResponseMap {
		if strings.ToLower(player.currentTextString) == strings.ToLower(keyword) {
			if response.progressUpdate != "" {
				s.progress = response.progressUpdate
			}
		}
	}

	narrator.setText(narratorTextSlice[0])
	player.setText("")
}
