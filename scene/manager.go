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

// ScenesMap maps scene identifiers (e.g. 'Beach') to their respective scene object
var ScenesMap map[string]*Scene

// CurrentScene holds the game's state, which can be a scene name like 'Beach' or a state like 'Quit' or 'Pause'.
var CurrentScene string

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

	json.Unmarshal(jsonBytes, &s.mapConfig)
}

// LoadFilesToSceneMap fills the global variable 'ScenesMap' with filepaths and contents
func LoadFilesToSceneMap() {
	ScenesMap = make(map[string]*Scene)

	sceneFileSlice, err := ioutil.ReadDir(ScenesDir)
	if err != nil {
		panic("Scenes directory '" + ScenesDir + "' couldn't be read!")
	}
	for _, sceneFile := range sceneFileSlice {
		sceneFileFilter := regexp.MustCompile(`^scene(\w+)\.(md|json)$`)

		fileMatchSlice := sceneFileFilter.FindStringSubmatch(sceneFile.Name())
		if len(fileMatchSlice) == 3 {
			filePath := ScenesDir + fileMatchSlice[0]
			fileScene := fileMatchSlice[1]
			fileExtension := fileMatchSlice[2]

			if ScenesMap[fileScene] == nil {
				ScenesMap[fileScene] = getSceneObjectWithDefaults()
				ScenesMap[fileScene].Name = fileScene
			}

			if fileExtension == `md` {
				ScenesMap[fileScene].script.filePath = filePath
				ScenesMap[fileScene].script.fileContent = fileio.LoadFileToString(filePath)
			} else if fileExtension == `json` {
				ScenesMap[fileScene].mapConfigPath = filePath
				ScenesMap[fileScene].loadMapConfig(filePath)
			}
		}
	}
}
