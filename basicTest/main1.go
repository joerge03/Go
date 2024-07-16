package main

import "fmt"

func print(texts ...interface{}, text ...interface{}) {
	fmt.Println()
	fmt.Println(text...)
	fmt.Println()
}

func main1() {
	var nameOne string = "tset"
	nameTwo := "NameTwo"
	var nameThree string

	nameFour := "@"

	fmt.Print(nameOne, nameTwo, nameThree, nameFour, "\n")

	nameOne = "1"
	nameTwo = "2"
	nameThree = "3"
	fmt.Print(nameOne, nameTwo, nameThree, "\n")

	var numberOne float64 = 1
	numberTwo := 2
	var numberThree int
	numberFour := 5132

	fmt.Print(numberOne, numberTwo, numberThree, numberFour, "\n")
	fmt.Printf("my age is %v \n", numberFour)
	fmt.Printf("my age is %q \n", nameOne)
	fmt.Printf("nameOne is type %T \n", nameOne)
	fmt.Printf("my age is %0.1f \n", numberOne)
	printMessage := fmt.Sprintf("my age is %0.1f \n", numberOne)
	println(printMessage)
}
