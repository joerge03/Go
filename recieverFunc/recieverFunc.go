package main

import "fmt"

type TextCo struct {
	X int
	Y int
}

func (test *TextCo) functionTest(x, y int) *TextCo {
	test.Y += y
	test.X += x
	return &TextCo{test.Y, test.X}
}

func recieverFunc() {
	test := TextCo{1, 2}

	test.functionTest(2, 4)

	fmt.Println(test)
}
