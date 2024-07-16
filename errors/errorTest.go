package main

import (
	"errors"
	"fmt"
)

type ErrorTest struct {
	Msg string
}

func (e *ErrorTest) Error() string {
	return fmt.Sprintf("error: %v", e.Msg)
}

func someFunc(t string) (string, error) {
	err := &ErrorTest{"test"}
	return t, err
}

func main1() {
	_, err := someFunc("test")
	if err != nil {
		inputError := &ErrorTest{"test"}

		fmt.Println(inputError)
		fmt.Println(err)

		if errors.Is(err, inputError) {
			fmt.Println("error is")
		} else {
			fmt.Println("error nah")
		}
	}
}
