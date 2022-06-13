package main

import (
	"bufio"
	"fmt"
	. "github.com/ordovician/mainframe"
	"io"
	"log"
	"os"
	"time"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

type Terminal struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (term *Terminal) Login() {
	scanner := bufio.NewScanner(term.Stdin)
	fmt.Fprint(term.Stdout, "Login: ")
	if !scanner.Scan() {
		fmt.Fprintln(term.Stderr, "Unable to get login name:", scanner.Err())
		os.Exit(1)
	}

	user := scanner.Text()

	fmt.Fprint(term.Stdout, "Password: ")
	if !scanner.Scan() {
		fmt.Fprintln(term.Stderr, "Unable to get password:", scanner.Err())
		os.Exit(1)
	}

	passwd := scanner.Text()

	if ok, err := CheckLogin(user, passwd); ok {
		curtime := time.Now().Local().Format("January 2, 2006")

		// Name from Soviet BESM-6 computer https://en.wikipedia.org/wiki/BESM
		fmt.Fprintf(term.Stdout, "Большая Электронно-Счётная Машина 6: %v\n", curtime)
		fmt.Fprintf(term.Stdout, "Welcome to БЭСМ-6 comrade %s\n", user)

		runShell()
	} else if err != nil {
		fmt.Fprintln(term.Stderr, "Could not log in because:", err)
	} else {
		fmt.Fprint(term.Stdout, "Username does not exist or password was wrong")
	}
}

func main() {
	term := &Terminal{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	term.Login()
}
