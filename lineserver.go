package main

// import "bufio"
import "fmt"
import "net"

func handleConnection(conn net.Conn) {
	fmt.Println("got a connection!", conn)
	buf := make([]byte, 32)
	n, _ := conn.Read(buf)
	fmt.Printf("read some bytes: %d %q\n", n, buf[:n])
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
