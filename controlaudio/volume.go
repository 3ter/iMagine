// Package controlaudio implements functions for further audio control with
// the beep library
package controlaudio

import "github.com/faiface/beep/effects"

// VolumeUp adds linear increments to the float controlling the volume of a track
func VolumeUp(track *effects.Volume) {
	if track.Silent {
		track.Volume = 0.5
		track.Silent = false
	} else {
		track.Volume += 0.5
	}
}

// VolumeDown subtracts linear increments to the float controlling the volume of a track
func VolumeDown(track *effects.Volume) {
	if track.Volume <= 0.5 {
		track.Silent = true
	} else {
		track.Volume -= 0.5
	}
}
