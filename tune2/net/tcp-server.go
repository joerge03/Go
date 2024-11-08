package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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
    c := make(chan []byte)
    go readData(conn,c)
    go writeData(conn,c)
}

func readData(conn net.Conn, c chan <-  []byte) {
    addr := conn.RemoteAddr().String()
    buff := make([]byte, 2048)
    for {
        fmt.Printf("string %s\n", buff )
        // fmt.Printf("%v byte , %s\n bytes", r, buff[:r])
        r, err := conn.Read(buff)
        // fmt.Printf("%v byte , %s\n bytes", r, buff[:r])

        if err == io.EOF {
            fmt.Printf("addr %v has been closed, ", addr)
            break
        }
        
        if err != nil {
            fmt.Print(err, "error reading")
            log.Fatal(err,r)
        }
        // fmt.Printf("READ: %s from %v\n\n",buff[:r], addr)        
        c <- buff[:r]
    }
}

func writeData(conn net.Conn, c <-chan  []byte){
    addr :=  conn.RemoteAddr().String()

    var buff []byte    
    for {
        fmt.Println(addr)
        buff = <-c
        _, err := conn.Write(buff)

        if err == io.EOF {
            fmt.Printf("addr %v has been closed, ", addr)
            
        }
        if err != nil {
            fmt.Printf("there is a problem writing: %v", err)
            break
        }
        // fmt.Printf("Write %v to: %v\n\n addr %v\n", buff , w, addr)  
    }
}

func main(){
    flag.Parse()

    formattedAddr := net.JoinHostPort(hostServer, portServer)    

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
        fmt.Print("tail")    
    }
}