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

func (s *Server) LoginHandler(conn Connection) error {
	// TODO: keep constantly reading the input until correct or the user interrupts the process
	fmt.Fprintf(conn, "Enter your name:\r\n")
	name, err := helpers.GetInput(conn)
	if err != nil {
		return err
	}
	password, ok := names[name]
	if ok {
		fmt.Fprintf(conn, "Enter your password:\r\n")
		// TODO: hide password input
		p, err := helpers.GetInput(conn)
		if err != nil {
			return err
		}
		if p == password {
			ClearScreen()
			fmt.Fprintf(conn, "You successfully logged in as %v!\r\n", name)
			// TODO: get user from DB
			//s.User :=
			s.StartGame(conn)
		} else {
			fmt.Fprintf(conn, "Sorry, but the password is not correct.\r\n")
		}
	} else {
		fmt.Fprintf(conn, "Sorry, but there is no such name.\r\n")
	}

	return nil
}

func (s *Server) RegisterHandler(conn Connection) error {
	fmt.Fprintf(conn, "Enter your name:\r\n")
	name, _ := helpers.GetInput(conn)
	_, ok := names[name]
	if ok {
		fmt.Fprintf(conn, "Sorry, but this name has already been taken.\r\n")
		// TODO: run this case again
	} else {
		fmt.Fprintf(conn, "Enter your password:\r\n")
		pass, _ := helpers.GetInput(conn)
		fmt.Fprintf(conn, "Confirm your password:\r\n")
		pass2, _ := helpers.GetInput(conn)
		if pass == pass2 {
			names[name] = pass
			// TODO: log the user in
			fmt.Fprintf(conn, "You have been successfully registered as %v!\r\n", name)
		}
	}

	return nil
}

func (s *Server) DBHandler(conn Connection) error {
	databases, err := s.DBClient.ListDatabaseNames(s.Context, bson.M{})
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, "databases list: %+v", databases)

	return nil
}

func (s *Server) ExitHandler(conn Connection) error {
	fmt.Fprintf(conn, "Goodbye!")
	conn.Close()

	return nil
}

func (s *Server) DefaultCommand(conn net.Conn, command string) error {
	fmt.Fprintf(conn, "What would you like to do?\r\n"+
		"[login, register, exit]\r\n")

	return nil
}

func (s *Server) StartGame(conn Connection) {
	// TODO: change user state to in-game
	//s.User.SwitchStatus()
	fmt.Fprintf(conn, "Your adventure starts here...\r\n")
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
