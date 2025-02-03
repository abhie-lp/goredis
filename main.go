package main

import (
	"fmt"
	"net"
)

var PORT = ":6379"

func main() {
	fmt.Println("Listening to port", PORT)

	// Create a new server
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		// read message from client
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: value.array[0].bulk})
	}

}
