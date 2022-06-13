package main

import (
	. "github.com/ordovician/mainframe"
	"log"
	"os"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

func main() {
	term := &Terminal{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	term.Login()
}
