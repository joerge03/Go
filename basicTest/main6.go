package main

import "fmt"


const  (
	Online 	=  0
	Offline = 1
	// Darkstar = 2
	Maintenance = 2
	Retired = 3
)


// type StatusType struct  
// 	{Online  Online
// 	Offline  Offline
// 	// Darkstar : 2
// 	Maintenance Maintenance
// 	Retired  Retired}

func printTest(status map[string]int){
	fmt.Printf("There are %d of tests \n", len(status))
	type Status struct {
		Name string
		value int
	}

	statuses := make(map[string]int)
	for _, element := range status {
		// fmt.Println(i,element)
		switch element {
		case Online:
			statuses["Online"] +=1 
		case Offline:
			statuses["Offline"] +=1
		case Maintenance:
			statuses["Maintenance"] +=1
		case Retired:
			statuses["Retired"] +=1		
		}
	}
	fmt.Println("statuses: " ,statuses )
}

func main6(){

	testData := []string{"darkstar", "lonestar","pepero","ongbakonawa"}

	statuseses := make(map[string]int)
	
	for _,element := range testData {
		statuseses[element] = Online

		if element == "darkstar"{
			statuseses[element] = Maintenance
		}
		// fmt.Println(element, "elemento")
	}
	printTest(statuseses)
	fmt.Println(statuseses)

	
	
	
}