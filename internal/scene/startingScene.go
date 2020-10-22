// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"

	"golang.org/x/image/colornames"
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func GetStartBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}
