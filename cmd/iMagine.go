package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/3ter/iMagine/internal/utils"
)

var gameState = "mainMenu"
var isMusicPlaying = false
var trackArray = [4]string{"Celesta.ogg", "Choir.ogg", "Harp.ogg", "Strings.ogg"}
var trackPath = "../assets/"
var bgColor = colornames.Black

var fragmentShader = utils.LoadFileToString("../assets/wavy_shader.glsl")
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

func toggleMusic(streamer beep.StreamSeekCloser) {

	speaker.Play(streamer)

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

type mainMenuItem struct {
	Text  string
	State string
}

var menuItems = []*mainMenuItem{
	&mainMenuItem{"Demo", "selected"},
	&mainMenuItem{"Quit", "unselected"}}

func returnMenuTexts(atlasRegular, atlasBold *text.Atlas) []*text.Text {
	menuTexts := make([]*text.Text, 2)
	for i, menuItem := range menuItems {
		txt := text.New(pixel.ZV, atlasRegular)
		if menuItem.State == "selected" {
			txt = text.New(pixel.ZV, atlasBold)
		}
		txt.WriteString(menuItem.Text)
		menuTexts[i] = txt
	}

	return menuTexts
}

func drawMainMenu(win *pixelgl.Window, atlasRegular, atlasBold *text.Atlas) {
	win.Clear(bgColor)

	menuTextVerticalOffset := 50 // pixels

	menuTexts := returnMenuTexts(atlasRegular, atlasBold)
	for i, menuText := range menuTexts {
		centerTextMatrix := pixel.IM.Moved(win.Bounds().Center().Sub(menuText.Bounds().Center()))
		verticalAdjustVector := pixel.V(0, float64(-menuTextVerticalOffset*i))
		menuText.Draw(win, centerTextMatrix.Moved(verticalAdjustVector))
	}
}

func handleMainMenuAndReturnState(win *pixelgl.Window) string {
	regularFace := utils.TtfFromBytesMust(goregular.TTF, 20)
	boldFace := utils.TtfFromBytesMust(gobold.TTF, 20)
	atlasRegular := text.NewAtlas(regularFace, text.ASCII)
	atlasBold := text.NewAtlas(boldFace, text.ASCII)
	drawMainMenu(win, atlasRegular, atlasBold)

	if win.JustPressed(pixelgl.KeyDown) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i < len(menuItems)-1 {
				menuItem.State = "unselected"
				menuItems[i+1].State = "selected"
			}
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i > 0 {
				menuItem.State = "unselected"
				menuItems[i-1].State = "selected"
			}
		}
	}

	if win.JustPressed(pixelgl.KeyEnter) {
		for _, menuItem := range menuItems {
			if menuItem.State == "selected" {
				return menuItem.Text
			}
		}
	}

	return "mainMenu"
}

func gameloop(win *pixelgl.Window) {
	face, err := utils.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(100, 500), atlas)
	title := text.New(pixel.ZV, atlas)
	footer := text.New(pixel.ZV, atlas)

	var typed string

	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience

	var trackMap = make(map[int]*effects.Volume)
	for index, element := range trackArray {
		fmt.Println(index, trackPath, element)
		var streamer = utils.GetStreamer(trackPath + element)
		trackMap[index] = streamer
		//defer streamer.Close()
	}

	var isShaderApplied = false

	start := time.Now()
	for !win.Closed() {

		switch gameState {
		case "Quit":
			win.SetClosed(true)
		case "mainMenu":
			gameState = handleMainMenuAndReturnState(win)
		case "Demo":
			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
				win.SetClosed(true)
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				gameState = "mainMenu"
			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyU) {
				//toggleMusic(trackMap[0])
				trackMap[0].Volume += 0.5
			}
			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyJ) {
				//toggleMusic(trackMap[0])
				trackMap[0].Volume -= 0.5
			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyI) {
				//toggleMusic(trackMap[1])
				trackMap[1].Volume += 0.5
			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyK) {
				//toggleMusic(trackMap[1])
				trackMap[1].Volume -= 0.5
			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyO) {
				trackMap[2].Volume += 0.5
				//toggleMusic(trackMap[2])

			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
				trackMap[2].Volume -= 0.5
				//toggleMusic(trackMap[2])

			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyP) {
				trackMap[3].Volume += 0.5
				//toggleMusic(trackMap[2])

			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeySemicolon) {
				trackMap[3].Volume -= 0.5
				//toggleMusic(trackMap[2])

			}

			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyP) {
				//toggleMusic(trackMap[3])
			}

			if win.JustPressed(pixelgl.KeyA) {
				/*toggleMusic(trackMap[0])
				toggleMusic(trackMap[1])
				toggleMusic(trackMap[2])
				toggleMusic(trackMap[3])
				*/
				allStreamer := beep.Mix(trackMap[0], trackMap[1], trackMap[2], trackMap[3])

				speaker.Play(allStreamer)
			}

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
		}
		win.Update()
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
