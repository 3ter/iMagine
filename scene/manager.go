// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
//
// This file should provide the orchestration of the game's scenes
package scene

import "io/ioutil"

// ScenesDir publishes the directory where its files are stored
const ScenesDir = `../scene/`

// scenesMap maps scene identifiers (e.g. 'Beach') to their respective scene object
var scenesMap map[string]Scene

func loadFilesToSceneMap() {
	sceneFileSlice, err := ioutil.ReadDir(ScenesDir)
	if err != nil {
		panic("Scenes directory '" + ScenesDir + "' couldn't be read!")
	}
	for _, sceneFile := range sceneFileSlice {
		_ = sceneFile
	}
}
