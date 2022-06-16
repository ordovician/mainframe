package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	. "github.com/ordovician/mainframe"
)

var debug *log.Logger = log.New(os.Stdout, "DEBUG: ", 0)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of server:\n")
		flag.PrintDefaults()
	}

	var (
		port     int
		protocol string
	)

	flag.IntVar(&port, "port", 1234, "Port number for server to listen for connections on")
	flag.StringVar(&protocol, "protocol", "tcp", "Protocol such as TCP or UDP for socket connection")
	flag.Parse()

	protocol = strings.ToLower(protocol)
	portStr := fmt.Sprintf(":%d", port)

	fmt.Printf("Starting server on port %d using protocol %s\n", port, protocol)
	fmt.Println("Hit Ctrl-C to stop")

	// listen on port 1234
	// you can connect using NetCat or telnet, see README.md
	server, err := net.Listen(protocol, portStr)
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
	// close socket connection when done or in case of panic
	defer conn.Close()

	reader := bufio.NewReader(conn)

	term := &Terminal{
		Stdin:  reader,
		Stdout: conn,
		Stderr: conn,
	}

	if err := term.Login(); err != nil {
		log.SetPrefix("ERROR: ")
		log.Print(err)
	}
}
