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

	fmt.Println("Starting server...")
	ln, _ := net.Listen(cfg.Server.Protocol, cfg.Server.Port)
	conn, _ := ln.Accept()

	defer conn.Close()

	// TODO: store context elsewhere
	s := &server.Server{
		Connection: conn,
		Context:    ctx,
		DBClient:   client,
	}

	err = s.RegisterCommands()
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(s.Connection, "Hello stranger! Welcome to GOMUD!\r\n"+
		"What would you like to do?\r\n"+
		"[login, register, exit]\r\n")
	if err != nil {
		return
	}

	for {
		message, _ := bufio.NewReader(s.Connection).ReadString('\n')
		message = strings.TrimRight(message, "\r\n")
		command, err := s.GetCommand(message)
		// TODO: learn how to handle errors properly
		// TODO: get rid of commands causing errors when closing connection
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		if command != nil {
			err = command()
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		}
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
