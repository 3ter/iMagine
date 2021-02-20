package scene

import (
	"testing"
)

func TestLoadFilesToSceneMap(t *testing.T) {
	LoadFilesToSceneMap()
	if len(GlobalScenes) == 0 {
		t.Fatalf("GlobalScenes is empty")
	}
	for sceneName, sceneObj := range GlobalScenes {
		if sceneName == `Void` || sceneName == `Demo` || sceneName == `MainMenu` {
			continue
		}
		if len(sceneObj.mapConfigPath) == 0 {
			t.Fatalf("The scene '" + sceneName + "' has no mapConfigPath")
		}
		if len(sceneObj.script.filePath) == 0 {
			t.Fatalf("The scene '" + sceneName + "' has no filePath")
		}

		directions := sceneObj.mapConfig.Directions
		_, directionsHasNorth := directions[`north`]
		if sceneObj.mapConfig != nil && !directionsHasNorth {
			t.Fatalf("The scene '" + sceneName + "' has no valid direction for 'north'")
		}

		if sceneName == `Beach` {
			_, hasKey := sceneObj.objects[`cupboardKey`]
			if !hasKey {
				t.Fatalf("The Beach scene doesn't have an object named 'cupboardKey'")
			}
		}
	}
}
