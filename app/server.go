package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleClientRequest(conn)
	}

}

func handleClientRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)

		argsRegex := regexp.MustCompile(`[\r\n]+`)
		bufferString := bytes.NewBuffer(buffer).String()
		args := argsRegex.Split(bufferString, -1)

		if len(args) == 0 || err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			os.Exit(1)
		}

		command := strings.TrimSpace(strings.ToLower(args[2]))
		switch command {
		case "ping":

			if len(args) > 4 {
				var value = "+" + strings.TrimSpace(strings.ToLower(args[4])) + "\r\n"
				conn.Write([]byte(value))
			} else {
				conn.Write([]byte("+PONG\r\n"))
			}
			break
		case "echo":
			if len(args) > 4 {
				var value = "+" + strings.TrimSpace(strings.ToLower(args[4])) + "\r\n"
				conn.Write([]byte(value))
			} else {
				conn.Write([]byte("+(error) ERR wrong number of arguments for command\r\n"))
			}
			break
		}

	}
}
