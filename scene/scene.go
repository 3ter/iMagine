// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"time"

	"golang.org/x/image/colornames"

	"github.com/3ter/iMagine/controltext"

	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/font"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"

	// TODO: These are still here if stuff from other scenes should move here.
	// These specifically seem to be more useful in 'controltext' I suppose.
	//"golang.org/x/image/font/gofont/gobold"
	//"golang.org/x/image/font/gofont/goregular"

	//"github.com/3ter/iMagine/controlaudio"
	"github.com/3ter/iMagine/fileio"
)

var player Player
var narrator Narrator
var window *pixelgl.Window

// Scene contains basic settings and assets (font, music, shaders, content)
type Scene struct {
	bgColor           color.RGBA //= colornames.Black
	fragmentShader    string     // =fileio.LoadFileToString("../assets/wavy_shader.glsl")
	passthroughShader string
	uTime, uSpeed     float32 // pointers to the two uniforms used by fragment shaders
	isShaderApplied   bool

	face      *font.Face
	atlas     *text.Atlas
	textColor color.RGBA

	// TODO: These should probably go or be used for real.
	txt, title, footer *controltext.SafeText

	// hints are used to provide the player with subtle help messages on screen.
	narratorBoxHint *controltext.SafeText
	playerBoxHint   *controltext.SafeText
	typed           string

	trackMap       map[int]*effects.Volume
	IsSceneSwitch  bool
	isPreventInput bool

	script   Script
	progress string
}

// Script groups all info from the (markdown) script to make it available to functions within a scene
type Script struct {
	file string
	// responseQueue contains the responses that still need to be delivered before player commands become active again.
	responseQueue []narratorResponse
	// keywordResponseMap contains a map from the player commands that are understood to a slice of narratorResponses.
	keywordResponseMap map[string][]narratorResponse
}

// narratorResponse groups the narrator text with it's ambience commands
// A player cmd is mapped onto a struct containing narrator text lines, ambience directives or progress updates.
type narratorResponse struct {
	narratorTextLine string
	progressUpdate   string
	ambienceCmdSlice []string
}

// This is called once when the package is imported for the first time
func init() {
	player.setDefaultAttributes()
	narrator.setDefaultAttributes()
}

// SetWindowForAllScenes initializes the global window variable for all scenes
func SetWindowForAllScenes(win *pixelgl.Window) {
	window = win
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

func (s *Scene) setSceneSwitchTrueInTime(duration time.Duration) {
	time.Sleep(duration)
	s.IsSceneSwitch = true
}

func toggleMusic(streamer beep.StreamSeekCloser) {

	speaker.Play(streamer)

}

func (s *Scene) applyShader(win *pixelgl.Window, start time.Time) {
	win.Canvas().SetUniform("uTime", &(s.uTime))
	win.Canvas().SetUniform("uSpeed", &(s.uSpeed))
	win.Canvas().SetFragmentShader(s.fragmentShader)
}

func (s *Scene) clearShader(win *pixelgl.Window, start time.Time) {
	win.Canvas().SetUniform("uTime", &(s.uTime))
	win.Canvas().SetUniform("uSpeed", &(s.uSpeed))
	win.Canvas().SetFragmentShader(s.passthroughShader)
}

func (s *Scene) updateShader(uSpeed float32, start time.Time) {
	s.uSpeed = uSpeed
	s.uTime = float32(time.Since(start).Seconds())
}

func (s *Scene) initHintText() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 18)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)

	s.narratorBoxHint = &controltext.SafeText{
		Text: text.New(pixel.ZV, atlas),
	}
	s.narratorBoxHint.Color = colornames.Gray
	s.narratorBoxHint.WriteString("Press Enter to continue.")

	s.playerBoxHint = &controltext.SafeText{
		Text: text.New(pixel.ZV, atlas),
	}
	s.playerBoxHint.Color = colornames.Gray
}

// Init loads text and music into the Scene struct.
func (s *Scene) Init() {
	s.bgColor = colornames.Black
	s.textColor = colornames.White

	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	s.atlas = text.NewAtlas(face, text.ASCII)
	s.txt = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}
	s.title = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}
	s.footer = &controltext.SafeText{
		Text: text.New(pixel.ZV, s.atlas),
	}
	s.initHintText()

	s.trackMap = make(map[int]*effects.Volume)

	s.fragmentShader = fileio.LoadFileToString("../assets/wavy_shader.glsl")
	//TODO: this shader does not do a true passthrough yet and only converts to grayscale
	s.passthroughShader = fileio.LoadFileToString("../assets/passthrough_shader.glsl")
	s.uSpeed = 5.0
	s.isShaderApplied = false

	s.IsSceneSwitch = true

	s.progress = "beginning"
}

// InitWithFile initializes a scene using a scene script file which then should be parsed.
func (s *Scene) InitWithFile(scriptFilepath string) {
	s.Init()
	s.script.file = fileio.LoadFileToString(scriptFilepath)
}

// TODO: redo backspace using the 'Repeating' event (see faiface/pixel Wiki for writing texts)
// handleBackspace is necessary to implement manually as we currently "misuse" the text library in having one text
// object holding all our text so it is currently replaced entirely though only one character should vanish.
func handleBackspace(win *pixelgl.Window, player *Player) {
	if win.JustPressed(pixelgl.KeyBackspace) && len(player.currentTextString) > 0 {
		player.setText(player.currentTextString[:len(player.currentTextString)-1])
		backspaceCounter = int(-120 * 0.5) // Framerate times seconds to wait until continuous backspace kicks in.
	} else if win.Pressed(pixelgl.KeyBackspace) && len(player.currentTextString) > 0 {
		backspaceCounter++
		backspaceDeletionSpeed := int(120 / 40) // Framerate divided by deletions per second.
		if backspaceCounter > 0 && backspaceCounter%backspaceDeletionSpeed == 0 {
			player.setText(player.currentTextString[:len(player.currentTextString)-1])
			backspaceCounter = 0
		}
	}
}

func (s *Scene) updateHintTexts() {
	if len(s.script.responseQueue) == 0 && len(s.script.keywordResponseMap) > 0 {
		s.playerBoxHint.Clear()
		s.narratorBoxHint.Clear()
		s.playerBoxHint.WriteString("Write a command and press Enter.")
	} else {
		s.narratorBoxHint.Clear()
		s.playerBoxHint.Clear()
		s.narratorBoxHint.WriteString("Press Enter to continue.")
	}
}

// OnUpdate listens and processes player input on every frame update.
func (s *Scene) OnUpdate(win *pixelgl.Window, gameState string) string {
	if s.isPreventInput {
		return gameState
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameState = "mainMenu"
	}
	handleBackspace(win, &player)
	if win.JustPressed(pixelgl.KeyEnter) {
		if len(s.script.responseQueue) == 0 && len(s.script.keywordResponseMap) == 0 {
			s.parseScriptFile()
		}
		s.executeScriptFromQueue()

		s.updateHintTexts()
	}

	if len(win.Typed()) > 0 {
		player.addText(win.Typed(), s)
	}

	return gameState
}

// Draw draws background and text to the window.
func (s *Scene) Draw(win *pixelgl.Window) {

	// TODO: I currently see the scene configs as package variables inside their respective files
	// but the struct initialization in main needs to support this.
	s.bgColor = getBeachBackgroundColor()
	win.Clear(s.bgColor)
	s.textColor = colornames.Black

	s.title.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.title.Bounds().Center())).Moved(pixel.V(0, 300)))
	s.narratorBoxHint.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.narratorBoxHint.Bounds().Center())).Moved(
		pixel.V(0, 2*s.narratorBoxHint.Bounds().H())))
	s.playerBoxHint.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.playerBoxHint.Bounds().Center())).Moved(
		pixel.V(0, -5.5*s.playerBoxHint.Bounds().H())))

	player.drawTextInBox(win)
	narrator.drawTextInBox(win)
}
