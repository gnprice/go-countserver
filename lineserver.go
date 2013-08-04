package main

import "bufio"
import "fmt"
import "net"
import "strings"

func serveConnection(ich chan string, och chan string) {
	fmt.Println("connection made")
	for line := range ich {
		fmt.Println("got line:", line)
		och <- fmt.Sprintf("Thanks for saying %q!", line)
	}
	fmt.Println("connection closed")
}

func parseConnection(conn net.Conn, ich chan string, och chan string) {
	bufreader := bufio.NewReader(conn)
	for {
		line, err := bufreader.ReadString('\n')
		if err != nil {
			close(ich)
			return
		}
		ich <- strings.TrimRight(line, "\n")
		reply := <-och
		conn.Write([]byte(reply+"\n"))
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
		ich, och := make(chan string), make(chan string)
		go parseConnection(conn, ich, och)
		go serveConnection(ich, och)
	}
}
