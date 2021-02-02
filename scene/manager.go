// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
//
// This file should provide the orchestration of the game's scenes
package scene

import (
	"io/ioutil"
	"regexp"
)

// ScenesDir publishes the directory where its files are stored
const ScenesDir = `../scene/`

// scenesMap maps scene identifiers (e.g. 'Beach') to their respective scene object
var scenesMap map[string]*Scene

func loadFilesToSceneMap() {
	scenesMap = make(map[string]*Scene)

	sceneFileSlice, err := ioutil.ReadDir(ScenesDir)
	if err != nil {
		panic("Scenes directory '" + ScenesDir + "' couldn't be read!")
	}
	for _, sceneFile := range sceneFileSlice {
		sceneFileFilter := regexp.MustCompile(`^scene(\w+)\.(md|json)$`)

		FileMatchSlice := sceneFileFilter.FindStringSubmatch(sceneFile.Name())
		if len(FileMatchSlice) == 3 {
			FileScene := FileMatchSlice[1]
			FileExtension := FileMatchSlice[2]

			if scenesMap[FileScene] == nil {
				scenesMap[FileScene] = getSceneObjectWithDefaults()
			}

			if FileExtension == `md` {
				scenesMap[FileScene].script.filePath = ScenesDir + FileScene + `.md`
			} else if FileExtension == `json` {
				scenesMap[FileScene].mapConfigPath = ScenesDir + FileScene + `.json`
			}
		}
	}

	// TODO: Load contents as well... let's see the performance hit...
}
