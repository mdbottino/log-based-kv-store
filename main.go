package main

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/mdbottino/log-based-kv-store/filesystem"
	"github.com/mdbottino/log-based-kv-store/store"
)

const MAX_BUFFER_SIZE = 1024

type Command struct {
	Name  string
	Key   string
	Value string
}

var InvalidCommandError = errors.New("invalid command")

func parseCommand(command string) (Command, error) {
	parts := strings.Split(command, " ")

	if len(parts) < 2 {
		return Command{}, InvalidCommandError
	}

	sanitizedCommandName := strings.ToLower(parts[0])

	if sanitizedCommandName == "get" {
		return Command{"GET", strings.Replace(parts[1], "\r\n", "", -1), ""}, nil
	}

	if sanitizedCommandName == "set" {
		if len(parts) < 3 {
			return Command{}, InvalidCommandError
		}

		return Command{"SET", parts[1], strings.Replace(parts[2], "\r\n", "", -1)}, nil
	}

	return Command{}, errors.New("unknown command")
}

func handleGetCommand(input string, conn net.Conn, s *store.Store) {
	command, err := parseCommand(input)
	if err != nil {
		_, err = conn.Write([]byte("-Error"))
		if err != nil {
			fmt.Println("Failed to respond", err)
		}
		return
	}

	val, err := s.Get(command.Key)

	if err != nil {
		_, err = conn.Write([]byte("-"))
		if err != nil {
			fmt.Println("Failed to respond", err)
		}
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("+%s", val)))
	if err != nil {
		fmt.Println("Failed to respond", err)
	}

}

func handleSetCommand(input string, conn net.Conn, s *store.Store) {
	command, err := parseCommand(input)
	if err != nil {
		_, err = conn.Write([]byte("-Error"))
		if err != nil {
			fmt.Println("Failed to respond", err)
		}
		return
	}

	err = s.Set(command.Key, command.Value)
	if err != nil {
		_, err = conn.Write([]byte("-Error"))
		if err != nil {
			fmt.Println("Failed to respond", err)
		}
		return
	}

	_, err = conn.Write([]byte("+Ok"))
	if err != nil {
		fmt.Println("Failed to respond", err)
	}

}

func handle(conn net.Conn, s *store.Store) {
	defer conn.Close()

	buffer := make([]byte, MAX_BUFFER_SIZE)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read the request", err)
		return
	}

	command := string(buffer[:n])

	if strings.HasPrefix(command, "GET") {
		handleGetCommand(command, conn, s)
	} else {
		handleSetCommand(command, conn, s)
	}
}

func main() {
	s := store.NewStore("./data", filesystem.FileSystem{})
	address := "127.0.0.1:10440"
	server, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind the server")
		panic(err)
	}
	fmt.Println("Listening on", address)
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming request", err)
			continue
		}

		fmt.Println("Handling a connection")
		go handle(conn, &s)
	}
}
