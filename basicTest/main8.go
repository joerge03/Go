package main

import "fmt"

type Status bool
type Items []string

const (
	Activated Status = true
	Deactivated Status = false
)

type Tag struct {
	Items 
	Status Status
}

func setTags(tag *Tag, status Status){
	if status {
		tag.Status = Activated
	}else {
		tag.Status = Deactivated
	}
}

func checkout(tags *[]Tag){
	for i := range *tags {
		(*tags)[i].Status = false
		fmt.Println((*tags)[i])
	}


}

func main8(){


	tags := []Tag{
		{Items: Items{"megatron","oktopouspride"}, Status: Activated},
		{Items: Items{"charizard"}, Status: Activated},
		{Items: Items{"Jimmy neutron"}, Status: Activated},
		{Items: Items{"Kickbattowski"}, Status: Activated},
	}

	

	checkout(&tags)
	fmt.Println(tags)

	

	// fmt.Println(tags)
	// setTags(&tags[0], Deactivated)
	

}