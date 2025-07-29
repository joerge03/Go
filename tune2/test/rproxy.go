package main

import (
	"fmt"
	"os"
)

var rand1 string
var rand2 string

func main() {
	// test := []byte("sasdfasdf")
	os.Create("./test/testt/test.txt")

	fmt.Println(rand1)
	fmt.Println(rand2)
	// io.ReadFull()
}
