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

// ContentDir publishes the directory where its files are stored
const ContentDir = `../scene/content/`

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
	// Directions maps north, east, south and west to their respective scene names
	Directions map[string]string
	Look       string
	// Number of times this scene has been entered
	Visited int
}

func (s *Scene) loadMapConfig(filename string) {
	jsonBytes := fileio.LoadFileToBytes(filename)

	json.Unmarshal(jsonBytes, &s.mapConfig)
}

func (s *Scene) loadObject(filename string, objectName string) {
	jsonBytes := fileio.LoadFileToBytes(filename)

	var objectData map[string]interface{}
	json.Unmarshal(jsonBytes, &objectData)

	s.objects[objectName] = objectData

	// TODO: I have no complete concept of how to deal with the unstructured data I've got here now.
	// Data needs to be cast: s.objects[objectName][`id`].(string)
	// I suppose we need to try and access certain keys in the map and determine the actions from there.
}

// isTestFile is a helper to skip go test files when looking for scene files
func isTestFile(filename string) bool {
	matchTestFile := regexp.MustCompile(`_test.go$`)
	return matchTestFile.MatchString(filename)
}

func buildSceneFromFolder(foldername string) {
	sceneName := foldername

	if GlobalScenes[sceneName] == nil {
		GlobalScenes[sceneName] = getSceneObjectWithDefaults()
		GlobalScenes[sceneName].Name = sceneName
		if sceneName == `Demo` {
			GlobalScenes[`Demo`].initDemo()
		} else if sceneName == `MainMenu` {
			GlobalScenes[`MainMenu`].initMainMenu()
		}
	}
}

// LoadFilesToSceneMap fills the global variable 'GlobalScenes' with filepaths and contents.
//
// Every file for a scene has its own directory named with the scene name (its identifier throughout the game).
// The files can be 'script.md', 'mapConfig.json' or '<objectName>.json' (not yet implemented):
// - JSON files contain the map config
// - MD files contain the scene's script
//
// GO files are outside this structure and contain special functions which don't fit in the generic 'OnUpdate' handling.
// For empty folders there will be an entry in the 'SceneMap' with default values.
//
// For some scenes special init functions are called (e.g. for the 'Demo' scene).
func LoadFilesToSceneMap() {
	GlobalScenes = make(map[string]*Scene)

	contentFolders, err := ioutil.ReadDir(ContentDir)
	if err != nil {
		panic("Content directory '" + ContentDir + "' couldn't be read!")
	}
	for _, contentFolder := range contentFolders {

		sceneName := contentFolder.Name()
		buildSceneFromFolder(sceneName)
		GlobalScenes[sceneName].objects = make(map[string]map[string]interface{})

		contentFiles, err := ioutil.ReadDir(ContentDir + contentFolder.Name())
		if err != nil {
			panic("Content directory '" + ContentDir + contentFolder.Name() + "' couldn't be read!")
		}
		for _, contentFile := range contentFiles {
			if isTestFile(contentFile.Name()) {
				continue
			}

			contentFileFilter := regexp.MustCompile(`^(\w+)\.(md|json)$`)

			fileMatchSlice := contentFileFilter.FindStringSubmatch(contentFile.Name())
			if len(fileMatchSlice) == 3 {
				filePath := ContentDir + sceneName + "/" + fileMatchSlice[0]
				fileName := fileMatchSlice[1]
				fileExtension := fileMatchSlice[2]

				if fileName == `script` && fileExtension == `md` {
					GlobalScenes[sceneName].script.filePath = filePath
					GlobalScenes[sceneName].script.fileContent = fileio.LoadFileToString(filePath)
				} else if fileName == `mapConfig` && fileExtension == `json` {
					GlobalScenes[sceneName].mapConfigPath = filePath
					GlobalScenes[sceneName].loadMapConfig(filePath)
				} else {
					GlobalScenes[sceneName].loadObject(filePath, fileName)
				}
			}
		}
	}
}
