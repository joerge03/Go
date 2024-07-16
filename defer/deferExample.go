package main

import "fmt"

func one() {
	fmt.Println("one")
}

func two() {
	fmt.Println("two")
}

func three() {
	fmt.Println("three")
}

func main() {
	fmt.Println("test")
	defer one()
	defer two()

	fmt.Println("test end")
	defer three()

	fmt.Println("test end")
}
