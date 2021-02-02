// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/3ter/iMagine/fileio"
	"github.com/faiface/beep/speaker"
)

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
func getCombinedAmbienceTextResponse(line string, ambienceCmdSlice *[]string) narratorResponse {
	ambienceCmd := getMatchedAmbienceCmd(line)
	if len(ambienceCmd) > 0 {
		*ambienceCmdSlice = append(*ambienceCmdSlice, ambienceCmd)
	}

	// Empty the previous ambience slice to fill the text (non-command) line with it
	cmdRegexp := regexp.MustCompile("^`")
	if !cmdRegexp.MatchString(line) {
		response := narratorResponse{
			narratorTextLine: line,
		}
		response.ambienceCmdSlice = *ambienceCmdSlice
		*ambienceCmdSlice = nil
		return response
	}
	return narratorResponse{}
}

// getActiveScriptSlice uses the already loaded script data and returns the current script based on s.progress.
func (s *Scene) getActiveScriptSlice() []string {

	if len(s.script.fileContent) <= 0 {
		panic("Script file hasn't been loaded into string.")
	}

	var activeScript string

	// Find currently active script part and remove progress line
	hashRegexp := regexp.MustCompile(`(?m:^# )`)
	scriptParts := hashRegexp.Split(s.script.fileContent, -1)
	if len(scriptParts) <= 0 {
		panic("Script doesn't contain at least one part marked by '#'.")
	}

	progressRegexp := regexp.MustCompile(`^` + s.progress + `\r?\n`)
	for _, scriptPart := range scriptParts {
		if progressRegexp.MatchString(scriptPart) {
			activeScript = progressRegexp.ReplaceAllString(scriptPart, ``)
			break
		}
	}
	if len(activeScript) <= 0 {
		panic("Active script empty after removal of progress marker.")
	}

	// Separate directives (ambience / text / keywords) by a blank line
	blankLineRegexp := regexp.MustCompile(`\r?\n\r?\n`)
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

		response := getCombinedAmbienceTextResponse(scriptLine, &ambienceCmdSlice)
		if response.narratorTextLine != "" {
			keywordResponseMap[currentKeyword] =
				append(keywordResponseMap[currentKeyword], response)
		}
	}

	return keywordResponseMap
}

func executeAmbienceCommands(ambienceCmdSlice []string) {
	for _, ambienceCmd := range ambienceCmdSlice {
		ambientTypeRegexp := regexp.MustCompile(`^(\w+):\s?`)
		ambientType := ambientTypeRegexp.FindStringSubmatch(ambienceCmd)[1]

		switch ambientType {
		case `Audio`:
			// Don't allow whitespace chars in filenames
			audioFileRegexp := regexp.MustCompile(`^Audio:\s?(\S+)`)
			audioFilename := audioFileRegexp.FindStringSubmatch(ambienceCmd)[1]
			var streamer = fileio.GetStreamer("../assets/" + audioFilename)
			speaker.Play(streamer)
		}
	}
}

// executeScriptFromQueue returns true if the queue is empty and all commands have been executed.
func (s *Scene) executeScriptFromQueue() {

	if len(s.script.responseQueue) > 0 {
		executeAmbienceCommands(s.script.responseQueue[0].ambienceCmdSlice)
	}

	// Set narrator text
	if len(s.script.responseQueue) > 0 {
		narrator.setTextLetterByLetter(s.script.responseQueue[0].narratorTextLine, s)
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
				executeAmbienceCommands(s.script.keywordResponseMap[keyword][0].ambienceCmdSlice)
				narrator.setTextLetterByLetter(s.script.keywordResponseMap[keyword][0].narratorTextLine, s)
			}
		}
	}
}

func (s *Scene) parseScriptFile() {

	activeScriptSlice := s.getActiveScriptSlice()

	var ambienceCmdSlice []string
	for lineNumber, scriptLine := range activeScriptSlice {

		response := getCombinedAmbienceTextResponse(scriptLine, &ambienceCmdSlice)
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

// MapConfig contains key/value-pairs for a scene that are intended to save
// * which scenes are adjacent to the current one
// * the state of the scene
type MapConfig struct {
	North string
	East  string
	South string
	West  string
	Look  string
}

func (s *Scene) loadMapConfig(filename string) {
	jsonBytes := fileio.LoadFileToBytes(filename)

	var mapConfig MapConfig
	json.Unmarshal(jsonBytes, &mapConfig)

}