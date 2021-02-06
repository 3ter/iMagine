package scene

import (
	"testing"
)

func TestLoadFilesToSceneMap(t *testing.T) {
	LoadFilesToSceneMap()
	if len(ScenesMap) == 0 {
		t.Errorf("scenesMap is empty")
	}
	for sceneName, sceneObj := range ScenesMap {
		if sceneName == `Void` {
			continue
		}
		if len(sceneObj.mapConfigPath) == 0 {
			t.Errorf("The scene '" + sceneName + "' has no mapConfigPath")
		}
		if len(sceneObj.script.filePath) == 0 {
			t.Errorf("The scene '" + sceneName + "' has no filePath")
		}
	}
}
