package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)


var (
    hostServer ,portServer string
)

func init(){
    flag.StringVar(&hostServer,"host", "localhost", "host")
    flag.StringVar(&portServer,"port", "1234", "port")
}

func handleConnectionLog(conn net.Conn) {
    addr := conn.RemoteAddr().String()

    defer fmt.Println("closed: ", addr)
    newConn := conn
    
    
    fmt.Println("tests")
    
    fmt.Fprint(newConn, "testsdasdsad")
    reader := bufio.NewReader(strings.NewReader("enlo"))
    _,err := io.Copy(conn,reader)
    if err != nil {
        log.Fatal(err, "copy error")
    }
}

func main(){
    flag.Parse()


    formattedAddr := net.JoinHostPort(hostServer, portServer)
    fmt.Println(formattedAddr)

    listener, err := net.Listen("tcp", formattedAddr)

    if err != nil {
        log.Fatal(err, "unable to connect")
    }

    for {
        conn, err  := listener.Accept() 
        if err != nil {
            log.Fatal(err, "errasdfasdf")
            break
        } 
        fmt.Print("test")
        go handleConnectionLog(conn)       
    }
    fmt.Print("tail")    
}