package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Type represents a Value type
type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

// Value represents the data of a valid RESP type.
type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

// String converts Value to a string.
//
// If Value cannot be converted, an empty string is returned.
func (val Value) String() string {
	if val.typ == BulkString || val.typ == SimpleString {
		return string(val.bytes)
	}

	return ""
}

// Array converts Value to an array.
//
// If Value cannot be converted, an empty array is returned.
func (val Value) Array() []Value {
	if val.typ == Array {
		return val.array
	}

	return []Value{}
}

// DecodeRESP parses a RESP message and returns a RedisValue
func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}
	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func decodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}
