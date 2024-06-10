package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const defaultHost string = "127.0.0.1"
const defaultPort string = "10440"

const MAX_READ_BUFFER_SIZE = 1024

func env(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

type Command struct {
	Value []byte
	Kind  string

	prettyValue string
}

func (c *Command) Pretty() string {
	if c.prettyValue == "" {
		c.prettyValue = strings.Replace(string(c.Value), "\r\n", "", -1)
	}

	return c.prettyValue
}

type ServerInfo struct {
	Host string
	Port string
}

func handleResponse(resp string, c Command) {
	if resp == "+Ok" {
		fmt.Println("Ok")
	} else if strings.HasPrefix(resp, "+") {
		// We retrieved the key
		fmt.Printf("Key '%s' => '%s'\n", c.Pretty()[4:], resp[1:])
	} else if strings.HasPrefix(resp, "-Error") {
		// Something went wrong
		fmt.Printf("Command '%s' returned an error\n", c.Pretty())
	} else {
		fmt.Printf("Key '%s' was not found\n", c.Pretty()[4:])
	}
}

func execute(info ServerInfo, c Command) {
	address := fmt.Sprintf("%s:%s", info.Host, info.Port)

	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println("Failed to connect", err)
		return
	}

	defer conn.Close()

	_, err = conn.Write(c.Value)
	if err != nil {
		fmt.Printf("Failed to send command '%s'\n", c.Pretty())
		return
	}

	response := make([]byte, MAX_READ_BUFFER_SIZE)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Printf("Failed to read response '%s'\n", c.Pretty())
		return
	}

	handleResponse(string(response[:n]), c)

}

func main() {
	info := ServerInfo{
		Host: env("KV_HOST", defaultHost),
		Port: env("KV_PORT", defaultPort),
	}

	commands := []Command{
		{[]byte("SET banana pijama\r\n"), "SET", ""},
		{[]byte("SET another value\r\n"), "SET", ""},
		{[]byte("SET more values\r\n"), "SET", ""},
		{[]byte("SET banana other pijama\r\n"), "SET", ""},
		{[]byte("GET banana\r\n"), "GET", ""},
		{[]byte("SET dummy some content\r\n"), "SET", ""},
		{[]byte("GET dummy\r\n"), "GET", ""},
		{[]byte("GET not.here\r\n"), "GET", ""},
		{[]byte("Banana\r\n"), "GET", ""},
	}

	for _, c := range commands {
		execute(info, c)
	}
}
