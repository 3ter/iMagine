// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"regexp"
	"strings"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// Texter is defined by its text
type Texter struct {
	fontFace font.Face

	// currentTextObjects contain text objects that define one letter of the current line of the Texter.
	// In the text library the color can be set via its attribute but for changing the font a new object is needed.
	currentTextObjects []*text.Text

	// currentTextString contains a string stripped from any markup symbols.
	currentTextString string

	textBox *TextBox
}

// setTextRangeFontFace sets the font face for every letter in the range specified by two indices.
//
// To change the whole string you can use 0 and len(str) as indices.
func (t *Texter) setTextRangeFontFace(face font.Face, indexStart, indexEnd int) {
	for idx, textObj := range t.currentTextObjects {
		if idx < indexStart {
			continue
		} else if idx >= indexEnd {
			break
		}
		textObj = text.New(textObj.Orig, text.NewAtlas(face, text.ASCII))
		// The newly created *text.Text doesn't contain any glyphs to draw yet
		currLetter := string(t.currentTextString[idx])
		textObj.WriteString(currLetter)
	}
}

func (t *Texter) setTextRangeColor(col color.Color, indexStart, indexEnd int) {
	for idx, textObj := range t.currentTextObjects {
		if idx < indexStart {
			continue
		} else if idx >= indexEnd {
			break
		}
		textObj.Color = col
	}
}

type markdownCommand struct {
	idxStart          int
	idxEnd            int
	attributeValueMap map[string]string
}

func getMarkdownCommandSliceFromString(str string) ([]markdownCommand, string) {

	var markdownCommandSlice []markdownCommand
	var strippedStr string

	htmlOpenRegexp := regexp.MustCompile(`\<[^\/]\w+ style="(?:(\w+)\s*:\s*([^"; ]+))+"\>`)
	htmlCloseRegexp := regexp.MustCompile(`\<\/\w+\s?\>`)
	htmlRegexp := regexp.MustCompile(`\<.+?\>`)

	styleMatchSlice := htmlOpenRegexp.FindAllStringSubmatch(str, -1)
	styleIndexStartSlice := htmlOpenRegexp.FindAllStringIndex(str, -1)
	styleIndexEndSlice := htmlCloseRegexp.FindAllStringIndex(str, -1)

	// I don't do error checking here (e.g. same number of opening/closing brackets) because it is expected to be seen
	// in the markdown preview in an editor.

	strippedStr = htmlRegexp.ReplaceAllString(str, "")

	// for adjusting indices after the replacement
	cumulativeOffset := 0

	var currMarkdownCommand markdownCommand
	for i := 0; i < len(styleIndexStartSlice); i++ {
		currMarkdownCommand.attributeValueMap = make(map[string]string)
		currMarkdownCommand.attributeValueMap[styleMatchSlice[i][1]] = styleMatchSlice[i][2]
		currMarkdownCommand.idxStart = styleIndexStartSlice[i][0] - cumulativeOffset
		cumulativeOffset += styleIndexStartSlice[i][1] - styleIndexStartSlice[i][0]
		currMarkdownCommand.idxEnd = styleIndexEndSlice[i][0] - cumulativeOffset
		cumulativeOffset += styleIndexEndSlice[i][1] - styleIndexEndSlice[i][0]
		markdownCommandSlice = append(markdownCommandSlice, currMarkdownCommand)
	}

	return markdownCommandSlice, strippedStr
}

// TODO: This already needs to know for which ranges what formatting should be used.
// So this string still has it all, all the md!

// case read until '<'
// use default atlas and color and write to text object (could be too long for the box though...)
// This needs to be checked right here! Otherwise the following letters could be disturbed
// So right here there needs to be checked, is it too long or not!

// case '<' html command, read it in until '>'

// TODO: Add nesting (up until then no nested html)

// execute it e.g. get the atlas or change the color
// Read in the rest of the text until '<'
// Check if this is correct closing token '</span>'
// Write letters into text object ...
func (t *Texter) convertStringToTextObjectsInBox(str string, scn *Scene) {

	markdownCommandSlice, str := getMarkdownCommandSliceFromString(str)

	t.currentTextObjects = nil
	leftIndent := t.textBox.topLeftCorner.X + t.textBox.margin
	nextWordRegexp := regexp.MustCompile(`^[^\s]+ `)
	var nextWord string

	// starting point for writing characters
	// TODO: Add debug mode, where you can see the coordinates of the mouse... pixel coordinates...
	currentOrig := pixel.V(leftIndent, t.textBox.topLeftCorner.Y-2*t.textBox.margin)
	currentAtlas := scn.atlas
	// TODO: Add scn.textColor to the scene struct.
	// This has to be in sync with the background, also defined by the scene
	currentColor := colornames.Black

	for idx, rune := range str {

		if len(markdownCommandSlice) == 0 {
		} else if idx == markdownCommandSlice[0].idxStart {
			for attribute, value := range markdownCommandSlice[0].attributeValueMap {
				switch attribute {
				case `color`:
					currentColor = colornames.Map[strings.ToLower(value)]
				case `font-size`:
					// TODO: implement font size change via atlas
				}
			}
		} else if idx == markdownCommandSlice[0].idxEnd {
			markdownCommandSlice = markdownCommandSlice[1:]
			currentAtlas = scn.atlas
			currentColor = colornames.Black
		}

		char := string(rune)
		switch char {
		case `\n`:
			// align at left indent and remove one line height to the current Y Position
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -t.currentTextObjects[idx].LineHeight))
		case ` `:
			nextWord = nextWordRegexp.FindString(str[(idx + 1):])
		}

		newTextObject := text.New(currentOrig, currentAtlas)
		t.currentTextObjects = append(t.currentTextObjects, newTextObject)
		newTextObject.Color = currentColor

		newTextObject.WriteString(char)
		currentOrig = newTextObject.Dot

		if newTextObject.BoundsOf(` `+nextWord).Max.X >
			(t.textBox.topLeftCorner.X + t.textBox.dimensions.X - 2*t.textBox.margin) {
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -t.currentTextObjects[idx].LineHeight))
		}
	}

	return
}

// setText accepts a string with potential markdown formatting containing HTML with inline CSS for text formatting.
// e.g. roses are <span style="color:red">red</span>
// Online LF "\n" is used to mark a new line.
func (t *Texter) setText(str string, scn *Scene) {

	t.convertStringToTextObjectsInBox(str, scn)
	// TODO: repopulate
	// wrappedString := t.getWrappedString(str)

	// t.currentTextString = wrappedString
	// t.currentTextObject.Clear()
	// t.currentTextObject.WriteString(wrappedString)
}

func (t *Texter) addText(str string) {
	// TODO: repopulate
	// wrappedString := t.getWrappedString(t.currentTextString + str)

	// t.currentTextString = wrappedString
	// t.currentTextObject.Clear()
	// t.currentTextObject.WriteString(wrappedString)
}

func (t *Texter) drawTextInBox(win *pixelgl.Window) {
	t.textBox.drawTextBox(win)

	// TODO: The y coordinate is guesswork and probably dependend on the font face used!
	// TODO: This can be either done in convertStringToTextObjectsInBox or here. One should go.
	// marginVec := pixel.V(t.textBox.margin, t.textBox.margin-t.textBox.thickness-2.4*t.currentTextObject.LineHeight)
	// t.currentTextObject.Draw(win, pixel.IM.Moved(t.textBox.topLeftCorner.Add(marginVec)))
	for _, textObj := range t.currentTextObjects {
		textObj.Draw(win, pixel.IM)
	}
}
