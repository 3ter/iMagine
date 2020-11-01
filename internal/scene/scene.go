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
	uTime, uSpeed   float32
	isShaderApplied bool

	face               font.Face
	txt, title, footer *text.Text
	typed              string

	trackMap    map[int]*effects.Volume
	sceneSwitch bool
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
	s.sceneSwitch = true
}

func toggleMusic(streamer beep.StreamSeekCloser) {

	speaker.Play(streamer)

}

// TODO: Add the commands back in.
func applyShader(win *pixelgl.Window, start time.Time) {
	//win.Canvas().SetUniform("uTime", &uTime)
	//win.Canvas().SetUniform("uSpeed", &uSpeed)
	//win.Canvas().SetFragmentShader(fragmentShader)
}

func updateShader(uTime *float32, uSpeed *float32, start time.Time) {
	*uSpeed = 5.0
	*uTime = float32(time.Since(start).Seconds())
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

	// TODO: This is probably here because there was the intent to make this generally available vs
	// declaring it in every scene anew.
	/*
		for index, element := range s.trackArray {
			fmt.Println(index, trackPath, element)
			var streamer = fileio.GetStreamer(trackPath + element)
			s.trackMap[index] = streamer

			//TODO: Why is this commented out?
			//defer streamer.Close()
		}
	*/

	//TODO: Apply shader
	//isShaderApplied = false
}
