package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	hostUDP, portUDP string
)

func createAddr(host string, port string) (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))

	if err != nil {
		return nil, err
	}

	return addr, nil
}

func init() {
	flag.StringVar(&hostUDP, "host", "example.com", "host")
	flag.StringVar(&portUDP, "port", "80", "port")
}

func main2() {
	flag.Parse()
	addr, err := createAddr(hostUDP, portUDP)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	message := []byte("Hello, example.com!")
	_, err = conn.Write(message)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(conn)

	scanner := bufio.NewScanner(reader)

	for {
		if scanner.Scan() {
			fmt.Println("hello")
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
