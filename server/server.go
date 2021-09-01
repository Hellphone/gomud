package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Connection net.Conn
	Context    context.Context
	DBClient   *mongo.Client
	commands   map[string]HandlerFunc
}

type HandlerFunc func() error

// TODO: probably store commands in the database (is this really needed?..)
func (s *Server) RegisterCommand(name string, handler HandlerFunc) error {
	if s.commands == nil {
		s.commands = map[string]HandlerFunc{}
	}
	if _, ok := s.commands[name]; ok {
		return errors.New("command with this name has already been registered")
	}
	s.commands[name] = handler

	return nil
}

func (s *Server) GetCommand(name string) (HandlerFunc, error) {
	if command, ok := s.commands[name]; ok {
		return command, nil
	}

	return nil, errors.New(fmt.Sprintf("%s: command with this name has not been registered", name))
}

func (s *Server) RegisterCommands() error {
	err := s.RegisterCommand("login", s.LoginHandler)
	if err != nil {
		return err
	}
	err = s.RegisterCommand("register", s.RegisterHandler)
	if err != nil {
		return err
	}
	err = s.RegisterCommand("db", s.DBHandler)
	if err != nil {
		return err
	}
	err = s.RegisterCommand("exit", s.ExitHandler)
	if err != nil {
		return err
	}

	return nil
}
