package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/3ter/iMagine/internal/scene"

	"golang.org/x/image/font"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/3ter/iMagine/internal/controlaudio"
	"github.com/3ter/iMagine/internal/fileio"
)

var (
	gameState = "mainMenu"

	isMusicPlaying = false
	trackArray     = [4]string{"Celesta.ogg", "Choir.ogg", "Harp.ogg", "Strings.ogg"}
	trackPath      = "../assets/"
	trackMap       map[int]*effects.Volume

	bgColor         = colornames.Black
	fragmentShader  = fileio.LoadFileToString("../assets/wavy_shader.glsl")
	uTime, uSpeed   float32
	isShaderApplied bool

	face   font.Face
	txt    *text.Text
	title  *text.Text
	footer *text.Text
	typed  string
)

func addStaticText() {
	// Add text only if it is empty
	if title.Dot == title.Orig {
		title.WriteString("Type in anything and press ENTER!\n\n")
		title.WriteString("CTRL + S: toggle shader\n")

		title.WriteString("CTRL + A: play music\n")
		title.WriteString("CTRL + U, I, O, P: increase volume of music layers\n")
		title.WriteString("CTRL + J, K, L, O-Umlaut (; for QWERTY): decrease volume of individual tracks")

	}
	if footer.Dot == footer.Orig {
		footer.WriteString("Use the arrow keys to change the background!\n")
	}
}

func init() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt = text.New(pixel.V(100, 500), atlas)
	title = text.New(pixel.ZV, atlas)
	footer = text.New(pixel.ZV, atlas)

	trackMap = make(map[int]*effects.Volume)
	for index, element := range trackArray {
		fmt.Println(index, trackPath, element)
		var streamer = fileio.GetStreamer(trackPath + element)
		trackMap[index] = streamer

		//TODO: Why is this commented out?
		//defer streamer.Close()
	}

	isShaderApplied = false

	addStaticText()
}

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
	&mainMenuItem{"Start", "unselected"},
	&mainMenuItem{"Quit", "unselected"}}

func returnMenuTexts(atlasRegular, atlasBold *text.Atlas) []*text.Text {
	menuTexts := make([]*text.Text, len(menuItems))
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
	bgColor = colornames.Black
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
	regularFace := fileio.TtfFromBytesMust(goregular.TTF, 20)
	boldFace := fileio.TtfFromBytesMust(gobold.TTF, 20)
	atlasRegular := text.NewAtlas(regularFace, text.ASCII)
	atlasBold := text.NewAtlas(boldFace, text.ASCII)
	drawMainMenu(win, atlasRegular, atlasBold)

	if win.JustPressed(pixelgl.KeyDown) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i < len(menuItems)-1 {
				menuItem.State = "unselected"
				menuItems[i+1].State = "selected"
				break
			}
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i > 0 {
				menuItem.State = "unselected"
				menuItems[i-1].State = "selected"
				break
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
	}
}

func drawDemoScene(win *pixelgl.Window, start time.Time) {
	if isShaderApplied {
		updateShader(&uTime, &uSpeed, start)
	}

	bgColor = colornames.Black
	win.Clear(bgColor)
	title.Color = colornames.White
	title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, 300)))
	footer.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(title.Bounds().Center())).Moved(pixel.V(0, -300)))
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}

func handleStartSceneInput(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
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

func gameloop(win *pixelgl.Window) {
	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience
	start := time.Now()

	for !win.Closed() {

		switch gameState {
		case "Quit":
			win.SetClosed(true)
		case "mainMenu":
			gameState = handleMainMenuAndReturnState(win)
		case "Start":
			handleStartSceneInput(win)
			typeStartTitle()
			drawStartScene(win)
		case "Demo":
			handleDemoInput(win, start)
			drawDemoScene(win, start)
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
