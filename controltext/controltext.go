// Package controltext implements functions that control where the text goes on
// the screen and how it is presented (e.g. coming in all at once or one letter at a time)
package controltext

import (
	"sync"
	"time"

	"github.com/faiface/pixel/text"
)

// SafeText adds a mutex to a text to control access to the text object
//
// This adds a little complexity to the creation of those objects.
//
// Example 1
// txt := text.Text.New(vec, atlas)
// safeTxt := SafeText{
// 	*text.Text: txt
// }
type SafeText struct {
	*text.Text
	sync.Mutex
}

// WriteToTextLetterByLetter accepts a text object that is filled with the provided message
// letter by letter so the eyes can follow as the text is being written.
//
// This should be called as a separate goroutine as it would be blocking otherwise.
//
// interval				time in milliseconds between writing two letters
// writingDoneChannel	channel which gets sent 1 when this function returns
//
// The function is waiting to receive a 'go' on the 'writingDoneChannel' channel.
// Accordingly at the end it doesn't return until the next writing function receives
// the baton to start writing.
//
// Example 1
//  msg := "Some text\n"
//  writingDoneChannel := make(chan int)
//  go controltext.WriteToTextLetterByLetter(s.title, msg, 60, writingDoneChannel)
//  writingDoneChannel <- 1 // init writing the first line
//  msg = "Press Ctrl + Q to quit or Escape for main menu.\n"
//  go controltext.WriteToTextLetterByLetter(s.title, msg, 10, writingDoneChannel)
func WriteToTextLetterByLetter(txt *SafeText, msg string, interval time.Duration, writingDoneChannel chan int) {
	<-writingDoneChannel
	txt.Lock()
	for _, char := range msg {
		_, err := txt.WriteString(string(char))
		if err != nil {
			panic(err)
		}
		time.Sleep(interval * time.Millisecond)
	}
	txt.Unlock()
	writingDoneChannel <- 1
	return
}
