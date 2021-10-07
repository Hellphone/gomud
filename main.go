package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hellphone/gomud/domain/models"
	"github.com/hellphone/gomud/helpers"
	"github.com/hellphone/gomud/server"
	"github.com/hellphone/gomud/server/mongo"

	"gopkg.in/yaml.v2"
)

var clients server.ClientList

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to database...")
	ctx, dbClient, err := mongo.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("error connecting to database:", err)
	}

	fmt.Println("Starting server...")
	ln, _ := net.Listen(cfg.Server.Protocol, cfg.Server.Port)

	// TODO: close connection if user is not active for 5-10 minutes by using goroutine
	// context should not be stored inside a struct type:
	// https://go.dev/blog/context-and-structs
	s := &server.Server{
		Context:  ctx,
		DBClient: dbClient,
		Clients:  &clients,
	}

	err = s.RegisterCommands()
	if err != nil {
		return
	}

	// TODO: for now store only the list of online users in memory
	for {
		client, err := acceptConnection(ln, &s.Clients.Clients)
		// TODO: use logger to handle errors (find an example in One Platform code)
		if err != nil {
			log.Println(err)
			return
		}

		// !!!!!
		// TODO: how to loop properly?
		// !!!!!
		go func(client *server.Client) {
			err = handleInput(s, client)
			if err != nil {
				log.Println(err)
				return
			}
		}(client)
	}
}

func getConfig() (*models.Config, error) {
	f, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg *models.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func acceptConnection(ln net.Listener, clients *[]server.Client) (*server.Client, error) {
	conn, _ := ln.Accept()
	client := server.Client{
		Connection: conn,
		User:       nil,
	}
	// TODO: delete client from the slice when closing connection
	*clients = append(*clients, client)

	_, err := fmt.Fprintf(conn, "Hello stranger! Welcome to GOMUD!\r\n"+
		"What would you like to do?\r\n"+
		"[login, register, exit]\r\n")
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func handleInput(s *server.Server, c *server.Client) error {
	for {
		message, err := helpers.GetInput(c.Connection)
		if err != nil {
			log.Fatal(err)
		}

		command, err := s.GetCommand(message)
		// TODO: learn how to handle errors properly
		// TODO: get rid of commands causing errors when closing connection
		if err != nil {
			// TODO: get rid of 'use of closed network connection' error
			log.Println(err)
		}

		if command != nil {
			err = command(c.Connection)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
