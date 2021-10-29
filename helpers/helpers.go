package helpers

import (
	"bufio"
	"net"
	"regexp"
)

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
	reg, err := regexp.Compile(`[^a-zA-Z0-9]+`)
	if err != nil {
		return str, err
	}

	processedString := reg.ReplaceAllString(str, "")

	return processedString, nil
}
