package main

import (
	"fmt"
	"time"
	"unicode"
)

func Sleep() {
	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		fmt.Println(i)
	}
}

func main13231534() {
	data := []rune{'s', 'd', 'f', 'g', 'h'}

	capitalize := []rune{}

	toUpper := func(str rune) {
		capitalize = append(capitalize, unicode.ToUpper(str))
	}

	for i := range data {
		go toUpper(data[i])
	}

	time.Sleep(1000 * time.Millisecond)

	fmt.Printf("%v\n", string(capitalize))
}
