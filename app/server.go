package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
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
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {

		value, err := DecodeRESP(bufio.NewReader(conn))

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}

		command := strings.ToLower(value.Array()[0].String())
		args := value.Array()[1:]

		switch command {
		case "ping":
			if len(value.Array()) >= 2 {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
			} else {
				conn.Write([]byte("+PONG\r\n"))
			}
			break
		case "echo":
			if len(value.Array()) >= 2 {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
			} else {
				conn.Write([]byte("+(error) ERR wrong number of arguments for command\r\n"))
			}
			break
		}

	}
}
