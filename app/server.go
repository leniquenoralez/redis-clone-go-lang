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
	store := make(map[string]string)
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
		go handleConnection(conn, store)
	}

}

func handleConnection(conn net.Conn, store map[string]string) {
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

		switch command {
		case "ping":
			handlePing(conn, value)
			break
		case "echo":
			handleEcho(conn, value)
			break
		case "set":
			key := value.Array()[1].String()
			store[key] = value.Array()[2].String()
			conn.Write([]byte("+OK\r\n"))
			break
		case "get":
			key := value.Array()[1].String()
			keyValue, exists := store[key]
			if !exists {
				conn.Write([]byte("+(error) Key does not exist in store!\r\n"))
				break
			}
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(keyValue), keyValue)))
			break
		}

	}
}
func handlePing(conn net.Conn, value Value) {
	args := value.Array()[1:]
	if len(value.Array()) >= 2 {
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
	} else {
		conn.Write([]byte("+PONG\r\n"))
	}
}

func handleEcho(conn net.Conn, value Value) {
	args := value.Array()[1:]
	if len(value.Array()) >= 2 {
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
	} else {
		conn.Write([]byte("+(error) ERR wrong number of arguments for command\r\n"))
	}
}
func handleSet(conn net.Conn, value Value) {

	readFile, err := os.Open("data.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}

	readFile.Close()
	println("handleSet")
}
func handleGet(conn net.Conn, value Value) {
	println("handleGet")
}
