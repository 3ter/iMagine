package main

import (
	//"fmt"
	
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/3ter/iMagine/internal/scene"
	

)

var (
	gameState     = "mainMenu"
	prevGameState = ""
)



func gameloop(win *pixelgl.Window) {
	fps := time.Tick(time.Second / 120) // 120 FPS provide a very smooth typing experience
	start := time.Now()

	var isSceneSwitch = true

	var demoScene = scene.DemoScene
	var beachScene = scene.BeachScene
	var mainScene = scene.MainScene
	var forestScene = scene.ForestScene
	demoScene.Init()
	demoScene.InitDemoScene()
	beachScene.Init()
	mainScene.Init()
	forestScene.Init()
	//scene.DemoScene.Init()
	//scene.BeachScene.Init()

	for !win.Closed() {

		switch gameState {
		case "Quit":
			win.SetClosed(true)
		case "mainMenu":
			prevGameState = gameState

			gameState = mainScene.HandleMainMenuAndReturnState(win)
			isSceneSwitch = (gameState != prevGameState)

		case "Start":
			prevGameState = gameState
			gameState = beachScene.HandleBeachSceneInput(win, gameState)
			if(isSceneSwitch){
				beachScene.TypeBeachTitle()
			}
			beachScene.DrawBeachScene(win)
			isSceneSwitch = (gameState != prevGameState)


		case "Forest":
			prevGameState = gameState
			gameState = forestScene.HandleForestSceneInput(win, gameState)
			forestScene.TypeForestTitle()
			forestScene.DrawForestScene(win)
			isSceneSwitch = (gameState != prevGameState)

		case "Demo":
			prevGameState = gameState
			gameState = demoScene.HandleDemoInput(win, start)
			demoScene.DrawDemoScene(win, start)
			isSceneSwitch = (gameState != prevGameState)


			
		}
		
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
