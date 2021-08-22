package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/hellphone/gomud/server"
)

const (
	CONN_PORT = ":23"
	CONN_TYPE = "tcp"
)

func main() {
	fmt.Println("Starting server...")
	ln, _ := net.Listen(CONN_TYPE, CONN_PORT)
	conn, _ := ln.Accept()

	s := &server.Server{
		Connection: conn,
	}

	s.RegisterCommands()

	fmt.Fprintf(s.Connection, "Hello stranger! Welcome to GOMUD!\r\n"+
		"What would you like to do?\r\n"+
		"[login, register, exit]\r\n")

	for {
		message, _ := bufio.NewReader(s.Connection).ReadString('\n')
		message = strings.TrimRight(message, "\r\n")
		command, err := s.GetCommand(message)
		// TODO: learn how to handle errors properly
		if err != nil {
			fmt.Println("Error:", err)
		}

		err = command(conn)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
