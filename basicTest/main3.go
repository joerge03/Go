package main

import "fmt"



type Room struct {
	Name string
	IsCleaned bool
}

func checkIsCleaned(room [4]Room){
	for i := 0; i < len(room); i++{
		roomCopy := room[i]

		if roomCopy.IsCleaned{
			fmt.Println("This room is cleaned")
			}else{
			fmt.Println("asdf")
		}
	}
}

func main3(){
	rooms := []Room{
		{
			Name: "testsfd",
			IsCleaned: true,
		},
	}

	fmt.Print(rooms)
}