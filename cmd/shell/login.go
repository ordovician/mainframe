package main

import (
	"log"
	"os"

	. "github.com/ordovician/mainframe"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

func main() {
	term := &Terminal{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := term.Login(); err != nil {
		log.Fatal(err)
	}
}
