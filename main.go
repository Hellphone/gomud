package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

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
		log.Fatalf("error connecting to database: %s", err)
	}

	var wg *sync.WaitGroup

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
		log.Fatalf("error registering commands: %s", err)
		return
	}

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// TODO: разобраться с Mutex, WaitGroup, Goroutines и тому подобным
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!

	// TODO: for now store only the list of online users in memory
	for {
		client, err := acceptConnection(ln, &s.Clients.Clients)
		// TODO: use logger to handle errors (find an example in One Platform code)
		if err != nil {
			log.Println(err)
		}

		// TODO: сделать команду, позволяющую выкинуть пользователя по имени (заодно разобраться с параметрами)
		// TODO: разобраться с состояниями пользователей (разрешать определённые команды только пользователям с определённым статусом)

		// TODO: stop the goroutine correctly (gracefully, using channels or context) when a connection is closed
		// TODO: check sync.Wg etc.
		//wg.Add(1)
		go func(client *server.Client, wg *sync.WaitGroup) {
			//defer wg.Done()
			for {
				err = handleInput(s, client)
				if err != nil {
					return
				}
			}
		}(client, wg)
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
		User: &models.User{
			Status: models.StatusOffline,
		},
	}
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
	message, parameter, err := helpers.GetCommandInput(c.Connection)
	if err != nil {
		return err
	}

	command, err := s.GetCommand(message)
	// TODO: maybe create a link between client and server error message to output them easier
	if err == models.ErrorCommandNotFound {
		fmt.Fprintf(c.Connection, "Sorry, but this command does not exist. Please repeat your input\r\n")
		return nil
	}
	if err != nil {
		return err
	}

	if command != nil {
		err = command(c, parameter)
		if err == models.ErrorUserNotFound {
			fmt.Fprintf(c.Connection, "Sorry, but this user can not be found\r\n")
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}
