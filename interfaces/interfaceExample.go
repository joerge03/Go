package main

import "fmt"

type   interface {
	Function1()
	Function2(x int) int
}

type TestType int

func (m TestType) Function1() {
	fmt.Println("func1")
}

func (m TestType) Function2(x int) int {
	return x + x
}

func executer(i FirstInterface) {
	i.Function1()
}

func main1() {
	m := TestType(1)

	executer(m)
	executer(&m)

	// m.Function1()
	// fmt.Println(m.Function2(2))
	// execute(m)
}
