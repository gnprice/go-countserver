package main

import "bufio"
import "fmt"
import "net"
import "strconv"
import "strings"

type state_t struct {
	count uint64
	max uint64
	sum uint64
}

var state = make(map[string](*state_t))

func getState(key string) (*state_t) {
	s := state[key]
	if s == nil {
		s = new(state_t)
		state[key] = s
	}
	return s
}

func serveConnection(ich chan string, och chan string) {
	for line := range ich {
		fields := strings.Fields(line)
		switch cmd := strings.ToLower(fields[0]); cmd {
		case "see":
			s := getState(fields[1])
			v, _ := strconv.ParseUint(fields[2], 10, 64)
			s.count++
			if v > s.max {
				s.max = v
			}
			s.sum += v
			fmt.Printf("saw %s %d: now %d %d %d\n",
				fields[1], v, s.count, s.max, s.sum)
			och <- "ok"
		case "count":
			s := getState(fields[1])
			fmt.Printf("state of %s: now %d %d %d\n",
				fields[1], s.count, s.max, s.sum)
			och <- strconv.FormatUint(s.count, 10)
		default:
			och <- "error"
		}
	}
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
