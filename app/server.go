package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	value  string
	expire int64
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	store := make(map[string]Data)
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

func handleConnection(conn net.Conn, store map[string]Data) {
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
			handleSet(conn, store, value)
			break
		case "get":
			handleGet(conn, store, value)
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
func handleSet(conn net.Conn, store map[string]Data, value Value) {

	valueArray := value.Array()
	key := valueArray[1].String()
	fmt.Printf("Args: %v %d", value.Array(), len(value.Array()))

	valueData := Data{
		value: value.Array()[2].String(),
	}

	argCount := len(value.Array())
	exprType := strings.ToLower(value.Array()[3].String())
	if argCount == 5 && exprType == "px" {
		val := value.Array()[4].String()
		exp, _ := strconv.ParseInt(val, 10, 64)
		fmt.Printf("\nexpires in: %d %d", time.Now().UnixMilli()+exp, exp)
		valueData.expire = time.Now().UnixMilli() + exp
	}

	store[key] = valueData
	conn.Write([]byte("+OK\r\n"))

}
func handleGet(conn net.Conn, store map[string]Data, value Value) {
	now := time.Now().UnixMilli()

	key := value.Array()[1].String()
	keyValue, exists := store[key]

	if !exists {
		conn.Write([]byte("+(error) Key does not exist in store!\r\n"))
		return
	}
	if now >= keyValue.expire {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(keyValue.value), keyValue.value)))
}
