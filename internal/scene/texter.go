// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"regexp"
	"strings"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

// Texter is defined by its text
type Texter struct {
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

// getWrappedString executed on a Texter returns a string wrapped inside the Texter's box.
// This currently depends on the first text object in Texter's currentTextObjects.
func (t *Texter) getWrappedString(str string) string {
	maxTextWidth := t.textBox.dimensions.X - t.textBox.margin - 20
	textObject := t.currentTextObjects[0]

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

	return wrappedString
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

func (t *Texter) applyMarkdownCommand(markdownCommandSlice []markdownCommand, idx int, scn *Scene) {

	if len(markdownCommandSlice) == 0 {
	} else if idx == markdownCommandSlice[0].idxStart {
		for attribute, value := range markdownCommandSlice[0].attributeValueMap {
			switch attribute {
			case `color`:
				t.color = colornames.Map[strings.ToLower(value)]
			case `font-size`:
				// TODO: implement font size change via atlas
			}
		}
	} else if idx == markdownCommandSlice[0].idxEnd {
		markdownCommandSlice = markdownCommandSlice[1:]
		t.atlas = scn.atlas
		t.color = scn.textColor
	}
}

func (t *Texter) convertMarkdownStringToTextObjectsInBox(str string, scn *Scene) {

	markdownCommandSlice, str := getMarkdownCommandSliceFromString(str)
	t.currentTextString = str

	t.currentTextObjects = nil
	leftIndent := t.textBox.topLeftCorner.X + t.textBox.margin
	nextWordRegexp := regexp.MustCompile(`^[^\s]+ `)
	var nextWord string

	// starting point for writing characters
	currentOrig := pixel.V(leftIndent, t.textBox.topLeftCorner.Y-2*t.textBox.margin)
	t.atlas = scn.atlas
	t.color = scn.textColor

	for idx, rune := range str {

		t.applyMarkdownCommand(markdownCommandSlice, idx, scn)

		char := string(rune)
		switch char {
		case `\n`:
			// align at left indent and remove one line height to the current Y Position
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -t.currentTextObjects[idx].LineHeight))
		case ` `:
			nextWord = nextWordRegexp.FindString(str[(idx + 1):])
		}

		newTextObject := text.New(currentOrig, t.atlas)
		t.currentTextObjects = append(t.currentTextObjects, newTextObject)
		newTextObject.Color = t.color

		newTextObject.WriteString(char)
		currentOrig = newTextObject.Dot

		if newTextObject.BoundsOf(nextWord).Max.X >
			(t.textBox.topLeftCorner.X + t.textBox.dimensions.X - 2*t.textBox.margin) {
			currentOrig = currentOrig.Add(pixel.V(leftIndent-currentOrig.X, -t.currentTextObjects[idx].LineHeight))
		}
		// to not check again until next space indicating the next word
		nextWord = ""
	}
}

// setText accepts a string with potential markdown formatting containing HTML with inline CSS for text formatting.
// e.g. roses are <span style="color:red">red</span>
// Online LF "\n" is used to mark a new line.
func (t *Texter) setTextLetterByLetter(str string, scn *Scene) {

	t.convertMarkdownStringToTextObjectsInBox(str, scn)
}

func (t *Texter) setTextFontFace(face font.Face) {
	textObject := t.currentTextObjects[0]
	textObject = text.New(textObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	textObject.WriteString(t.currentTextString)
}

func (t *Texter) setTextColor(col color.Color) {
	t.currentTextObjects[0].Color = col
}

func (t *Texter) setText(str string) {
	wrappedString := t.getWrappedString(str)

	t.currentTextString = wrappedString
	t.currentTextObjects[0].Clear()
	t.currentTextObjects[0].WriteString(wrappedString)
}

func (t *Texter) addText(str string, scn *Scene) {

	wrappedString := t.getWrappedString(t.currentTextString + str)

	t.currentTextString = wrappedString
	t.currentTextObjects[0].Clear()
	t.currentTextObjects[0].WriteString(wrappedString)
}

func (t *Texter) drawTextLetterByLetter(win *pixelgl.Window) {
	t.textBox.drawTextBox(win)

	for _, textObj := range t.currentTextObjects {
		textObj.Draw(win, pixel.IM)
	}
}

func (t *Texter) drawTextInBox(win *pixelgl.Window) {
	t.textBox.drawTextBox(win)

	// TODO: The y coordinate is guesswork and probably dependend on the font face used!
	marginVec := pixel.V(t.textBox.margin, t.textBox.margin-t.textBox.thickness-2.4*t.currentTextObjects[0].LineHeight)
	t.currentTextObjects[0].Draw(win, pixel.IM.Moved(t.textBox.topLeftCorner.Add(marginVec)))
}
