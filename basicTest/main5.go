package main

import "fmt"

func main5(){
	// slice := []string{"test","test1","test2"}

	// for _,element :=  range slice {

	// 	for _,ch := range element {

	// 		fmt.Printf("%q \n",  string(ch))
	// 	}

	// }

	mapTest := make(map[string]int)

	mapTest["sdaf"] = 23
	mapTest["test"] = 5

	// mapTest := map[string]int{
	// 	"asdf": 2,
	// 	"asdf2": 2232,
	// 	"asdf3": 23213,
	// 	"asd4": 4352,
	// 	"asdf23": 54376722,
	// 	"asdf321": 64823442,
	// }

	delete(mapTest, "sdaf")


	sdaf, exist := mapTest["sdaf"]

	if !exist{
		fmt.Println("sdaf does not exist")
		} 
	fmt.Println(sdaf)
	fmt.Println(exist)


	// for i,key := range mapTest {

	// 	if key == 2  || key == 23213{
	// 		mapTest[i] = 6
	// 	}

	// 	fmt.Println	(key)
		
	// }

	fmt.Println(mapTest)
}