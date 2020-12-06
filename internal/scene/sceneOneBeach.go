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

// getCombinedAmbienceTextResponse collects ambience commands in a slice until the first text line gets
// the accumulated commands added to the narratorResponse struct.
//
// This could also be done with a closure but the additional parameter seemed easier.
func getCombinedAmbienceTextResponse(line string, ambienceCmdSlice []string) narratorResponse {
	ambienceCmd := getMatchedAmbienceCmd(line)
	if len(ambienceCmd) > 0 {
		ambienceCmdSlice = append(ambienceCmdSlice, ambienceCmd)
	}

	// Empty the previous ambience slice to fill the text (non-command) line with it
	cmdRegexp := regexp.MustCompile("^`")
	if !cmdRegexp.MatchString(line) {
		response := narratorResponse{
			narratorTextLine: line,
		}
		response.ambienceCmdSlice = ambienceCmdSlice
		ambienceCmdSlice = nil
		return response
	}
	return narratorResponse{}
}

// getActiveScriptSlice uses the already loaded script data and returns the current script based on s.progress.
func (s *Scene) getActiveScriptSlice() []string {
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

	return activeScriptSlice
}

// getKeywordResponseMap gobbles up the rest of the lines from lineNumber onwards when it encounters the first
// player command marker (parentheses).
func getKeywordResponseMap(line string, lineNumber int, activeScriptSlice []string) map[string][]narratorResponse {

	playerCmdMarkerRegexp := regexp.MustCompile("^`\\((.+)\\)(?: > (.+))?`$")
	submatchSlice := playerCmdMarkerRegexp.FindStringSubmatch(line)
	if len(submatchSlice) == 0 {
		return map[string][]narratorResponse{}
	}

	var keywordResponseMap = make(map[string][]narratorResponse)
	var ambienceCmdSlice []string

	var currentKeyword string
	for _, scriptLine := range activeScriptSlice[lineNumber:] {
		submatchSlice := playerCmdMarkerRegexp.FindStringSubmatch(scriptLine)
		if len(submatchSlice) > 0 {
			currentKeyword = submatchSlice[1]
		}
		if len(submatchSlice) > 2 && submatchSlice[2] != "" {
			// Two (sub)matches mean no more messages to come but a jump to a new progress state
			keywordResponseMap[currentKeyword] =
				append(keywordResponseMap[currentKeyword], narratorResponse{
					progressUpdate: submatchSlice[2],
				})
			continue
		}

		response := getCombinedAmbienceTextResponse(scriptLine, ambienceCmdSlice)
		if response.narratorTextLine != "" {
			keywordResponseMap[currentKeyword] =
				append(keywordResponseMap[currentKeyword], response)
		}
	}

	return keywordResponseMap
}

// executeScriptFromQueue returns true if the queue is empty and all commands have been executed.
func (s *Scene) executeScriptFromQueue() {

	// Execute ambience directives
	// TODO: Implement it

	// Set narrator text
	if len(s.script.responseQueue) > 0 {
		narrator.setText(s.script.responseQueue[0].narratorTextLine)
		s.script.responseQueue = s.script.responseQueue[1:]
		return
	}

	playerProvidedKeyword := player.currentTextString
	player.setText("")

	// Check for progress change
	for keyword, responseSlice := range s.script.keywordResponseMap {
		if strings.ToLower(playerProvidedKeyword) == strings.ToLower(keyword) {
			// If there's a progressUpdate then there's only one response in the slice
			if responseSlice[0].progressUpdate != "" {
				s.progress = responseSlice[0].progressUpdate
				// Empty keywordResponseMap to prepare for jump to new script section.
				s.script.keywordResponseMap = map[string][]narratorResponse{}
				s.parseScriptFile()
				s.executeScriptFromQueue()
			} else {
				narrator.setText(s.script.keywordResponseMap[keyword][0].narratorTextLine)
			}
		}
	}
}

func (s *Scene) parseScriptFile() {

	activeScriptSlice := s.getActiveScriptSlice()

	var ambienceCmdSlice []string
	for lineNumber, scriptLine := range activeScriptSlice {

		response := getCombinedAmbienceTextResponse(scriptLine, ambienceCmdSlice)
		if response.narratorTextLine != "" {
			s.script.responseQueue = append(s.script.responseQueue, response)
		}

		keywordResponseMap := getKeywordResponseMap(scriptLine, lineNumber, activeScriptSlice)
		if len(keywordResponseMap) > 0 {
			s.script.keywordResponseMap = keywordResponseMap
			break
		}
	}
}
