// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"sync"
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

var globalPlayer Player
var globalNarrator Narrator
var globalWindow *pixelgl.Window

type threadSafeBool struct {
	value bool
	sync.Mutex
}

// Scene contains basic settings and assets (font, music, shaders, content)
type Scene struct {
	// Name is the scene identifier which is used in 'OnUpdate' to determine the functions to call
	Name string

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

	trackMap          map[int]*effects.Volume
	IsSceneSwitch     bool
	isPreventInput    threadSafeBool
	isImmediateReveal threadSafeBool

	script        Script
	progress      string
	mapConfigPath string
	mapConfig     *MapConfig
	objects       map[string]map[string]interface{}
}

// Script groups all info from the (markdown) script to make it available to functions within a scene
type Script struct {
	filePath    string
	fileContent string
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
	globalPlayer.setDefaultAttributes()
	globalNarrator.setDefaultAttributes()
}

// SetWindowForAllScenes initializes the global window variable for all scenes
func SetWindowForAllScenes(win *pixelgl.Window) {
	globalWindow = win
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

func (s *Scene) applyShader(win *pixelgl.Window) {
	win.Canvas().SetUniform("uTime", &(s.uTime))
	win.Canvas().SetUniform("uSpeed", &(s.uSpeed))
	win.Canvas().SetFragmentShader(s.fragmentShader)
}

func (s *Scene) clearShader(win *pixelgl.Window) {
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

func getSceneObjectWithDefaults() *Scene {

	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	defaultScene := &Scene{
		bgColor:   colornames.White,
		textColor: colornames.Black,
		atlas:     text.NewAtlas(face, text.ASCII),

		trackMap: make(map[int]*effects.Volume),

		fragmentShader: fileio.LoadFileToString("../assets/wavy_shader.glsl"),
		//TODO: this shader does not do a true passthrough yet and only converts to grayscale
		passthroughShader: fileio.LoadFileToString("../assets/passthrough_shader.glsl"),
		uSpeed:            5.0,
		isShaderApplied:   false,

		progress: "beginning",
	}

	defaultScene.initHintText()

	return defaultScene
}

func handleBackspace(win *pixelgl.Window) {
	if len(globalPlayer.currentTextString) > 0 &&
		(win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace)) {
		globalPlayer.setText(globalPlayer.currentTextString[:len(globalPlayer.currentTextString)-1])
	}
}

func (s *Scene) toggleIsInteractiveUI() {
	if len(s.script.responseQueue) == 0 && len(s.script.keywordResponseMap) > 0 {
		globalPlayer.isInteractiveUI = true
	} else {
		globalPlayer.isInteractiveUI = false
	}
}

func (s *Scene) updateHintTexts() {
	if globalPlayer.isInteractiveUI {
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
func (s *Scene) OnUpdate(win *pixelgl.Window) {

	switch GlobalCurrentScene {
	case "MainMenu":
		s.onUpdateMainMenu(win)
		return
	case `Demo`:
		s.onUpdateDemo(win)
		return
	}

	if s.isPreventInput.value {
		if win.JustPressed(pixelgl.KeySpace) {
			s.isImmediateReveal.Lock()
			s.isImmediateReveal.value = true
			s.isImmediateReveal.Unlock()
		}
		return
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
		win.SetClosed(true)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		GlobalCurrentScene = "MainMenu"
	}
	handleBackspace(win)
	if win.JustPressed(pixelgl.KeyEnter) || (globalPreviousScene != GlobalCurrentScene) {
		globalPreviousScene = GlobalCurrentScene
		if len(s.script.responseQueue) == 0 && len(s.script.keywordResponseMap) == 0 {
			s.parseScriptFile()
		}
		s.executeScriptFromQueue()

		s.toggleIsInteractiveUI()
		s.updateHintTexts()
	}
	globalPlayer.cycleWordInventory(win)

	if len(s.script.responseQueue) == 0 && len(win.Typed()) > 0 {
		globalPlayer.addText(win.Typed(), s)
	}
}

// Draw draws background and text to the window.
func (s *Scene) Draw(win *pixelgl.Window, start time.Time) {

	switch GlobalCurrentScene {
	case `MainMenu`:
		s.drawMainMenu(win)
		return
	case `Demo`:
		s.drawDemo(win, start)
		return
	case `Quit`:
		return
	}

	// TODO: I currently see the scene configs as package variables inside their respective files
	// but the struct initialization in main needs to support this.
	s.bgColor = colornames.White
	win.Clear(s.bgColor)
	s.textColor = colornames.Black

	s.narratorBoxHint.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.narratorBoxHint.Bounds().Center())).Moved(
		pixel.V(0, 2*s.narratorBoxHint.Bounds().H())))
	s.playerBoxHint.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(s.playerBoxHint.Bounds().Center())).Moved(
		pixel.V(0, -5.5*s.playerBoxHint.Bounds().H())))

	globalPlayer.drawTextInBox(win)
	globalNarrator.drawTextInBox(win)

	if globalPlayer.isInteractiveUI {
		globalPlayer.drawWordInventory(win)
	}
}
