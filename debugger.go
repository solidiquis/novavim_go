package main

import (
	"fmt"
	"log"
	"os"
)

// Write logs to a separate terminal
// Ex:
// wr := debug("ttys012")
// wr("I am writing into another terminal")
func debug(charDevice string) func(string) {
	dev := fmt.Sprintf("/dev/%s", charDevice)
	f, err := os.OpenFile(dev, os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	f.Write([]byte("\n"))

	return func(txt string) {
		f.Write([]byte(txt))
	}
}
