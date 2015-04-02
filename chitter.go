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

func handleConnection(conn net.Conn) {
    // Read the message from the connection channel 
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Failed reading")
	}
    conn.Write([]byte(msg + "\n"))
    conn.Close()
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run chitter [port_num]")
        return
    }
	port := os.Args[1]
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
        fmt.Println("Failed to connect to port" + port)
        return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// error
		}
		go handleConnection(conn)
	}
}
