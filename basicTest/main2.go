package main

import "fmt"

func main2(){
	// var ages2 [3]int = [3]int{3,25,534}

	var ages = [3]int {2,3,5}
	
	names := [4]string{"carl", "sheen", "conan"}

	fmt.Println(ages, len(ages))
	fmt.Println(names, len(names))


	// SLICES (use arrays under the hood)
	var scores = []int{10,5324,2,5}
	
	scores[1] = 34
	
	addedScores := append(scores, 35)
	fmt.Println(scores, addedScores)
}