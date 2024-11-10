package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	desIp,
	desPort,
	proxyIp,
	proxyPort string
)

func init() {
	flag.StringVar(&proxyIp, "ip", "localhost", "ip")
	flag.StringVar(&proxyPort, "port", "1234", "port")
	flag.StringVar(&desIp, "desIp", "localhost", "desIp")
	flag.StringVar(&desPort, "desPort", "1234", "desPort")
}

func readSocket(conn *net.Conn, c chan<- []byte) {

	buf := make([]byte, 2048)
	address := (*conn).RemoteAddr().String()

	for {
		r, err := (*conn).Read(buf)
		fmt.Printf("Has been sent to: %v\n", address)
		if err != nil {
			log.Fatal("there's something wrong with reading the connection")
			return
		}
		fmt.Printf("read from %v\n", address)
		c <- buf[:r]
	}
}

func writeSocket(conn *net.Conn, c <-chan []byte) {
	address := (*conn).RemoteAddr().String()
	var buf []byte

	// buf := make([]byte
	for {
		buf = <-c
		r, err := (*conn).Write(buf)

		fmt.Printf("Written to %v\n", address)
		if err != nil {
			log.Fatal("There's something wrong writing data to connection")
			return
		}
		fmt.Printf("Log into %v\n", buf[:r])
	}
}

func forwardConnection(conn *net.Conn) {

	c2s := make(chan []byte, 2048)
	s2c := make(chan []byte, 2048)

	proxyAddress := net.JoinHostPort(proxyIp, proxyPort)

	proxyConn, err := net.Dial("tcp", proxyAddress)

	if err != nil {
		log.Fatalf("there's something wrong dialing to proxy %v\n", err)
	}

	go readSocket(conn, c2s)
	go writeSocket(&proxyConn, c2s)
	go readSocket(conn, s2c)
	go writeSocket(&proxyConn, s2c)
}

func main3() {
	flag.Parse()

	clientAddress := net.JoinHostPort(desIp, desPort)

	listener, err := net.Listen("tcp", clientAddress)

	if err != nil {
		log.Fatalf("there's something wrong connection client, %v\n", err)
	}
	for {
		fmt.Println(" ******************** if you connect this to telnet or etc... This will infinitely loop ****************")
		c, err := listener.Accept()
		if err == io.EOF {
			log.Fatalf("closed from %v\n", c.RemoteAddr().String())
		}
		if err != nil {
			log.Fatalf("there's something wrong connectiong to %v\n", c.RemoteAddr().String())
		}
		fmt.Println("Client address", clientAddress)
		go forwardConnection(&c)
	}
}
