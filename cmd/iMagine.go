// taken from the tutorial here: https://github.com/faiface/pixel/wiki/Typing-text-on-the-screen

package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"

	"github.com/3ter/iMagine/internal/utils"
)

var isMusicPlaying = false
var trackArray = [4]string{"Celesta.ogg", "Choir.ogg", "Harp.ogg", "Strings.ogg"}
var trackPath = "../assets/"

var fragmentShader = utils.LoadFileToString("./assets/wavy_shader.glsl")
var uTime, uSpeed float32

func convertTextToRGB(txt string) [3]uint8 {
	var rgb = [3]uint8{0, 0, 0}

	for pos, char := range txt {
		switch {
		case pos <= 1/3*len(txt):
			rgb[0] = uint8(rgb[0] + uint8(char)*10)
		case pos <= 2/3*len(txt):
			rgb[1] = uint8(rgb[0] + uint8(char)*20)
		case pos <= len(txt):
			rgb[2] = uint8(rgb[0] + uint8(char)*30)
		}
	}

	return rgb
}

func getStreamer(trackpath string) beep.StreamSeekCloser {
	f, err := os.Open(trackpath)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	return streamer
}

func toggleMusic(streamer beep.StreamSeekCloser) {
	if isMusicPlaying {
		speaker.Clear()
		isMusicPlaying = false
	} else {
		speaker.Play(streamer)
		isMusicPlaying = true
	}
}

func applyShader(win *pixelgl.Window, start time.Time) {
	win.Canvas().SetUniform("uTime", &uTime)
	win.Canvas().SetUniform("uSpeed", &uSpeed)
	win.Canvas().SetFragmentShader(fragmentShader)
}

func updateShader(uTime *float32, uSpeed *float32, start time.Time) {
	*uSpeed = 5.0
	*uTime = float32(time.Since(start).Seconds())
}

func gameloop(win *pixelgl.Window) {
	face, err := utils.LoadTTF("./assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(100, 500), atlas)
	title := text.New(pixel.ZV, atlas)
	footer := text.New(pixel.ZV, atlas)

	var typed string
	var bgColor = colornames.Black

	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience

	var trackMap = make(map[int]beep.StreamSeekCloser)
	for index, element := range trackArray {
		fmt.Println(index, trackPath, element)
		var trackStreamer = getStreamer(trackPath + element)
		trackMap[index] = trackStreamer
		defer trackStreamer.Close()
	}

	//var streamer = getStreamer()

	//defer streamer.Close()
	var streamer = utils.GetStreamer("./assets/track1.ogg")
	defer streamer.Close()

	var isShaderApplied = false

	start := time.Now()
	for !win.Closed() {

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
		}
		/*if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyM) {
			toggleMusic(streamer)
		}
		*/
		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
			// TODO: Make it a toggle (set a default fragment shader..?)
			applyShader(win, start)
			isShaderApplied = true
		}

		if isShaderApplied {
			updateShader(&uTime, &uSpeed, start)
		}

		if title.Dot == title.Orig {
			title.WriteString("Type in anything and press ENTER!")
		}
		if footer.Dot == footer.Orig {
			footer.WriteString("Use the arrow keys to change the background!")
		}

		if win.Pressed(pixelgl.KeyDown) {
			bgColor.R--
		}
		if win.Pressed(pixelgl.KeyUp) {
			bgColor.R++
		}

		txt.WriteString(win.Typed())

		typed = typed + win.Typed()

		// b/c GLFW doesn't support {Enter} (and {Tab}) (yet)
		if win.JustPressed(pixelgl.KeyEnter) {
			// txt.WriteRune('\n')
			title.Clear()
			rgb := convertTextToRGB(typed)
			title.Color = color.RGBA{rgb[0], rgb[1], rgb[2], 0xff}
			title.WriteString("That worked quite well! (You can do that again)")
			typed = ""
			txt.Clear()
		}
		// TODO: Add backspace (e.g. use this as reference
		// https://github.com/faiface/pixel-examples/blob/master/typewriter/main.go

		win.Clear(bgColor)
		title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, 300)))
		footer.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, -300)))
		txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
		win.Update()

		// TODO: Understand exactly how this realizes the framerate
		<-fps
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true) // remove potential artifacts

	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}
