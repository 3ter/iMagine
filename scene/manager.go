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
	scenesMap := make(map[string]*Scene)

	sceneFileSlice, err := ioutil.ReadDir(ScenesDir)
	if err != nil {
		panic("Scenes directory '" + ScenesDir + "' couldn't be read!")
	}
	for _, sceneFile := range sceneFileSlice {
		sceneMarkdownFileFilter := regexp.MustCompile(`^scene(\w+)\.md$`)
		sceneJSONFileFilter := regexp.MustCompile(`^scene(\w+)\.json$`)

		markdownFileMatchSlice := sceneMarkdownFileFilter.FindStringSubmatch(sceneFile.Name())
		if markdownFileMatchSlice != nil {
			markdownFileScene := markdownFileMatchSlice[1]
			scenesMap[markdownFileScene].script.filePath = ScenesDir + `/` + markdownFileScene + `.md`
			continue
		}
		jsonFileMatchSlice := sceneJSONFileFilter.FindStringSubmatch(sceneFile.Name())
		if jsonFileMatchSlice != nil {
			jsonFileScene := jsonFileMatchSlice[1]
			scenesMap[jsonFileScene].mapConfigPath = ScenesDir + `/` + jsonFileScene + `.md`
			continue
		}
	}

	// TODO: Init the Scene (key) in the map at the right time (and only once)
	// TODO: Load contents as well... let's see the performance hit...
	_ = sceneFileSlice
}
