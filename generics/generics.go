package main

import (
	"fmt"
)

type GenericIntType interface {
	~int | ~string
}

func Sum[t GenericIntType](s t) t {
	return s
}

type myS[t ~int] struct {
	test []t
}

func (m *myS[t]) myArr() t {
	max := m.test[0]
	for i, test := range m.test {
		if m.test[i] > max {
			max = test
		}
	}
	return max
}

type myI int

func main7() {
	testArr := myS[myI]{[]myI{23, 42, 32, 23}}

	// testString := myI(232)
	fmt.Println(testArr.myArr())
}
