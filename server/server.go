package server

import (
	"errors"
	"net"
)

type Server struct {
	Connection net.Conn
	Commands   map[string]string
}

// TODO: register commands with functions, not strings
func (s *Server) RegisterCommand(name string, handler string) error {
	if _, ok := s.Commands[name]; ok {
		return errors.New("command with this name has already been registered")
	}
	s.Commands[name] = handler

	return nil
}
