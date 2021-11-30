package helpers

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

func GetCommandInput(conn net.Conn) (string, string, error) {
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", "", err
	}
	input, err = processInputString(input)
	if err != nil {
		return "", "", err
	}

	words := strings.Fields(input)
	if len(words) < 2 {
		return words[0], "", nil
	}

	return words[0], words[1], nil
}

func GetInput(conn net.Conn) (string, error) {
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	input, err = processInputString(input)
	if err != nil {
		return "", err
	}

	return input, nil
}

// TODO: search for golang sanitizing library
func processInputString(str string) (string, error) {
	// TODO: sanitize input (make the regexp right)
	// TODO: get rid of these symbols: ←[A, ←[B, ←[C, ←[D (arrows, backspaces etc.)
	reg, err := regexp.Compile(`[^a-zA-Z0-9 ]+`)
	if err != nil {
		return str, err
	}

	processedString := reg.ReplaceAllString(str, "")

	return processedString, nil
}

func PrintToServer(str string) error {
	fmt.Printf("%s\r\n", str)
	return nil
}

func PrintToClient(conn net.Conn, str string, a ...interface{}) error {
	str = fmt.Sprintf("%s\r\n", str)
	fmt.Fprintf(conn, str, a)
	return nil
}
