// chitter.go
// A simple messaging server written in Go
// Written by James Whang

package main

import (
	"fmt"
	"net"
	"os"
    "strings"
    "bufio"
    "strconv"
    "log"
    "io"
)

type Client struct {
	connection net.Conn
	id string
	ch chan string
}

var done = make(chan bool)

func handleConnection(conn net.Conn, msgChan chan<- string, addClientChan chan<- Client, rmchan chan<- Client, clientID string) {
	//reader := bufio.NewReader(conn)
	defer conn.Close()

	client := Client {
		connection: conn,
		id: clientID,
		ch: make(chan string),
	}

	addClientChan <- client


	//ch := make(chan string)
	//addClientChan <- Client{conn, ch}
	// io.WriteString(conn, fmt.Sprintf("Welcome "))

	go client.ReadLinesInto(msgChan)
	client.WriteLinesFrom(client.ch)

}

func (c Client) ReadLinesInto(ch chan <- string) {
	reader := bufio.NewReader(c.connection)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if line == "whoami:\n" {
			ch <- fmt.Sprintf("%s: %s: %s", c.id, "chitter", c.id)
		} else {
			if strings.Contains(line, ":") {
				msgArray := strings.SplitN(line, ":", 2)

				if strings.TrimSpace(msgArray[0]) == "all" {
					ch <- fmt.Sprintf("%s: %s: %s", "all", c.id, msgArray[1])
				} else {
					ch <- fmt.Sprintf("%s: %s: %s", msgArray[0], c.id, msgArray[1])
				}
			} else {
				ch <- fmt.Sprintf("%s: %s: %s", "all", c.id, line)
			}
		}

		
		//ch <- fmt.Sprintf("%s: %s", c.id, line)
	}
}

func (c Client) WriteLinesFrom(ch <- chan string) {
	for msg := range ch {
		msgArray := strings.SplitN(msg, ":", 3)
		target := strings.TrimSpace(msgArray[0])

		if target == "all" {
			messages := []string{msgArray[1], ": ", msgArray[2]}
			_, err := io.WriteString(c.connection, strings.Join(messages, ""))
			if err != nil {
				return
			}
		} else if target == c.id {
			messages := []string{msgArray[1], ": ", msgArray[2]}
			_, err := io.WriteString(c.connection, strings.Join(messages, ""))
			if err != nil {
				return
			}
		}

		//_, err := io.WriteString(c.connection, msg)
		
	}
}

func initIdChan(idchan chan string) {
	for i:=0; i < 100; i++ {
		idchan <- strconv.Itoa(i)
	}
}

func handleMessage(msgChan chan string, clientChan <- chan Client, rmchan <- chan Client) {
	clients := make(map[net.Conn]chan <- string)
	for {
		select {
		case msg:= <- msgChan:
			for _, ch := range clients {
				go func (mch chan <- string) {mch <- msg}(ch)
			}
		case client:= <-clientChan:
			clients[client.connection] = client.ch
		case client := <- rmchan:
			delete(clients, client.connection)
		}
	}
	for msg := range msgChan {
		log.Printf("new message %s", msg)
	}
}



func main() {
	// usage & argument sanitization
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run chitter [port_num]")
        os.Exit(1)
    }

	port := os.Args[1]

	//msgChan := make(chan string, 100) // channel for communication	
	idChan := make(chan string, 100) // channel for clientID
	// TODO: What happens if there are more than 100 clients..?

	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
        fmt.Println("Failed to connect to port" + port)
        os.Exit(1)
	}

	msgChan := make(chan string)
	addClientChan := make(chan Client)
	removeClientChan := make(chan Client)


	initIdChan(idChan)

	go handleMessage(msgChan, addClientChan, removeClientChan)


	for {
		conn, err := ln.Accept()
		if err == nil {
			clientID := <- idChan
			go handleConnection(conn, msgChan, addClientChan, removeClientChan, clientID)
		}
	}
}
