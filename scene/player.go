// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"regexp"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"

	"github.com/3ter/iMagine/fileio"
)

type WordInventory struct {
	textObjects []*text.Text
	// index of the currently displayed word
	currentIndex int
}

// Player is defined by its text and contains the game progression in
// the form of its inventory.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Player struct {
	wordInventory   WordInventory
	isInteractiveUI bool

	atlas    *text.Atlas
	fontFace font.Face
	color    color.RGBA

	// currentTextObjects contain text objects that define one letter of the current line of the Texter.
	// In the text library the color can be set via its attribute but for changing the font a new object is needed.
	currentTextObjects []*text.Text

	// currentTextString contains a string stripped from any markup symbols.
	currentTextString string

	textBox *TextBox
}

func (p *Player) setTextFontFace(face font.Face) {
	textObject := p.currentTextObjects[0]
	textObject = text.New(textObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	textObject.WriteString(p.currentTextString)
}

func (p *Player) setTextColor(col color.Color) {
	p.currentTextObjects[0].Color = col
}

// getWrappedString executed on a Texter returns a string wrapped inside the Texter's box.
// This currently depends on the first text object in Texter's currentTextObjects.
func (p *Player) getWrappedString(str string) string {
	maxTextWidth := p.textBox.dimensions.X - p.textBox.margin - 20
	textObject := p.currentTextObjects[0]

	m := regexp.MustCompile(`\n`) // GoogleDocs uses LF to mark its line endings
	currLinesSlice := m.Split(str, -1)
	var wrappedString string
	for idx, currLine := range currLinesSlice {
		if textObject.BoundsOf(str).W() <= maxTextWidth {
			wrappedString = str
			break
		}
		if textObject.BoundsOf(currLine).W() <= maxTextWidth {
			if idx < len(currLinesSlice)-1 {
				wrappedString += currLine + "\n"
			} else {
				wrappedString += currLine
			}
			continue
		}
		newLineBreakIndex := int((maxTextWidth) * (float64(len(currLine))) / textObject.BoundsOf(currLine).W())
		lastSpaceMatch := regexp.MustCompile(` [^ ]*?$`)
		if lastSpaceMatch.FindStringIndex(currLine[:newLineBreakIndex]) != nil {
			newLineBreakIndex = lastSpaceMatch.FindStringIndex(currLine[:newLineBreakIndex])[0]
		}
		if idx < len(currLinesSlice)-1 {
			wrappedString += currLine[:newLineBreakIndex] + "\n" + currLine[newLineBreakIndex+1:] + " "
		} else {
			wrappedString += currLine[:newLineBreakIndex] + "\n" + currLine[newLineBreakIndex+1:]
		}
	}

	matchOnlyWhitespaceString := regexp.MustCompile(`^\s+$`)
	if matchOnlyWhitespaceString.MatchString(wrappedString) {
		wrappedString = ""
	}

	return wrappedString
}

func (p *Player) setText(str string) {
	wrappedString := p.getWrappedString(str)

	p.currentTextString = wrappedString
	p.currentTextObjects[0].Clear()
	p.currentTextObjects[0].WriteString(wrappedString)
}

func (p *Player) addText(str string, scn *Scene) {

	wrappedString := p.getWrappedString(p.currentTextString + str)

	p.currentTextString = wrappedString
	p.currentTextObjects[0].Clear()
	p.currentTextObjects[0].WriteString(wrappedString)
}

func (p *Player) drawTextInBox(win *pixelgl.Window) {
	p.textBox.drawTextBox(win)

	// TODO: The y coordinate is guesswork and probably dependend on the font face used!
	marginVec := pixel.V(p.textBox.margin, p.textBox.margin-p.textBox.thickness-2.4*p.currentTextObjects[0].LineHeight)
	p.currentTextObjects[0].Draw(win, pixel.IM.Moved(p.textBox.topLeftCorner.Add(marginVec)))
}

func (p *Player) drawWordInventory(win *pixelgl.Window) {
	wordBankTextBox := new(TextBox)

	wordBankTextBox.dimensions = pixel.V(130, 30)
	wordBankTextBox.topLeftCorner = pixel.V(1024/2-p.textBox.dimensions.X/2, 768-500+50)
	wordBankTextBox.thickness = 5
	wordBankTextBox.margin = 20

	wordBankTextBox.drawTextBox(win)

	marginVec := pixel.V(wordBankTextBox.margin,
		wordBankTextBox.margin-wordBankTextBox.thickness-1.8*
			p.wordInventory.textObjects[p.wordInventory.currentIndex].LineHeight)
	p.wordInventory.textObjects[p.wordInventory.currentIndex].Draw(win,
		pixel.IM.Moved(wordBankTextBox.topLeftCorner.Add(marginVec)))

	upArrowPic, err := fileio.LoadPicture(`../assets/kb_up_arrow.png`)
	if err != nil {
		panic(err)
	}
	downArrowPic, err := fileio.LoadPicture(`../assets/kb_down_arrow.png`)
	if err != nil {
		panic(err)
	}
	upArrowSprite := pixel.NewSprite(upArrowPic, upArrowPic.Bounds())
	downArrowSprite := pixel.NewSprite(downArrowPic, downArrowPic.Bounds())
	spriteMatrix := pixel.IM.Scaled(pixel.ZV, 0.2)
	spriteMatrix = spriteMatrix.Moved(wordBankTextBox.topLeftCorner.Add(marginVec).Add(pixel.V(wordBankTextBox.dimensions.X, 0)))
	upArrowSprite.Draw(win, spriteMatrix.Moved(pixel.V(0, 17)))
	downArrowSprite.Draw(win, spriteMatrix.Moved(pixel.V(0, -7)))
}

func (p *Player) initWordInventory(face font.Face, col color.Color) {
	firstWordObject := text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	firstWordObject.Color = col
	firstWordObject.WriteString("inspect")

	// FIXME: remove secondWordObject
	secondWordObject := text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	secondWordObject.Color = col
	secondWordObject.WriteString("foobar")

	p.wordInventory = WordInventory{currentIndex: 0}
	// FIXME: remove secondWordObject
	p.wordInventory.textObjects = append(p.wordInventory.textObjects, firstWordObject, secondWordObject)
}

func (p *Player) cycleWordInventory(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyUp) {
		globalPlayer.wordInventory.currentIndex =
			(globalPlayer.wordInventory.currentIndex + 1) % len(globalPlayer.wordInventory.textObjects)
	} else if win.JustPressed(pixelgl.KeyDown) {
		globalPlayer.wordInventory.currentIndex =
			(globalPlayer.wordInventory.currentIndex - 1) % len(globalPlayer.wordInventory.textObjects)
		if globalPlayer.wordInventory.currentIndex < 0 {
			globalPlayer.wordInventory.currentIndex *= -1
		}
	}
}

// SetDefaultAttributes initializes the Player struct
func (p *Player) setDefaultAttributes() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	p.fontFace = face

	// pixel.ZV is the zero vector representing the orig(in) (i.e. beginning of the line)
	p.currentTextObjects = append(p.currentTextObjects, text.New(pixel.ZV, text.NewAtlas(face, text.ASCII)))
	p.setTextColor(colornames.Blueviolet)

	p.textBox = new(TextBox)
	// TODO: Find a good way to know the window dimensions here...
	// I've used a potential hack, now we only have to use 'window' here to get relative bounds
	p.textBox.dimensions = pixel.V(900, 230)
	p.textBox.topLeftCorner = pixel.V(1024/2-p.textBox.dimensions.X/2, 768-500)
	p.textBox.thickness = 5
	p.textBox.margin = 20

	p.initWordInventory(face, colornames.Blueviolet)
}
