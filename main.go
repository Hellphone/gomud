package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/hellphone/gomud/domain/models"
	"github.com/hellphone/gomud/server"
	"github.com/hellphone/gomud/server/mongo"

	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to database...")
	ctx, client, err := mongo.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("error connecting to database:", err)
	}

	var clients []server.Client

	fmt.Println("Starting server...")
	ln, _ := net.Listen(cfg.Server.Protocol, cfg.Server.Port)

	// TODO: close connection if user is not active for 5-10 minutes by using goroutine
	s := &server.Server{
		Context:  ctx,
		DBClient: client,
	}

	// TODO: for now store only the list of online users in memory
	go func() {
		for {
			err := acceptConnections(ln, &clients)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
		}
	}()

	err = s.RegisterCommands()
	if err != nil {
		return
	}

	for _, c := range clients {
		go func(c *server.Client) {
			for {
				message, _ := bufio.NewReader(c.Connection).ReadString('\n')
				message = strings.TrimRight(message, "\r\n")
				command, err := s.GetCommand(message)
				// TODO: learn how to handle errors properly
				// TODO: get rid of commands causing errors when closing connection
				if err != nil {
					fmt.Printf("error: %v\n", err)
				}

				if command != nil {
					err = command(c.Connection)
					if err != nil {
						fmt.Printf("error: %v\n", err)
					}
				}
			}
		}(&c)
	}

	for {

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

func acceptConnections(ln net.Listener, clients *[]server.Client) error {
	conn, _ := ln.Accept()
	*clients = append(*clients, server.Client{
		Connection: conn,
		User:       nil,
	})
	_, err := fmt.Fprintf(conn, "Hello stranger! Welcome to GOMUD!\r\n"+
		"What would you like to do?\r\n"+
		"[login, register, exit]\r\n")
	if err != nil {
		return err
	}

	return nil
}
