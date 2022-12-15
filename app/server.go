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
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	for {

		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		re := regexp.MustCompile(`\$[0-9]+\r\n`)
		bufferString := bytes.NewBuffer(buf).String()
		args := re.Split(bufferString, -1)

		if len(args) == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			os.Exit(1)
		}

		message := strings.TrimSpace(strings.ToLower(args[1]))

		if message == "ping" && len(args) > 2 {
			var value = "+" + strings.TrimSpace(strings.ToLower(args[2])) + "\r\n"
			conn.Write([]byte(value))
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}

}

// func handleRequest(conn net.Conn) {

// 	conn.Close()
// }
