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

func TestRemoveMarkdownComments(t *testing.T) {
	testString := `You <span style="text-speed:500">open</span> your ` +
		`<span style="color:Red; font-size:16px;">eyes</span>.
		<!-- This is a great comment. -->

		Text foobar embedded hey!<!-- This is a multiline comment

		embedded into text. -->Mark it down baby!!!

	You find yourself at a beach. <span style="text-speed:2000">You hear the waves come and go </span>, the ` +
		`<span style="color:red">red</span> sunset reflects on the <span style="color:blue">water’s</span> surface.`

	replacedStringLength := len(`<!-- This is a great comment. -->` +
		`<!-- This is a multiline comment

		embedded into text. -->`)

	strippedString := stripMarkdownComments(testString)
	if len(testString)-replacedStringLength != len(strippedString) {
		t.Fatal("Markdown comments have not been removed as planned!")
	}
}
