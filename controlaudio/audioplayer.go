//The audio player should be able to open an audio file via file name and play multiple tracks/stems in sync
package controlaudio

type AudioPlayer struct {
	trackPath string
}

func (a *AudioPlayer) loadTrack(path string) {
	a.trackPath = path
}
