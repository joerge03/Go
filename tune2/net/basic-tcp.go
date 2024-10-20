package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

var host, port string

func init() {
	flag.StringVar(&host, "host", "example.com", host)
	flag.StringVar(&port, "port", "80", port)
}

func main1() {
	flag.Parse()
	hostPort := net.JoinHostPort(host, port)
	fmt.Printf("test %v\n", hostPort)
	// fmt.Printf("go run . \r\n go run . ")

	conn, err := net.Dial("tcp", hostPort)

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	// req := "GET /  HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\n"

	// _, err = conn.Write([]byte(req))
	fmt.Fprint(conn, "req")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	reader := bufio.NewReader(conn)

	scanner := bufio.NewScanner(reader)

	for {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
	}

	// for {
	// 	response, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		if err.Error() == "EOF" {
	// 			log.Fatal(err)
	// 			break
	// 		}
	// 		fmt.Println("error reading the res", err)
	// 	}
	// 	fmt.Println(response)

	// }
}
