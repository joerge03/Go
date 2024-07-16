package main

import (
	"fmt"
	"strings"
)

func printIt(function func(str string), str string) {
	function(strings.ToUpper(str))
}

type calculateDiscount func(price float64) float64

func calculateDiscountFunc(price float64) float64 {
	discount := 0.1
	if price > 100 {
		discount += 0.1
	}
	return discount
}

func calculatePrice(price float64, discount calculateDiscount) float64 {
	return price - (price * discount(price))
}

func testFunction() (function func(str string)) {
	return func(str string) {
		fmt.Println(str)
	}
}

func main213152144() {
	// printIt(testFunction(), "test")

	fmt.Println(calculatePrice(300, calculateDiscountFunc))
}
