// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"regexp"
	"strings"

	"github.com/3ter/iMagine/internal/fileio"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

// NarratorText includes additional attributes to the implementation by github.com/faiface/text.
// The first being textSpeed for gradually revealed text.
type NarratorText struct {
	*text.Text

	// textSpeed should be measured in CPM (characters per minute)
	// slow: 50, normal: 100, fast: 150, faster: 200
	textSpeed int
}

// Narrator is defined by its text.
//
// Every attributes are private atm so functions to interact with them are
// expected to be created in this package.
type Narrator struct {
	atlas    *text.Atlas
	fontFace font.Face
	color    color.RGBA

	// currentTextObjects contain text objects that define one letter of the current line of the Texter.
	// In the text library the color can be set via its attribute but for changing the font a new object is needed.
	currentTextObjects []*NarratorText

	// currentTextString contains a string stripped from any markup symbols.
	currentTextString string

	textBox *TextBox
}

// SetDefaultAttributes initializes the Player struct
func (n *Narrator) setDefaultAttributes() {
	face, err := fileio.LoadTTF("../assets/intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	n.fontFace = face
	// pixel.ZV is the zero vector representing the orig(in) (i.e. beginning of the line)
	// TODO: see player.go same function
	// n.currentTextObject = text.New(pixel.ZV, text.NewAtlas(face, text.ASCII))
	// n.setTextColor(colornames.Blueviolet)

	n.textBox = new(TextBox)
	// TODO: Find a good way to know the window dimensions here...
	// I've used a potential hack, now we only have to use 'window' here to get relative bounds
	n.textBox.dimensions = pixel.V(900, 230)
	n.textBox.topLeftCorner = pixel.V(1024/2-n.textBox.dimensions.X/2, 768-100)
	n.textBox.thickness = 5
	n.textBox.margin = 20
}

type markdownCommand struct {
	idxStart          int
	idxEnd            int
	attributeValueMap map[string]string
}

// I don't do markdown error checking here (e.g. same number of opening/closing brackets) because it is expected
// to be seen in the markdown preview in an editor.
func getMarkdownCommandSliceFromString(str string) ([]markdownCommand, string) {

	var markdownCommandSlice []markdownCommand
	var strippedStr string

	htmlOpenRegexp := regexp.MustCompile(`\<([^\/].+?)\>`)
	htmlCloseRegexp := regexp.MustCompile(`\<\/\w+\s?\>`)
	styleAttrValRegexp := regexp.MustCompile(`([^":\s]+)\s*:\s*([^;"\s]+)`)

	styleMatchSlice := htmlOpenRegexp.FindAllStringSubmatch(str, -1)
	styleIndexStartSlice := htmlOpenRegexp.FindAllStringIndex(str, -1)
	styleIndexEndSlice := htmlCloseRegexp.FindAllStringIndex(str, -1)

	anyHTMLRegexp := regexp.MustCompile(`\<.+?\>`)
	strippedStr = anyHTMLRegexp.ReplaceAllString(str, "")

	// for adjusting indices after the replacement
	cumulativeOffset := 0

	var currMarkdownCommand markdownCommand
	for i := 0; i < len(styleIndexStartSlice); i++ {

		currMarkdownCommand.attributeValueMap = make(map[string]string)
		currHTMLStyleSlicesSlice := styleAttrValRegexp.FindAllStringSubmatch(styleMatchSlice[i][1], -1)
		for _, HTMLstyleSlice := range currHTMLStyleSlicesSlice {
			currMarkdownCommand.attributeValueMap[HTMLstyleSlice[1]] = HTMLstyleSlice[2]
		}

		currMarkdownCommand.idxStart = styleIndexStartSlice[i][0] - cumulativeOffset
		cumulativeOffset += styleIndexStartSlice[i][1] - styleIndexStartSlice[i][0]

		currMarkdownCommand.idxEnd = styleIndexEndSlice[i][0] - cumulativeOffset
		cumulativeOffset += styleIndexEndSlice[i][1] - styleIndexEndSlice[i][0]

		markdownCommandSlice = append(markdownCommandSlice, currMarkdownCommand)
	}

	return markdownCommandSlice, strippedStr
}

func (n *Narrator) applyMarkdownCommand(markdownCommandSlice []markdownCommand, idx int, scn *Scene) {

	if len(markdownCommandSlice) == 0 {
	} else if idx == markdownCommandSlice[0].idxStart {
		for attribute, value := range markdownCommandSlice[0].attributeValueMap {
			switch attribute {
			case `color`:
				n.color = colornames.Map[strings.ToLower(value)]
			case `font-size`:
				// TODO: Cast value to int accordingly
				face, err := fileio.LoadTTF("../assets/intuitive.ttf", 30)
				if err != nil {
					panic(err)
				}
				n.atlas = text.NewAtlas(face, text.ASCII)
			case `text-speed`:
				// TODO: Implement
			}
		}
	} else if idx == markdownCommandSlice[0].idxEnd {
		markdownCommandSlice = markdownCommandSlice[1:]
		n.atlas = scn.atlas
		n.color = scn.textColor
	}
}

func (n *Narrator) convertMarkdownStringToTextObjectsInBox(str string, scn *Scene) {

	markdownCommandSlice, str := getMarkdownCommandSliceFromString(str)
	n.currentTextString = str

	n.currentTextObjects = nil
	leftIndent := n.textBox.topLeftCorner.X + n.textBox.margin
	nextWordRegexp := regexp.MustCompile(`^[^\s]+ `)
	var nextWord string

	// starting point for writing characters
	currentOrig := pixel.V(leftIndent, n.textBox.topLeftCorner.Y-2*n.textBox.margin)
	n.atlas = scn.atlas
	n.color = scn.textColor

	for idx, rune := range str {

		n.applyMarkdownCommand(markdownCommandSlice, idx, scn)

		char := string(rune)
		switch char {
		case `\n`:
			// align at left indent and remove one line height to the current Y Position
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -n.currentTextObjects[idx].LineHeight))
		case ` `:
			nextWord = nextWordRegexp.FindString(str[(idx + 1):])
		}

		newTextObject := &NarratorText{Text: text.New(currentOrig, n.atlas)}
		n.currentTextObjects = append(n.currentTextObjects, newTextObject)
		newTextObject.Color = n.color

		newTextObject.WriteString(char)
		currentOrig = newTextObject.Dot

		if newTextObject.BoundsOf(nextWord).Max.X >
			(n.textBox.topLeftCorner.X + n.textBox.dimensions.X - 2*n.textBox.margin) {
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -n.currentTextObjects[idx].LineHeight))
		}
		// to not check again until next space indicating the next word
		nextWord = ""
	}
}

// setTextRangeFontFace sets the font face for every letter in the range specified by two indices.
//
// To change the whole string you can use 0 and len(str) as indices.
func (n *Narrator) setTextRangeFontFace(face font.Face, indexStart, indexEnd int) {
	for idx, textObj := range n.currentTextObjects {
		if idx < indexStart {
			continue
		} else if idx >= indexEnd {
			break
		}
		textObj = &NarratorText{Text: text.New(textObj.Orig, text.NewAtlas(face, text.ASCII))}
		// The newly created *text.Text doesn't contain any glyphs to draw yet
		currLetter := string(n.currentTextString[idx])
		textObj.WriteString(currLetter)
	}
}

func (n *Narrator) setTextRangeColor(col color.Color, indexStart, indexEnd int) {
	for idx, textObj := range n.currentTextObjects {
		if idx < indexStart {
			continue
		} else if idx >= indexEnd {
			break
		}
		textObj.Color = col
	}
}

// setText accepts a string with potential markdown formatting containing HTML with inline CSS for text formatting.
// e.g. roses are <span style="color:red">red</span>
// Online LF "\n" is used to mark a new line.
func (n *Narrator) setTextLetterByLetter(str string, scn *Scene) {

	n.convertMarkdownStringToTextObjectsInBox(str, scn)
}

func (n *Narrator) drawTextLetterByLetter(win *pixelgl.Window) {
	n.textBox.drawTextBox(win)

	for _, textObj := range n.currentTextObjects {
		textObj.Draw(win, pixel.IM)
	}
}
