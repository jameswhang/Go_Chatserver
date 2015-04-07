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
)

var done = make(chan bool)

func handleConnection(conn net.Conn, msgChan *chan string, clientID int) {
    // Read the message from the connection channel 
    reader := bufio.NewReader(conn)
    /*
    for {
    	// echo it back to the connected user
	   	//for i := 0; i < 2; i++ {
	   		go func () {
		        select {
			       	case msg1 := <-*msgChan:
			       		msgArray := strings.SplitN(msg1, ":", 2)
			       		target := strings.TrimSpace(msgArray[0])
			       		message := msgArray[1]

			       		fmt.Println(target)
			       		fmt.Println(message)

			       		
			       		if strconv.Itoa(clientID) == target {
			            	fmt.Println("received from channel", msg1)
			            	conn.Write([]byte(message))
			            }
			        default:
			        	msg2, err := reader.ReadString('\n')
			        	if err != nil {
			        		fmt.Println("Bye!")
			        		conn.Close()
			        		break
			        	}

			        	msgArray := strings.SplitN(msg2, ":", 2)
			        	target := strings.TrimSpace(msgArray[0])
			        	message := msgArray[2]
			        	fmt.Println(target)
			        	fmt.Println(message)
			        	
			        	if target == "all" {
			        		// Do something here to broadcast lol
			        	} else if message == msg2 {
			        		// Do something here to broadcast lol 
			        	} else if message == "whoami" {
			        		conn.Write([]byte(strconv.Itoa(clientID)))
			        	} else {
			        		fmt.Println("received from user", msg2, clientID)
		            		*msgChan <- msg2
			        	}
			    }
			    done <- true
		    }()
    	//}

   		//conn.Write([]byte(rcvMsg))
    } 
    */
    //reader := bufio.NewReader(conn)
    for {
    	go func () {
    		message , _:= reader.ReadString('\n') // TODO: ERROR CHECK
    		conn.Write([]byte(message))
    	}()
    }
    
}

func initIdChan(idchan chan int) {
	for i:=0; i < 100; i++ {
		idchan <- i
	}
}

func main() {
	// usage & argument sanitization
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run chitter [port_num]")
        return
    }

	port := os.Args[1]

	msgChan := make(chan string, 100) // channel for communication	
	idChan := make(chan int, 100) // channel for clientID
	// TODO: What happens if there are more than 100 clients..?

	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
        fmt.Println("Failed to connect to port" + port)
        return
	}

	initIdChan(idChan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// error
		} else {
			clientID := <- idChan
			s := []string{"Hello, ", strconv.Itoa(clientID), "\n"}
			welcomeMessage := strings.Join(s, "")
			conn.Write([]byte(welcomeMessage))
		    go handleConnection(conn, &msgChan, clientID)
		   // <-done
        }
	}
}
