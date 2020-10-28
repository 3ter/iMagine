package scene
 
import (
    //"image/color"
    //"fmt"
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



type Scene struct {
	
    bgColor         = colornames.Black
    fragmentShader  = fileio.LoadFileToString("../assets/wavy_shader.glsl")
    uTime, uSpeed   float32
    isShaderApplied bool

    face   font.Face
    txt    *text.Text
    title  *text.Text
    footer *text.Text
    typed  string

}

/*
var (
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
*/

func (s *Scene) convertTextToRGB(txt string) [3]uint8 {
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

func (s *Scene) toggleMusic(streamer beep.StreamSeekCloser) {

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


 /*
func (s *Scene) init() int {
    return p.Sides
}
 
type Triangle struct {
    Polygon // anonymous field
}
 */



func (s *Scene) init() {
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
}

