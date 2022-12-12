package main

import (
	"bytes"
	"fmt"
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

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	handleRequest(conn)
}

func handleRequest(conn net.Conn) {

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	data := strings.Split(bytes.NewBuffer(buf).String(), "$4\r\n")
	message := strings.TrimSpace(strings.ToLower(data[1]))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	if message == "ping" && len(data) > 2 {
		var value = "+" + strings.TrimSpace(strings.ToLower(data[2])) + "\r\n"
		conn.Write([]byte(value))
	} else {
		conn.Write([]byte("+PONG\r\n"))
	}

	conn.Close()
}
