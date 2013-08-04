package main

import "bufio"
import "fmt"
import "net"
import "strings"

func parseConnection(conn net.Conn, ich chan string, och chan string) {
	bufreader := bufio.NewReader(conn)
	for {
		cmd := <-ich
		conn.Write([]byte(cmd+"\n"))

		line, err := bufreader.ReadString('\n')
		if err != nil {
			close(och)
			return
		}
		och <- strings.TrimRight(line, "\n")
	}
}

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		return
	}
	ich, och := make(chan string), make(chan string)
	go parseConnection(conn, ich, och)
	req := func (cmd string) (reply string) {
		ich <- cmd
		reply = <-och
		return
	}
	req("see a 1")
	fmt.Println(req("count a"))
}
