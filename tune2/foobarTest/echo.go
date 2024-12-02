package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)


func echo(nc net.Conn){
	defer nc.Close()

	// b := make([]byte, 1024)
	// cmd:= "ls\n"
	
	// fmt.Printf("remote %s, local %v\n", nc.RemoteAddr().String(), nc.LocalAddr().String())
	// config := &ssh.ClientConfig{
	// 	User: "",
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.Password(""),
	// 	},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }	
	for {		
		w:= io.MultiWriter(nc,os.Stdout)
		r := io.MultiReader(nc)
		io.Copy(w,r)
		// fmt.Println("go")
	}
	
}

func main2(){
	
	l, err := net.Listen("tcp", "localhost:8082")
	if err != nil{
		log.Panic(err)
	}
		for {
			conn, err := l.Accept()
			// fmt.Println("test")
			if err != nil {
				log.Panic(err)
			}
			go echo(conn)
			fmt.Println("accept")
		}
	}