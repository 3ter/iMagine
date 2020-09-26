// taken from the tutorial here: https://github.com/faiface/pixel/wiki/Typing-text-on-the-screen

package main

import (
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"log"
	"os"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("track1.ogg")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)

}
