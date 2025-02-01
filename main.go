package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Create a server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a AOF
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

  // Read the AOF
  aof.Read(func(value Value) {
    fmt.Println(value)
    command := strings.ToUpper(value.array[0].bulk)
    args := value.array[1:]

    handler, ok := Handlers[command]

    if !ok {
      fmt.Println("Invalid command: ", command)
      return
    }
    handler(args)
  })

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return

	}
	defer conn.Close()

	fmt.Println("Connection open...")

	for {
		resp := NewResp(conn)

		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("invalid request, expected array of length more than 1")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)

			writer.Write(Value{typ: "string", str: "OK"})
			continue
		}

		if command == "SET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
