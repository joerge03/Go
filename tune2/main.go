package main

import "fmt"


func main(){
	// fmt.Printf("%v", strings.Repeat("-",30))
	var str  any


	str = 1

	switch str.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	case bool:
		fmt.Println("bool")
	}
	
}