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

// scenesMap maps scene identifiers (e.g. 'Beach') to their respective scene object
var scenesMap map[string]*Scene

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

func loadFilesToSceneMap() {
	scenesMap = make(map[string]*Scene)

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

			if scenesMap[fileScene] == nil {
				scenesMap[fileScene] = getSceneObjectWithDefaults()
			}

			if fileExtension == `md` {
				scenesMap[fileScene].script.filePath = filePath
				scenesMap[fileScene].script.fileContent = fileio.LoadFileToString(filePath)
			} else if fileExtension == `json` {
				scenesMap[fileScene].mapConfigPath = filePath
				scenesMap[fileScene].loadMapConfig(filePath)
			}
		}
	}

	// TODO: Load contents as well... let's see the performance hit...
}
