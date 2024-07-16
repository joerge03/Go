package main

import "fmt"


type Counter struct {
	Hits int
}

func increment(num *Counter){
	num.Hits += 1
	fmt.Printf("%p", num)
} 

func replace(str *string,newStr string){
	*str = newStr
	
}
	
func main7(){
		
	// num := Counter{1}
	strings := []string{"Asdf","gdas"}
	
	replace(&strings[0],"ongbakunawa")
	replace(&strings[1],"ongbakunawa2")
	// increment(&num)

	fmt.Printf("%p",strings )
		
	// fmt.Println(num.Hits)
		
}