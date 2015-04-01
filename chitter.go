package main

import (
    "os"
    "fmt"
    "net"
)

func handleConnection(conn net.Conn) {
  //var msg = make([]byte, 0, 4096)
  var tmp = make([]byte, 0, 1024)
  msgSize, err := conn.Read(tmp)
  if err != nil {
    fmt.Println("Failed reading")
  }
  fmt.Println(msgSize)
  fmt.Println(tmp)
}

func main() {
  port := os.Args[1]
  ln, err := net.Listen("tcp", ":" + port)
  if err != nil {
    // error handle
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
    // error 
    }
    go handleConnection(conn)
  }
}

