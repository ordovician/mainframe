package mainframe

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Terminal struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Perform login assuming user is connected through the stdin and stdout defined
// by the terminal object.
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

		term.runShell()
	} else if err != nil {
		fmt.Fprintln(term.Stderr, "Could not log in because:", err)
	} else {
		fmt.Fprint(term.Stdout, "Username does not exist or password was wrong")
	}
}

//go:embed data
var storage embed.FS

// runShell creates a command line where you can issue
// simple commands such as ls and cat in a pretend Unix shell.
func (term *Terminal) runShell() {

	scanner := bufio.NewScanner(term.Stdin)
	fmt.Fprint(term.Stdout, "> ")

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		cmd := Lookup(fields[0])
		if cmd == nil {
			fmt.Fprintln(term.Stdout, "Unknown command:", fields[0])
		} else {
			if cmd.Run(term, fields[1:]) == Quit {
				break
			}
		}
		fmt.Fprint(term.Stdout, "> ")
	}
}
