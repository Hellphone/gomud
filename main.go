package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	CONN_PORT = ":23"
	CONN_TYPE = "tcp"
)

func main() {
	fmt.Println("Start server...")

	ln, _ := net.Listen(CONN_TYPE, CONN_PORT)
	conn, _ := ln.Accept()

	fmt.Fprintf(conn, "Hello stranger! Welcome to GOMUD!\r\n"+
		"What would you like to do?\r\n"+
		"[login, register, exit]\r\n")

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Fprintf(conn, message)
	}
}
