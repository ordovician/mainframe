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
	// you can connect using NetCat or telnet, see README.md
	server, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer server.Close()

	// accept multiple connections
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleClientConnection(conn)
	}
}

// handles a single connection from a client
func handleClientConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	term := &Terminal{
		Stdin:  reader,
		Stdout: conn,
		Stderr: conn,
	}

	term.Login()

	// close socket connection
	conn.Close()
}
