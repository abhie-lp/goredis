package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string  // Data type carried by value
	str   string  // Value of string received from simple strings
	num   int     // Value of integer received
	bulk  string  // Value of string received from bulk strings
	array []Value // Holds all the values received from arrays
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	// Read one byte at a time till we reach '\r' and return without last 2 bytes
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY: // *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
		return r.readArray()
	case BULK: // $5\r\nhello\r\n
		return r.readBulk()
	default:
		fmt.Println("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	v := Value{}
	v.typ = "array"

	// Read length of the array which is the next byte
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// Iterate over the array
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// Add parsed value to array
		v.array[i] = val
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	// $5\r\nhello\r\n
	v := Value{}

	v.typ = "bulk"
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// read trailing CRLF
	r.readLine()

	return v, nil
}
