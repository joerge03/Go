package main

import (
	"fmt"
	"log"
	"net"
)


func main(){
	conn, err := net.Dial("tcp","scanme.nmap.org:22")

	if err != nil {
		log.Panic(err)
	}

	data := make([]byte, 2048) 
	
	d, err := conn.Read(data)
	
	fmt.Printf("%v", string(data[:d]))

}