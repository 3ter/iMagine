package main

import (
	//"fmt"
	
	"time"
	"github.com/faiface/pixel/pixelgl"

	"github.com/3ter/iMagine/internal/scene/sceneMain"
	"github.com/3ter/iMagine/internal/scene/sceneDemo"

)

var (
	gameState     = "mainMenu"
	prevGameState = ""
	isSceneSwitch = true
)



func setSceneSwitchTrueInTime(duration time.Duration) {
	time.Sleep(duration)
	isSceneSwitch = true
}

func gameloop(win *pixelgl.Window) {
	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience
	start := time.Now()

	for !win.Closed() {

		switch gameState {
		case "Quit":
			win.SetClosed(true)
		case "mainMenu":
			prevGameState = gameState
			main := sceneMain{}
			gameState = main.handleMainMenuAndReturnState(win)
			isSceneSwitch = (gameState != prevGameState)
/*
		case "Start":
			prevGameState = gameState
			handleStartSceneInput(win)
			typeStartTitle()
			drawStartScene(win)
			isSceneSwitch = (gameState != prevGameState)
		case "Demo":
			prevGameState = gameState
			handleDemoInput(win, start)
			drawDemoScene(win, start)
			isSceneSwitch = (gameState != prevGameState)
		}
*/		
		win.Update()
		<-fps
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title: "Pixel Rocks!",
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
