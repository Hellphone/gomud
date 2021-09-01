package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/hellphone/gomud/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

var names = map[string]string{
	"john":    "12345",
	"abraham": "lincoln",
	"rose":    "MaRy",
	"chloe":   "Zjk%d82*ja)",
}

func (s *Server) LoginHandler() error {
	// TODO: keep constantly reading the input until correct or the user interrupts the process
	fmt.Fprintf(s.Connection, "Enter your name:\r\n")
	name, err := helpers.GetInput(s.Connection)
	if err != nil {
		return err
	}
	password, ok := names[name]
	if ok {
		fmt.Fprintf(s.Connection, "Enter your password:\r\n")
		// TODO: hide password input
		p, err := helpers.GetInput(s.Connection)
		if err != nil {
			return err
		}
		if p == password {
			ClearScreen()
			fmt.Fprintf(s.Connection, "You successfully logged in as %v!\r\n", name)
			// TODO: get user from DB
			//s.User :=
			s.StartGame()
		} else {
			fmt.Fprintf(s.Connection, "Sorry, but the password is not correct.\r\n")
		}
	} else {
		fmt.Fprintf(s.Connection, "Sorry, but there is no such name.\r\n")
	}

	return nil
}

func (s *Server) RegisterHandler() error {
	fmt.Fprintf(s.Connection, "Enter your name:\r\n")
	name, _ := helpers.GetInput(s.Connection)
	_, ok := names[name]
	if ok {
		fmt.Fprintf(s.Connection, "Sorry, but this name has already been taken.\r\n")
		// TODO: run this case again
	} else {
		fmt.Fprintf(s.Connection, "Enter your password:\r\n")
		pass, _ := helpers.GetInput(s.Connection)
		fmt.Fprintf(s.Connection, "Confirm your password:\r\n")
		pass2, _ := helpers.GetInput(s.Connection)
		if pass == pass2 {
			names[name] = pass
			// TODO: log the user in
			fmt.Fprintf(s.Connection, "You have been successfully registered as %v!\r\n", name)
		}
	}

	return nil
}

func (s *Server) DBHandler() error {
	databases, err := s.DBClient.ListDatabaseNames(s.Context, bson.M{})
	if err != nil {
		return err
	}
	fmt.Fprintf(s.Connection, "databases list: %+v", databases)

	return nil
}

func (s *Server) ExitHandler() error {
	fmt.Fprintf(s.Connection, "Goodbye!")
	s.Connection.Close()

	return nil
}

func (s *Server) DefaultCommand(conn net.Conn, command string) error {
	fmt.Fprintf(conn, "What would you like to do?\r\n"+
		"[login, register, exit]\r\n")

	return nil
}

func (s *Server) StartGame() {
	// TODO: change user state to in-game
	s.User.SwitchStatus()
	fmt.Fprintf(s.Connection, "Your adventure starts here...\r\n")
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
