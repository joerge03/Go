package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(conn *net.Conn) *Flusher {
	return &Flusher{w: bufio.NewWriter(*conn)}
}

func (f *Flusher) Write(b []byte) (int, error) {
	i, err := f.w.Write(b)
	if err != nil {
		return -1, err
	}
	err = f.w.Flush()

	if err != nil {
		return -1, err
	}

	return i, err
}

// cmd := exec.Command("/bin/sh", "-i", "-v")
func Handle(c *net.Conn, stdout *os.File) {
	os := runtime.GOOS

	var cmd *exec.Cmd

	if os == "linux" {
		fmt.Println("you are in linux")
		cmd = exec.Command("/bin/sh", "-i")
	} else {
		fmt.Println("you are in windows")
		cmd = exec.Command("cmd.exe")
	}

	rp, wp := io.Pipe()

	// test:= os.std

	// os.std
	// r := io.MultiReader(*c,,rp)
	w := io.MultiWriter(*c)

	// test := make([]byte, 2048)
	// (*c).Read(test)

	cmd.Stdin = *c
	// cmd.Stdin = *c
	// cmd.Stdout = NewFlusher(c)
	cmd.Stdout = wp

	go io.Copy(w, rp)
	// fmt.Printf("test %s",  test)

	// fmt.Println(n, "bytes")

	fmt.Println("run")
	err := cmd.Run()
	defer (*c).Close()

	if err != nil {
		fmt.Print(err)
	}
}

func main3() {

	l, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(c.RemoteAddr())
		Handle(&c, os.Stdout)
	}

	// cmd.Stdout = NewFlusher(&conn).w

	// buffedOutData := make([]byte, 2049)

	// b, err := conn.Read(buffedOutData)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// reader.Flush(buffedOutData[:b])

	// cmd.Stdin = conn
	// cmd.Stdout = conn
}
