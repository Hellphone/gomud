package helpers

import (
	"bufio"
	"net"
	"strings"
)

func GetInput(conn net.Conn) (string, error) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	message = strings.TrimRight(message, "\r\n")

	return message, nil
}
