package main

import (
	"log"
	"time"

	"github.com/3ter/iMagine/scene"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	// import packages purely for their initialization side effects.
	// see https://golang.org/pkg/image
	_ "image/jpeg"
	_ "image/png"
)

func gameloop(win *pixelgl.Window) {
	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience
	start := time.Now()

	scene.SetWindowForAllScenes(win)

	scene.LoadFilesToSceneMap()
	scene.GlobalCurrentScene = `MainMenu`

	for !win.Closed() {

		switch scene.GlobalCurrentScene {

		case "Quit":
			win.SetClosed(true)

		default:
			scene.GlobalScenes[scene.GlobalCurrentScene].OnUpdate(win)
			scene.GlobalScenes[scene.GlobalCurrentScene].Draw(win, start)
		}

		win.Update()
		<-fps
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "iMagine",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true) // remove potential artifacts

	gameloop(win)
}

func main() {
	// to change the flags on the default logger to also print the location (e.g. log.Fatal("Foo"))
	log.SetFlags(log.LstdFlags | log.Llongfile)

	pixelgl.Run(run)
}
