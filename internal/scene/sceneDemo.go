//Demo Scene to test basic capabilities we need
package scene

import "github.com/3ter/iMagine/internal/scene/Scene"

var (
	isMusicPlaying = false
	trackArray     = [4]string{"Celesta.ogg", "Choir.ogg", "Harp.ogg", "Strings.ogg"}
	trackPath      = "../assets/"
	trackMap       map[int]*effects.Volume
)




func handleDemoInput(win *pixelgl.Window, start time.Time) {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyU) {
		controlaudio.VolumeUp(trackMap[0])
	}
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyJ) {
		controlaudio.VolumeDown(trackMap[0])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyI) {
		controlaudio.VolumeUp(trackMap[1])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyK) {
		controlaudio.VolumeDown(trackMap[1])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyO) {
		controlaudio.VolumeUp(trackMap[2])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
		controlaudio.VolumeDown(trackMap[2])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyP) {
		controlaudio.VolumeUp(trackMap[3])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeySemicolon) {
		controlaudio.VolumeDown(trackMap[3])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyA) {
		allStreamer := beep.Mix(trackMap[0], trackMap[1], trackMap[2], trackMap[3])
		speaker.Play(allStreamer)
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
		// TODO: Make it a toggle (set a default fragment shader..?)
		applyShader(win, start)
		isShaderApplied = true
	}

	if win.Pressed(pixelgl.KeyDown) {
		bgColor.R--
	}
	if win.Pressed(pixelgl.KeyUp) {
		bgColor.R++
	}

	txt.WriteString(win.Typed())
	typed = typed + win.Typed()

	if win.JustPressed(pixelgl.KeyEnter) {
		// Mind that GLFW doesn't support {Enter} (and {Tab}) (yet)
		// txt.WriteRune('\n')
		title.Clear()
		rgb := convertTextToRGB(typed)
		title.Color = color.RGBA{rgb[0], rgb[1], rgb[2], 0xff}
		title.WriteString("That worked quite well! (You can do that again)")
		typed = ""
		txt.Clear()

		go setSceneSwitchTrueInTime(2 * time.Second)
	}
}

func writeDemoText() {
	title.Color = colornames.White
	footer.Color = colornames.White

	title.Clear()
	title.WriteString("Type in anything and press ENTER!\n\n")
	title.WriteString("CTRL + S: toggle shader\n")

	title.WriteString("CTRL + A: play music\n")
	title.WriteString("CTRL + U, I, O, P: increase volume of music layers\n")
	title.WriteString("CTRL + J, K, L, O-Umlaut (; for QWERTY): decrease volume of individual tracks")

	footer.Clear()
	footer.WriteString("Use the arrow keys to change the background!\n")
}




func drawDemoScene(win *pixelgl.Window, start time.Time) {
	if isSceneSwitch {
		writeDemoText()
	}

	if isShaderApplied {
		updateShader(&uTime, &uSpeed, start)
	}

	win.Clear(bgColor)
	title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, 300)))
	footer.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, -300)))
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}