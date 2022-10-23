package mainframe

// NOTE: If you are missing the stringer command, you can install it
// $ go install golang.org/x/tools/cmd/stringer@latest

//go:generate stringer -type=ExitCode exitcode.go

// Code returned when running commands such as ListCmd and CatCmd
type ExitCode int

const (
	Ok      ExitCode = iota // normal result
	Failure                 // Command could not complete task
	Quit                    // request to quit shell
)
