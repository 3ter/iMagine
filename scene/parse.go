// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
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

// translateDirectionToSceneName uses the Directions map from mapConfig to return the sceneName
//
// It defaults back to returning the original string in case there was no direction matching the string
// e.g. Directions[`foobar`] -> no direction with foobar so print out foobar's  not a valid direction
func translateDirectionToSceneName(direction string) string {

	sceneName := GlobalScenes[GlobalCurrentScene].mapConfig.Directions[direction]
	if sceneName == `` {
		return direction
	}
	return sceneName
}

func (s *Scene) handleSpecialPlayerCommands(playerWords []string) {

	if len(playerWords) < 2 {
		globalNarrator.setTextLetterByLetter("Specify your command in the format: '[verb] [object]'", s)
		return
	}
	verb := strings.ToLower(playerWords[0])
	object := strings.ToLower(playerWords[1])
	switch verb {
	case `go`:
		sceneName := translateDirectionToSceneName(object)
		if GlobalScenes[sceneName] == nil || sceneName == `Void` {
			globalNarrator.setTextLetterByLetter("You can't go to '"+sceneName+"'! (Enter a direction: e.g. North)", s)
			return
		}
		// To allow parsing of the newly selected current script file (see 'scene.OnUpdate')
		s.script.keywordResponseMap = nil
		GlobalCurrentScene = sceneName
	case `look`:
		sceneName := translateDirectionToSceneName(object)
		globalNarrator.setTextLetterByLetter(GlobalScenes[sceneName].mapConfig.Look, s)
	}
}

func (s *Scene) handlePlayerCommand(playerInput string) {

	playerWords := strings.Split(playerInput, ` `)

	if len(playerWords[0]) > 0 && (playerWords[0] == `go` || playerWords[0] == `look`) {
		s.handleSpecialPlayerCommands(playerWords)
		return
	}

	// Check for progress change
	for keyword, responseSlice := range s.script.keywordResponseMap {
		if strings.ToLower(playerInput) == strings.ToLower(keyword) {
			// If there's a progressUpdate then there's only one response in the slice
			if responseSlice[0].progressUpdate != "" {
				s.progress = responseSlice[0].progressUpdate
				// Empty keywordResponseMap to prepare for jump to new script section.
				s.script.keywordResponseMap = map[string][]narratorResponse{}
				s.parseScriptFile()
				s.executeScriptFromQueue()
			} else {
				executeAmbienceCommands(s.script.keywordResponseMap[keyword][0].ambienceCmdSlice)
				globalNarrator.setTextLetterByLetter(s.script.keywordResponseMap[keyword][0].narratorTextLine, s)
			}
		}
	}
}

// executeScriptFromQueue modfies the scene according to scene script and player input.
//
// If the scene modifications are still to be fed from the 'responseQueue' the function returns without checking player
// input.
// The check for progress change provides a way to jump from section within a scene script.
func (s *Scene) executeScriptFromQueue() {

	if len(s.script.responseQueue) > 0 {
		executeAmbienceCommands(s.script.responseQueue[0].ambienceCmdSlice)
	}

	// Set narrator text
	if len(s.script.responseQueue) > 0 {
		globalNarrator.setTextLetterByLetter(s.script.responseQueue[0].narratorTextLine, s)
		s.script.responseQueue = s.script.responseQueue[1:]
		return
	}

	playerInput := globalPlayer.currentTextString
	globalPlayer.setText("")

	s.handlePlayerCommand(playerInput)
}

func (s *Scene) parseScriptFile() {

	activeScriptSlice := s.getActiveScriptSlice()

	var ambienceCmdSlice []string
	for lineNumber, scriptLine := range activeScriptSlice {

		response := getCombinedAmbienceTextResponse(scriptLine, &ambienceCmdSlice)
		if response.narratorTextLine != "" {
			s.script.responseQueue = append(s.script.responseQueue, response)
		}

		// break loop if the first player keyword is found and gobble up the rest of the lines
		keywordResponseMap := getKeywordResponseMap(scriptLine, lineNumber, activeScriptSlice)
		if len(keywordResponseMap) > 0 {
			s.script.keywordResponseMap = keywordResponseMap
			break
		}
	}
}
