//Main Scene
package scene

import "github.com/3ter/iMagine/internal/scene/Scene"

func drawMainMenu(win *pixelgl.Window, atlasRegular, atlasBold *text.Atlas) {
	bgColor = colornames.Black
	win.Clear(bgColor)

	menuTextVerticalOffset := 50 // pixels

	menuTexts := returnMenuTexts(atlasRegular, atlasBold)
	for i, menuText := range menuTexts {
		centerTextMatrix := pixel.IM.Moved(win.Bounds().Center().Sub(menuText.Bounds().Center()))
		verticalAdjustVector := pixel.V(0, float64(-menuTextVerticalOffset*i))
		menuText.Draw(win, centerTextMatrix.Moved(verticalAdjustVector))
	}
}

func handleMainMenuAndReturnState(win *pixelgl.Window) string {

	regularFace := fileio.TtfFromBytesMust(goregular.TTF, 20)
	boldFace := fileio.TtfFromBytesMust(gobold.TTF, 20)
	atlasRegular := text.NewAtlas(regularFace, text.ASCII)
	atlasBold := text.NewAtlas(boldFace, text.ASCII)
	drawMainMenu(win, atlasRegular, atlasBold)

	if win.JustPressed(pixelgl.KeyDown) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i < len(menuItems)-1 {
				menuItem.State = "unselected"
				menuItems[i+1].State = "selected"
				break
			}
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		for i, menuItem := range menuItems {
			if menuItem.State == "selected" && i > 0 {
				menuItem.State = "unselected"
				menuItems[i-1].State = "selected"
				break
			}
		}
	}

	if win.JustPressed(pixelgl.KeyEnter) {
		for _, menuItem := range menuItems {
			if menuItem.State == "selected" {
				return menuItem.Text
			}
		}
	}

	return "mainMenu"
}

type mainMenuItem struct {
	Text  string
	State string
}

var menuItems = []*mainMenuItem{
	&mainMenuItem{"Demo", "selected"},
	&mainMenuItem{"Start", "unselected"},
	&mainMenuItem{"Quit", "unselected"}}

func returnMenuTexts(atlasRegular, atlasBold *text.Atlas) []*text.Text {
	menuTexts := make([]*text.Text, len(menuItems))
	for i, menuItem := range menuItems {
		txt := text.New(pixel.ZV, atlasRegular)
		if menuItem.State == "selected" {
			txt = text.New(pixel.ZV, atlasBold)
		}
		txt.WriteString(menuItem.Text)
		menuTexts[i] = txt
	}

	return menuTexts
}