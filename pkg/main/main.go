package main

import (
	"fmt"
	"regexp"
)

var EmailExpr *regexp.Regexp

func init() {
	test, err := regexp.Compile(`.+@.+\..+`)
	if err != nil {
		panic("failed to compile regular expression")
	}

	fmt.Println(err)
	EmailExpr = test

	fmt.Println("regular")
}

func IsValidEmail(email string) bool {
	return EmailExpr.Match([]byte(email))
}

func main() {
	fmt.Println(IsValidEmail("testjbhjhb"))
	fmt.Println(IsValidEmail("email@gmail.com"))
	fmt.Println(IsValidEmail("email..@gma"))
	fmt.Println(IsValidEmail("testinadafiguguaka"))
}
