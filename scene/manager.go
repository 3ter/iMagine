// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
//
// This file should provide the orchestration of the game's scenes
package scene

import (
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/3ter/iMagine/fileio"
)

// ScenesDir publishes the directory where its files are stored
const ScenesDir = `../scene/`

// GlobalScenes maps scene identifiers (e.g. 'Beach') to their respective scene object
var GlobalScenes map[string]*Scene

// GlobalCurrentScene holds the game's state, which can be a scene name like 'Beach' or a state like 'Quit' or 'Pause'.
var GlobalCurrentScene string

// globalPreviousScene is used to determine a scene switch for calculating the number of times a scene has been visited
var globalPreviousScene string

// MapConfig contains key/value-pairs for a scene that are intended to save
// * which scenes are adjacent to the current one
// * the state of the scene
type MapConfig struct {
	North string
	East  string
	South string
	West  string
	Look  string
	// Number of times this scene has been entered
	Visited int
}

func (s *Scene) loadMapConfig(filename string) {
	jsonBytes := fileio.LoadFileToBytes(filename)

	json.Unmarshal(jsonBytes, &s.mapConfig)
}

// LoadFilesToSceneMap fills the global variable 'GlobalScenes' with filepaths and contents.
//
// The file format is 'scene<sceneName>.<fileExtension>':
// - JSON files contain the map config
// - MD files contain the scene's script
// - GO files contain special functions which don't fit in the generic 'OnUpdate' handling
//		- For '.go' files there will be an entry in the 'SceneMap' with default values
//
// For some scenes special init functions are called (e.g. for the 'Demo' scene).
func LoadFilesToSceneMap() {
	GlobalScenes = make(map[string]*Scene)

	sceneFileSlice, err := ioutil.ReadDir(ScenesDir)
	if err != nil {
		panic("Scenes directory '" + ScenesDir + "' couldn't be read!")
	}
	for _, sceneFile := range sceneFileSlice {
		sceneFileFilter := regexp.MustCompile(`^scene(\w+)\.(md|json|go)$`)

		fileMatchSlice := sceneFileFilter.FindStringSubmatch(sceneFile.Name())
		if len(fileMatchSlice) == 3 {
			filePath := ScenesDir + fileMatchSlice[0]
			fileScene := fileMatchSlice[1]
			fileExtension := fileMatchSlice[2]

			if GlobalScenes[fileScene] == nil {
				GlobalScenes[fileScene] = getSceneObjectWithDefaults()
				GlobalScenes[fileScene].Name = fileScene
				if fileScene == `Demo` {
					GlobalScenes[`Demo`].initDemo()
				} else if fileScene == `MainMenu` {
					GlobalScenes[`MainMenu`].initMainMenu()
				}
			}

			if fileExtension == `md` {
				GlobalScenes[fileScene].script.filePath = filePath
				GlobalScenes[fileScene].script.fileContent = fileio.LoadFileToString(filePath)
			} else if fileExtension == `json` {
				GlobalScenes[fileScene].mapConfigPath = filePath
				GlobalScenes[fileScene].loadMapConfig(filePath)
			}
		}
	}
}
