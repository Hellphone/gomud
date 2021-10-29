package server

import (
	"context"
	"errors"
	"net"

	"github.com/hellphone/gomud/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Context  context.Context
	DBClient *mongo.Client
	Clients  *ClientList
	commands map[string]HandlerFunc
}

type HandlerFunc func(c Connection) error

type ClientList struct {
	Clients []Client
}

// TODO: add states
type Client struct {
	Connection net.Conn
	User       *models.User
}

type Connection net.Conn

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

	return nil, models.ErrorCommandNotFound
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

func (c *ClientList) CloseConnection(conn net.Conn) error {
	for k, client := range c.Clients {
		if client.Connection == conn {
			// TODO: add mutex
			c.Clients = RemoveIndex(c.Clients, k)
			// TODO: find out why all connections are closing here
			return conn.Close()
		}
	}
	return nil
}

func RemoveIndex(s []Client, index int) []Client {
	result := make([]Client, 0)
	result = append(result, s[:index]...)
	return append(result, s[index+1:]...)
}
