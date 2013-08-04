package main

import "bufio"
import "fmt"
import "net"
import "strings"

func serveConnection(ch chan string) {
	fmt.Println("connection made")
	for line := range ch {
		fmt.Println("got line:", line)
	}
	fmt.Println("connection closed")
}

func parseConnection(conn net.Conn, ch chan string) {
	bufreader := bufio.NewReader(conn)
	for {
		line, err := bufreader.ReadString('\n')
		if err != nil {
			close(ch)
			return
		}
		ch <- strings.TrimRight(line, "\n")
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
		ch := make(chan string)
		go parseConnection(conn, ch)
		go serveConnection(ch)
	}
}
