// Package controltext implements functions that control where the text goes on
// the screen and how it is presented (e.g. coming in all at once or one letter at a time)
package controltext

import (
	"time"

	"github.com/faiface/pixel/text"
)

// WriteToTextLetterByLetter accepts a text object that is filled with the provided message
// letter by letter so the eyes can follow as the text is being written.
//
// This should be called as a separate goroutine as it would be blocking otherwise.
//
// interval				time in milliseconds between writing two letters
// writingDoneChannel	channel which gets sent 1 when this function returns
func WriteToTextLetterByLetter(txt *text.Text, msg string, interval time.Duration, writingDoneChannel chan int) {
	<-writingDoneChannel
	for _, char := range msg {
		_, err := txt.WriteString(string(char))
		if err != nil {
			panic(err)
		}
		time.Sleep(interval * time.Millisecond)
	}
	writingDoneChannel <- 1
	return
}
