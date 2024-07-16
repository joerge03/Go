package main

import "fmt"

type math uint

const (
	add math = iota
	mult
	sub
	div
)

func (test math) com(num, num1 uint) uint {
	switch test {
	case add:
		return num + num1
	case mult:
		return num * num1
	case sub:
		return num - num1
	case div:
		if num < num1 {
			return num % num1
		} else {
			return num / num1
		}
	default:
		return 0
	}
}

func main() {
	fmt.Println(div.com(4, 6))
}
