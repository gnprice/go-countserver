package main

import "bufio"
import "fmt"
import "net"
import "os"
import "strconv"
import "strings"
import "time"

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
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s COUNT\n", os.Args[0])
		return
	}
	count, _ := strconv.Atoi(os.Args[1])

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		return
	}
	ich, och := make(chan string), make(chan string)
	go parseConnection(conn, ich, och)
	req := func (cmd string) (reply string) {
		ich <- cmd
		return <-och
	}

	start := time.Now()
	for i := 0; i < count; i++ {
		req("see a 1")
	}
	result := req("count a")
	tm := time.Since(start)
	qps := float64(count)*float64(time.Second)/float64(tm)

	fmt.Println("result:", result)
	fmt.Printf("%d iters in %.3fs, %.1f kqps\n",
		count, float64(tm)/float64(time.Second), qps / 1000.0)
}
