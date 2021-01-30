package scene

import (
	"testing"
)

func TestLoadFilesToSceneMap(t *testing.T) {
	loadFilesToSceneMap()
	if len(scenesMap) == 0 {
		t.Errorf("scenesMap is empty")
	}
}
