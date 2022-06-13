package main

import (
	"bufio"
	"log"
	"net"
	"os"

	. "github.com/ordovician/mainframe"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

func main() {
	// listen on port 1234
	server, _ := net.Listen("tcp", ":1234")

	// accept connection
	conn, _ := server.Accept()
	reader := bufio.NewReader(conn)

	term := &Terminal{
		Stdin:  reader,
		Stdout: conn,
		Stderr: conn,
	}

	term.Login()
	conn.Close()
}
