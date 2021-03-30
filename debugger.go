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
func Debug(charDevice string) func([]byte) {
	dev := fmt.Sprintf("/dev/%s", charDevice)
	f, err := os.OpenFile(dev, os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	f.Write([]byte("\n"))

	return func(b []byte) {
		output := fmt.Sprintf("%v -> %s\n", b, b)
		f.Write([]byte(output))
	}
}
