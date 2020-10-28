// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music
package scene

import (
	"image/color"
	
	"github.com/3ter/iMagine/internal/scene"
	"golang.org/x/image/colornames"
)

// GetStartBackgroundColor is a placeholder... this probably should go somewhere else
func GetStartBackgroundColor() color.RGBA {
	return colornames.Aliceblue
}


func typeStartTitle() {
	if title.Dot != title.Orig {
		title.Clear()
		title.Color = colornames.Darkgoldenrod
	}
	titleString := "Welcome to the START. Here is nothing... (yet)!\n"
	titleString += "Press Ctrl + Q to quit or Escape for main menu."
	title.WriteString(titleString)
}

func drawStartScene(win *pixelgl.Window) {
	bgColor = scene.GetStartBackgroundColor()
	win.Clear(bgColor)
	// TODO: I don't think it's a good idea to reuse the title from the demo... is it?
	title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, 300)))
}

func handleStartSceneInput(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
}
