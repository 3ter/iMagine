package scene

import (
	"image/color"
	"time"

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

// Scene contains basic settings and assets (font, music, shaders, content)
type Scene struct {
	bgColor         color.RGBA //= colornames.Black
	fragmentShader  string     // =fileio.LoadFileToString("../assets/wavy_shader.glsl")
	uTime, uSpeed   float32    // pointers to the two uniforms used by fragment shaders
	isShaderApplied bool

	face               font.Face
	txt, title, footer *text.Text
	typed              string

	trackMap      map[int]*effects.Volume
	isSceneSwitch bool
}

// TODO: This has probably been copied here as a reference.
/*
var (
    bgColor         = colornames.Black
    fragmentShader  = fileio.LoadFileToString("../assets/wavy_shader.glsl")
    uTime, uSpeed   float32
    isShaderApplied bool
	isSceneSwitch = true

    face   font.Face
    txt    *text.Text
    title  *text.Text
    footer *text.Text
    typed  string
)
*/

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
	s.isSceneSwitch = true
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
	s.txt = text.New(pixel.V(100, 500), atlas)
	s.title = text.New(pixel.ZV, atlas)
	s.footer = text.New(pixel.ZV, atlas)

	s.trackMap = make(map[int]*effects.Volume)

	s.fragmentShader = fileio.LoadFileToString("../assets/wavy_shader.glsl")
	s.uSpeed = 5.0
	s.isShaderApplied = false

}
