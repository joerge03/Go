package main

import (
	"fmt"
)

func main() {
	// slice := []string{"test","test1","test2"}

	// for _,element :=  range slice {

	// 	for _,ch := range element {

	// 		fmt.Printf("%q \n",  string(ch))
	// 	}

	// }

	mapTest := make(map[struct{string}][]int)

	mapTest[struct{string}{"asdf"}] = []int{1, 2}
	// mapTest[] = []int{3, 4}

	// mapTest := map[string]int{
	// 	"asdf": 2,
	// 	"asdf2": 2232,
	// 	"asdf3": 23213,
	// 	"asd4": 4352,
	// 	"asdf23": 54376722,
	// 	"asdf321": 64823442,
	// }

	// delete(mapTest, "sdaf")

	sdaf, exist := mapTest[struct{string}{"asdf"}]

	if !exist {
		fmt.Println("sdaf does not exist")
	}
	fmt.Println(sdaf,exist,mapTest)
	// fmt.Println(exist)

	// for i,key := range mapTest {

	// 	if key == 2  || key == 23213{
	// 		mapTest[i] = 6
	// 	}

	// 	fmt.Println	(key)

	// }

	// fmt.Println(mapTest)
}
