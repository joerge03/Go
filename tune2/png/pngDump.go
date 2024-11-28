package main

import "fmt"

func main(){
	const pngHeader = "\x89PNG\r\n\x1a\n"

	fmt.Printf("%v", pngHeader)
}