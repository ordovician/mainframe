package mainframe

import (
	"bufio"
	"fmt"
	"io"
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
func (term *Terminal) Login() error {
	scanner := bufio.NewScanner(term.Stdin)
	fmt.Fprint(term.Stdout, "Login: ")
	if !scanner.Scan() {
		return fmt.Errorf("unable to get login name: %w", scanner.Err())
	}

	user := scanner.Text()

	fmt.Fprint(term.Stdout, "Password: ")
	if !scanner.Scan() {
		return fmt.Errorf("unable to get password: %w", scanner.Err())
	}

	passwd := scanner.Text()

	if ok, err := CheckLogin(user, passwd); ok {
		curtime := time.Now().Local().Format("January 2, 2006")

		// Name from Soviet BESM-6 computer https://en.wikipedia.org/wiki/BESM
		fmt.Fprintf(term.Stdout, "Большая Электронно-Счётная Машина 6: %v\n", curtime)
		fmt.Fprintf(term.Stdout, "Welcome to БЭСМ-6 comrade %s\n", user)

		term.runShell()
	} else if err != nil {
		return fmt.Errorf("could not log in because: %w", err)
	} else {
		return fmt.Errorf("username does not exist or password was wrong")
	}
	return nil
}

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
