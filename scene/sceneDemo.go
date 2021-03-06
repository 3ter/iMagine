// Package scene implements a small demo to test basic capabilities we need
package scene

import (
	"image/color"
	"time"

	"github.com/3ter/iMagine/controlaudio"
	"github.com/3ter/iMagine/controltext"
	"github.com/3ter/iMagine/fileio"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

// initDemo loads all demo specific assets for the scene
func (s *Scene) initDemo() {

	s.IsSceneSwitch = true

	s.bgColor = colornames.Black
	s.textColor = colornames.Wheat

	s.txt = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}
	s.title = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}
	s.footer = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}

	var trackArray = [4]string{"Celesta.ogg", "Choir.ogg", "Harp.ogg", "Strings.ogg"}
	var trackPath = "../assets/"
	s.trackMap = make(map[int]*effects.Volume)
	for index, element := range trackArray {
		//fmt.Println(index, trackPath, element)
		var streamer = fileio.GetStreamer(trackPath + element)
		s.trackMap[index] = streamer
	}
}

// OnUpdateDemo listens and processes player input.
func (s *Scene) onUpdateDemo(win *pixelgl.Window) {

	if s.IsSceneSwitch {
		s.writeDemoText()
	}
	s.IsSceneSwitch = (globalPreviousScene != GlobalCurrentScene)
	globalPreviousScene = GlobalCurrentScene

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		GlobalCurrentScene = "MainMenu"
		return
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyU) {
		controlaudio.VolumeUp(s.trackMap[0])
	}
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyJ) {
		controlaudio.VolumeDown(s.trackMap[0])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyI) {
		controlaudio.VolumeUp(s.trackMap[1])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyK) {
		controlaudio.VolumeDown(s.trackMap[1])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyO) {
		controlaudio.VolumeUp(s.trackMap[2])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
		controlaudio.VolumeDown(s.trackMap[2])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyP) {
		controlaudio.VolumeUp(s.trackMap[3])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeySemicolon) {
		controlaudio.VolumeDown(s.trackMap[3])
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyA) {
		//TODO: This should be a toggle as well.
		allStreamer := beep.Mix(s.trackMap[0], s.trackMap[1], s.trackMap[2], s.trackMap[3])
		speaker.Play(allStreamer)
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {

		if s.isShaderApplied {
			s.clearShader(win)
			s.isShaderApplied = false
		} else {
			s.applyShader(win)
			s.isShaderApplied = true

		}
	}

	if win.Pressed(pixelgl.KeyDown) {
		s.bgColor.R--
	}
	if win.Pressed(pixelgl.KeyUp) {
		s.bgColor.R++
	}

	s.txt.WriteString(win.Typed())
	s.typed += win.Typed()

	if win.JustPressed(pixelgl.KeyEnter) {
		// Mind that GLFW doesn't support {Enter} (and {Tab}) (yet)
		// txt.WriteRune('\n')
		s.title.Clear()
		rgb := convertTextToRGB(s.typed)
		// TODO: Fix convertTextToRGB function (at the moment it's pretty black)
		s.title.Color = color.RGBA{rgb[0], rgb[1], rgb[2], 0xff}
		s.title.WriteString("That worked quite well! (You can do that again)")
		s.typed = ""
		s.txt.Clear()

		go s.setSceneSwitchTrueInTime(2 * time.Second)
	}
}

func (s *Scene) writeDemoText() {

	s.title.Color = colornames.White
	s.footer.Color = colornames.White

	s.title.Clear()
	s.title.WriteString("\n\nSHADER\n")
	s.title.WriteString("CTRL + S: toggle shader\n\n")

	s.title.WriteString("MUSIC\n")
	s.title.WriteString("CTRL + A: play music\n")
	s.title.WriteString("CTRL + U, I, O, P: increase volume of music layers\n")
	s.title.WriteString("CTRL + J, K, L, O-Umlaut (; for QWERTY): decrease volume of individual tracks\n\n")

	s.title.WriteString("TYPING\n")
	s.title.WriteString("Type in anything and press ENTER!\n")

	s.footer.Clear()

	s.footer.WriteString("BG COLOR\n")
	s.footer.WriteString("Use the UP and DOWN arrow keys to change the background!\n\n")
	s.footer.WriteString("REVEALED TEXT\n")

	writingDoneChannel := make(chan int)
	var revealedText = "Here is some gradually revealed text.\n"
	go controltext.WriteToTextLetterByLetter(s.footer, revealedText, 60, writingDoneChannel)
	writingDoneChannel <- 1 // init writing the first line

	revealedText = "Fast reveal, quite thrilling. \n"
	go controltext.WriteToTextLetterByLetter(s.footer, revealedText, 10, writingDoneChannel)
	revealedText = "This text will be revealed slooooooooooooooooooooowly."
	go controltext.WriteToTextLetterByLetter(s.footer, revealedText, 100, writingDoneChannel)

}

// DrawDemo draws background and text to the window.
func (s *Scene) drawDemo(win *pixelgl.Window, start time.Time) {

	if s.isShaderApplied {
		s.updateShader(s.uSpeed, start)
	}

	win.Clear(s.bgColor)
	s.title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.title.Bounds().Center())).Moved(pixel.V(0, 250)))
	s.footer.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.footer.Bounds().Center())).Moved(pixel.V(0, -150)))
	s.txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.txt.Bounds().Center())).Moved(pixel.V(0, 50)))
}
