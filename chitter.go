// chitter.go
// A simple messaging server written in Go
// Written by James Whang

package main

import (
	"fmt"
	"net"
	"os"
//    "strings"
    "bufio"
)

func handleConnection(conn net.Conn, msgChan chan string, idChan chan int) {
    // Read the message from the connection channel 
    reader := bufio.NewReader(conn)
    
    for {
    	msg, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Bye!")
			break
		}

    	// echo it back to the connected user
    	conn.Write([]byte(msg))
   		//size, err := writer.WriteString([]byte(msg))
    	if err != nil {
      	  fmt.Println("Bye!")
      	  break
    	}
    }
}

func main() {
	// usage & argument sanitization
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run chitter [port_num]")
        return
    }

	port := os.Args[1]

	msgChan := make(chan string) // channel for communication
	idChan := make(chan int, 100) // channel for clientID
	// TODO: What happens if there are more than 100 clients..?

	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
        fmt.Println("Failed to connect to port" + port)
        return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// error
		} else {
		    go handleConnection(conn, msgChan, idChan)
        }
	}
}
