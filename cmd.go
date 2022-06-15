package mainframe

import (
	_ "embed"
	"fmt"
	"path"
	"strings"
)

//go:embed key.base64.txt
var cryptKey string

// Code returned when running commands such as ListCmd and CatCmd
type ExitCode int

const (
	Ok      ExitCode = iota // normal result
	Failure                 // Command could not complete task
	Quit                    // request to quit shell
)

type Command interface {
	Name() string
	Help(term *Terminal)
	Run(term *Terminal, args []string) ExitCode
}

type HelpCmd struct{}
type ListCmd struct{}
type CatCmd struct{}
type DecryptCmd struct{}
type ExitCmd struct{}

func (cmd *HelpCmd) Name() string {
	return "help"
}

func (cmd *HelpCmd) Help(term *Terminal) {
	fmt.Fprintln(term.Stdout, `NAME
    help -- give help for a command
SYNOPSIS
    help command
DESCRIPTION
    The help command displays help for a command.`)
}

func (cmd *HelpCmd) Run(term *Terminal, args []string) ExitCode {
	if len(args) == 0 {
		fmt.Fprintln(term.Stdout, `Valid commands:
    ls   cat   decrypt	exit   help`)
		return Ok
	}

	helpCmd := Lookup(args[0])
	if helpCmd == nil {
		fmt.Fprintln(term.Stdout, "Unknown command:", args[0])
		return Failure
	} else {
		helpCmd.Help(term)
	}
	return Ok
}

func (cmd *ListCmd) Name() string {
	return "ls"
}

func (cmd *ListCmd) Help(term *Terminal) {
	fmt.Fprintln(term.Stdout, `NAME
	ls -- list content of current directory
SYNOPSIS
    ls [directory]
DESCRIPTION
    The ls command shows all files and directories in working directory`)
}

func (cmd *ListCmd) Run(term *Terminal, args []string) ExitCode {
	entries, _ := storage.ReadDir("data")
	for _, entry := range entries {
		fmt.Fprintln(term.Stdout, entry.Name())
	}
	fmt.Fprintln(term.Stdout)
	return Ok
}

func (cmd *CatCmd) Name() string {
	return "cat"
}

func (cmd *CatCmd) Help(term *Terminal) {
	fmt.Fprintln(term.Stdout, `NAME
    cat -- concatenate and print files
SYNOPSIS
    cat [file ...]
DESCRIPTION
    The cat utility reads files sequentually and write them to standard output.`)
}

func (cmd *CatCmd) Run(term *Terminal, args []string) ExitCode {
	if len(args) < 1 {
		fmt.Fprintln(term.Stdout, "Missing file argument")
	} else {
		for _, arg := range args {
			filepath := path.Join("data", arg)
			content, err := storage.ReadFile(filepath)
			if err != nil {
				fmt.Fprintln(term.Stderr, "Could not read file:", err)
				return Failure
			}
			fmt.Fprintln(term.Stdout, string(content))
		}
	}
	return Ok
}

func (cmd *DecryptCmd) Name() string {
	return "decrypt"
}

func (cmd *DecryptCmd) Help(term *Terminal) {
	fmt.Fprintln(term.Stdout, `NAME
    decrypt -- decrypts and file and print contents
SYNOPSIS
	decrypt file
DESCRIPTION
    The decrypt utility decrypts a single file and write it to standard output.`)
}

func (cmd *DecryptCmd) Run(term *Terminal, args []string) ExitCode {
	if len(args) < 1 {
		fmt.Fprintln(term.Stdout, "Missing file argument")
		return Failure
	} else {
		keystorage := strings.NewReader(cryptKey)
		key, err := LoadEncodedKey(keystorage, Base64)
		if err != nil {
			fmt.Fprintf(term.Stderr, "Key required to decrypt data is missing: %v\n", err)
			return Failure
		}

		cip, err := NewCipher(key)
		if err != nil {
			return Failure
		}
		for _, arg := range args {
			msg, err := cip.DecryptFile(arg)

			if err != nil {
				fmt.Fprintf(term.Stderr, "Could not read file %s: %s\n", arg, err)
				return Failure
			}
			fmt.Fprintln(term.Stdout, string(msg))
		}
	}
	return Ok
}

func (cmd *ExitCmd) Name() string {
	return "exit"
}

func (cmd *ExitCmd) Help(term *Terminal) {
	fmt.Fprintln(term.Stdout, `NAME
    exit -- exit the shell
SYNOPSIS
    exit
DESCRIPTION
    exit the mainframe shell. Logs the user out.`)
}

func (cmd *ExitCmd) Run(term *Terminal, args []string) ExitCode {
	return Quit
}

var commands = [...]Command{
	new(HelpCmd),
	new(ListCmd),
	new(DecryptCmd),
	new(CatCmd),
	new(ExitCmd),
}

// Lookup command with given name. Returns nil if command
// of that name is not found
func Lookup(cmdName string) Command {
	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			return cmd
		}
	}
	return nil
}
