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
//
// The function is waiting to receive a 'go' on the 'writingDoneChannel' channel.
// Accordingly at the end it doesn't return until the next writing function receives
// the baton to start writing.
//
// Example 1
//  titleString := "Welcome to the START. Here is nothing... (yet)!\n"
//  writingDoneChannel := make(chan int)
//  go controltext.WriteToTextLetterByLetter(s.title, titleString, 60, writingDoneChannel)
//  writingDoneChannel <- 1 // init writing the first line
//  titleString = "Press Ctrl + Q to quit or Escape for main menu.\n"
//  go controltext.WriteToTextLetterByLetter(s.title, titleString, 10, writingDoneChannel)
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
