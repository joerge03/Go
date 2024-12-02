package main

import (
	"fmt"
	"io"
	"os"
)


type Foo struct{}
type Bar struct{}

func (f *Foo) Read(p []byte)(int, error){
	fmt.Print("Foo input > ")
	
	return os.Stdin.Read(p)
}

func (b *Bar) Write(p []byte) (int, error){
	fmt.Println("Bar output > ")
	
	return os.Stdout.Write(p)
} 

func main1(){
	read := Foo{}
	write := Bar{}
	// in := []byte("")
	// in := make([]byte, 2049)

	
	
	// reader, _ :=read.Read(in)
	// writer, _ := write.Write(in)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Printf("bytes read %d, ", read)
	// writer, err:= write.Write(in[:reader])	
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Printf("bytes wrote %d, ", in[:writer], )	

	io.Copy(&write,&read)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Printf("%d", w)
	
}