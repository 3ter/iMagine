package scene

import (
	"image/color"
	"time"

	"github.com/3ter/iMagine/internal/controltext"

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

	//"github.com/3ter/iMagine/internal/controlaudio"
	"github.com/3ter/iMagine/internal/fileio"
)

var player Player
var narrator Narrator
var window *pixelgl.Window

// Scene contains basic settings and assets (font, music, shaders, content)
type Scene struct {
	bgColor         color.RGBA //= colornames.Black
	fragmentShader  string     // =fileio.LoadFileToString("../assets/wavy_shader.glsl")
	uTime, uSpeed   float32    // pointers to the two uniforms used by fragment shaders
	isShaderApplied bool

	face               font.Face
	txt, title, footer *controltext.SafeText
	typed              string

	trackMap      map[int]*effects.Volume
	IsSceneSwitch bool

	scriptFile    string
	progress string
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

func (s *Scene) updateShader(uSpeed float32, start time.Time) {
	s.uSpeed = uSpeed
	s.uTime = float32(time.Since(start).Seconds())
}

// Init loads text and music into the Scene struct.
func (s *Scene) Init() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	s.txt = &controltext.SafeText{
		Text: text.New(pixel.ZV, atlas),
	}
	s.title = &controltext.SafeText{
		Text: text.New(pixel.ZV, atlas),
	}
	s.footer = &controltext.SafeText{
		Text: text.New(pixel.ZV, atlas),
	}

	s.trackMap = make(map[int]*effects.Volume)

	s.fragmentShader = fileio.LoadFileToString("../assets/wavy_shader.glsl")
	s.uSpeed = 5.0
	s.isShaderApplied = false

	s.IsSceneSwitch = true

	s.progress = "beginning"
}

// InitWithFile initializes a scene using a scene script file which then should be parsed.
func (s *Scene) InitWithFile(scriptFilepath string) {
	s.Init()
	s.scriptFile = fileio.LoadFileToString(scriptFilepath)
}
