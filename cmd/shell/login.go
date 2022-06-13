package main

import (
	"bufio"
	"fmt"
	. "github.com/ordovician/mainframe"
	"log"
	"os"
	"time"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Login: ")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Unable to get login name:", scanner.Err())
		os.Exit(1)
	}

	user := scanner.Text()

	fmt.Print("Password: ")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Unable to get password:", scanner.Err())
		os.Exit(1)
	}

	passwd := scanner.Text()

	if ok, err := CheckLogin(user, passwd); ok {
		// Name from Soviet BESM-6 computer https://en.wikipedia.org/wiki/BESM
		fmt.Printf("Большая Электронно-Счётная Машина 6: %v\n", time.Now().Local().Format("January 2, 2006"))
		fmt.Printf("Welcome to БЭСМ-6 comrade %s\n", user)

		runShell()
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Could not log in becase:", err)
	} else {
		fmt.Print("Username does not exist or password was wrong")
	}
}
