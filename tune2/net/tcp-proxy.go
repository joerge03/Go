package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	desIp,
	desPort,
	proxyIp,
	proxyPort string
)

func init(){
	flag.StringVar(&proxyIp, "ip","127.0.0.1", "ip")
	flag.StringVar(&proxyPort,"port", "80","port")
	flag.StringVar(&desIp,"desIp", "127.0.0.1","desIp")
	flag.StringVar(&desPort,"desPort", "80","desPort")
}

func readSocket(conn *net.Conn, c chan<- []byte){
		
	buf := make([]byte, 2048)
	address := (*conn).RemoteAddr().String()

	for {
		r, err := (*conn).Read(buf)
		fmt.Printf("Has been sent to: %v\n", address)
		if err != nil {
			log.Fatal("there's something wrong with reading the connection")			
			return 
		}
		c <- buf[:r]
	}
}


func writeSocket(conn *net.Conn, c <-chan []byte){
	address := (*conn).RemoteAddr().String()
	var buf []byte

	for {		
		buf =  <- c		
		r, err := (*conn).Write(buf)

		fmt.Printf("Written to %v\n", address)
		if err != nil {
			log.Fatal("There's something wrong writing data to connection")
			return 
		}
		fmt.Printf("Log into %v\n", buf[:r])
	}
}

func main3(){
	flag.Parse()
}