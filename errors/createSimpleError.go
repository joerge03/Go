package main

import (
	"errors"
	"fmt"
)

type Test struct {
	stuff []string
}

func (t *Test) Get(i int) (string, error) {
	if i > len(t.stuff) {
		return "", errors.New("Error the index is less than the stuff, cannot find the stuff with the given index")
	} else {
		fmt.Println("test")
		return t.stuff[i], nil
	}
}

func main3() {
	test := Test{}

	data, err := test.Get(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}
}
