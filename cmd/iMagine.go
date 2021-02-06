package main

import (
	"log"
	"time"

	"github.com/3ter/iMagine/scene"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	// import package purely for its initialization side effects.
	// see https://golang.org/pkg/image
	_ "image/jpeg"
	_ "image/png"
)

var (
	gameState     = "mainMenu"
	prevGameState = ""
)

func gameloop(win *pixelgl.Window) {
	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience
	start := time.Now()

	scene.SetWindowForAllScenes(win)

	// TODO: Remove old way to change scenes
	var demoScene = scene.DemoScene
	var mainScene = scene.MainScene
	demoScene.Init()
	demoScene.InitDemoScene()
	mainScene.Init()

	scene.LoadFilesToSceneMap()
	scene.CurrentScene = `Beach`

	for !win.Closed() {

		switch scene.CurrentScene {

		case "Quit":
			win.SetClosed(true)

		case "mainMenu":
			prevGameState = gameState
			gameState = mainScene.HandleMainMenuAndReturnState(win)

		case "Demo":
			prevGameState = gameState
			gameState = demoScene.HandleDemoInput(win, start)
			demoScene.DrawDemoScene(win, start)
			demoScene.IsSceneSwitch = (gameState != prevGameState)

		default:
			gameState = scene.ScenesMap[scene.CurrentScene].OnUpdate(win, gameState)
			scene.ScenesMap[scene.CurrentScene].Draw(win)
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
