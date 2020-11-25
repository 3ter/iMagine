// Package scene implements functions to provide the contents of a scene
// like the backgrounds and texts and music.
package scene

import (
	"image/color"
	"regexp"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// Texter is defined by its text
type Texter struct {
	fontFace font.Face
	// The color in this object is settable, for changing the font a new object
	// is needed.
	currentTextObject *text.Text
	currentTextString string

	textBox *TextBox
}

func (t *Texter) setTextFontFace(face font.Face) {
	t.currentTextObject = text.New(t.currentTextObject.Orig, text.NewAtlas(face, text.ASCII))
	// The newly created *text.Text doesn't contain any glyphs to draw yet
	t.currentTextObject.WriteString(t.currentTextString)
}

func (t *Texter) setTextColor(col color.Color) {
	t.currentTextObject.Color = col
}

func (t *Texter) getWrappedString(str string) string {
	maxTextWidth := t.textBox.dimensions.X - t.textBox.margin - 20

	m := regexp.MustCompile(`\n`) // GoogleDocs uses LF to mark its line endings
	currLinesSlice := m.Split(str, -1)
	var wrappedString string
	for idx, currLine := range currLinesSlice {
		if t.currentTextObject.BoundsOf(str).W() <= maxTextWidth {
			wrappedString = str
			break
		}
		if t.currentTextObject.BoundsOf(currLine).W() <= maxTextWidth {
			if idx < len(currLinesSlice)-1 {
				wrappedString += currLine + "\n"
			} else {
				wrappedString += currLine
			}
			continue
		}
		newLineBreakIndex := int((maxTextWidth) * (float64(len(currLine))) / t.currentTextObject.BoundsOf(currLine).W())
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

func (t *Texter) setText(str string) {
	wrappedString := t.getWrappedString(str)

	t.currentTextString = wrappedString
	t.currentTextObject.Clear()
	t.currentTextObject.WriteString(wrappedString)
}

func (t *Texter) addText(str string) {
	wrappedString := t.getWrappedString(t.currentTextString + str)

	t.currentTextString = wrappedString
	t.currentTextObject.Clear()
	t.currentTextObject.WriteString(wrappedString)
}

func (t *Texter) drawTextInBox(win *pixelgl.Window) {
	t.textBox.drawTextBox(win)

	// TODO: The y coordinate is guesswork and probably dependend on the font face used!
	marginVec := pixel.V(t.textBox.margin, t.textBox.margin-t.textBox.thickness-2.4*t.currentTextObject.LineHeight)
	t.currentTextObject.Draw(win, pixel.IM.Moved(t.textBox.topLeftCorner.Add(marginVec)))
}
