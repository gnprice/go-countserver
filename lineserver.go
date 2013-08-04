package main

import "bufio"
import "fmt"
import "net"

func handleConnection(conn net.Conn) {
	fmt.Println("got a connection!", conn)
	bufreader := bufio.NewReader(conn)
	for {
		line, err := bufreader.ReadString('\n')
		if err != nil {
			fmt.Printf("error on %s: %s\n", conn, err)
			return
		}
		fmt.Printf("read a line from %s: %q\n", conn, line)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen failed:", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept failed:", err)
			continue
		}
		go handleConnection(conn)
	}
}
