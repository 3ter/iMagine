package utils

import (
	"io/ioutil"
	"log"
)

// LoadFileToString loads the contents of a file into a string or dies
func LoadFileToString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
