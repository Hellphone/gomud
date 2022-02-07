package server

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hellphone/gomud/domain/models"
	"github.com/hellphone/gomud/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

var names = map[string]string{
	"john":    "12345",
	"abraham": "lincoln",
	"rose":    "MaRy",
	"chloe":   "Zjk%d82*ja)",
}

// TODO: maybe pass a struct storing a connection
func (s *Server) LoginHandler(client *Client, args ...string) error {
	// TODO: keep constantly reading the input until correct or the user interrupts the process
	fmt.Fprintf(client.Connection, "Enter your name:\r\n")
	name, err := helpers.GetInput(client.Connection)
	if err != nil {
		return err
	}
	password, ok := names[name]
	if ok {
		fmt.Fprintf(client.Connection, "Enter your password:\r\n")
		// TODO: hide password input
		p, err := helpers.GetInput(client.Connection)
		if err != nil {
			return err
		}
		if p == password {
			ClearScreen()
			fmt.Fprintf(client.Connection, "You successfully logged in as %v!\r\n", name)
			// TODO: create a function linking a user to a connection
			for k, v := range s.Clients.Clients {
				if v.Connection == client.Connection {
					s.Clients.Clients[k].User.Login = name
					s.Clients.Clients[k].User.Status = models.StatusOnline
					break
				}
			}
			// TODO: get user from DB
			//s.User :=
			s.StartGame(client)
		} else {
			fmt.Fprintf(client.Connection, "Sorry, but the password is not correct.\r\n")
		}
	} else {
		fmt.Fprintf(client.Connection, "Sorry, but there is no such name.\r\n")
	}

	return nil
}

func (s *Server) RegisterHandler(client *Client, args ...string) error {
	fmt.Fprintf(client.Connection, "Enter your name (only alphabetical and numeric symbols are allowed):\r\n")
	name, _ := helpers.GetInput(client.Connection)
	_, ok := names[name]
	if ok {
		fmt.Fprintf(client.Connection, "Sorry, but this name has already been taken.\r\n")
		// TODO: run this case again
	} else {
		fmt.Fprintf(client.Connection, "Enter your password:\r\n")
		pass, _ := helpers.GetInput(client.Connection)
		fmt.Fprintf(client.Connection, "Confirm your password:\r\n")
		pass2, _ := helpers.GetInput(client.Connection)
		if pass == pass2 {
			names[name] = pass
			// TODO: log the user in
			// TODO: highlight the name in different color if possible
			fmt.Fprintf(client.Connection, "You have been successfully registered as %v!\r\n", name)
		}
	}

	return nil
}

func (s *Server) DBHandler(client *Client, args ...string) error {
	databases, err := s.DBClient.ListDatabaseNames(s.Context, bson.M{})
	if err != nil {
		return err
	}
	fmt.Fprintf(client.Connection, "databases list: %+v", databases)

	return nil
}

func (s *Server) KickoutHandler(client *Client, args ...string) error {
	// TODO: add confirmation ("Are you sure to kick %USERNAME% out?)"

	if client.User.Status != models.StatusOnline {
		fmt.Fprintf(client.Connection, "Not enough rights to proceed with this operation\r\n")
		return nil
	}

	if len(args) < 1 || args[0] == "" {
		return models.ErrorNotEnoughArguments
	}

	username := args[0]
	kicked := false
	for _, c := range s.Clients.Clients {
		if c.User.Login == username {
			c.User.SwitchStatus()
			kicked = true
			fmt.Fprintf(client.Connection, "You have successfully kicked %s out\r\n", username)
			fmt.Fprintf(c.Connection, "You have been kicked out by %s\r\n", client.User.Login)
			err := s.Clients.CloseConnection(c.Connection)
			if err != nil {
				return err
			}

			break
		}
	}

	if !kicked {
		return models.ErrorUserNotFound
	}

	return nil
}

func (s *Server) ExitHandler(client *Client, args ...string) error {
	fmt.Fprintf(client.Connection, "Goodbye!\r\n")
	// TODO: figure out how to run this method (not by passing Clients to Server)
	// (it's okay)
	err := s.Clients.CloseConnection(client.Connection)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DefaultCommand(client *Client, command string) error {
	fmt.Fprintf(client.Connection, "What would you like to do?\r\n"+
		"[login, register, exit]\r\n")

	return nil
}

func (s *Server) StartGame(client *Client) {
	// TODO: change user state to in-game (online?)
	//s.User.SwitchStatus()
	fmt.Fprintf(client.Connection, "Your adventure starts here...\r\n")
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
