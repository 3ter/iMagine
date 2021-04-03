//The audio player should be able to open an audio file via file name and play multiple tracks/stems in sync
package controlaudio

import (
	"path/filepath"

	"github.com/3ter/iMagine/fileio"
	"github.com/faiface/beep/effects"
)

//AudioPlayer handles all music played throughout the game
type AudioPlayer struct {
	trackDir string
	trackMap map[string]*effects.Volume
}

func (a *AudioPlayer) loadTrack(trackName string) {
	var trackPath = filepath.Join(a.trackDir, trackName)
	a.trackMap[trackName] = fileio.GetStreamer(trackPath)
}

func (a *AudioPlayer) playAllTracks() {

}

func (a *AudioPlayer) trackVolUp(trackName string) {

}

func (a *AudioPlayer) trackVolDown(trackName string) {

}

func (a *AudioPlayer) initAudioPlayer(path string) {
	a.trackDir = path
	a.trackMap = make(map[string]*effects.Volume)

}
