package main

import "fmt"

const (
	SmallLift = iota
	MediumLift
	LargeLift
)

const (
	Small = iota
	Medium
	Large
)

type Lift interface {
	CheckVehicle()
}

type Vehicle struct {
	Name string
	Size int
}

func (v Vehicle) CheckVehicle() {
	switch v.Size {
	case SmallLift:
		fmt.Println("- Going in small lift")
	case MediumLift:
		fmt.Println("- Going in medium lift")
	case LargeLift:
		fmt.Println("- Going in large lift")
	default:
		fmt.Println("- Too large or small to lift")
	}
}

func check(v Lift) {
	v.CheckVehicle()
}

func main() {
	vehicle := []Vehicle{
		{Name: "toy car", Size: Large},
		{Name: "hot wheels", Size: Medium},
		{Name: "batman car", Size: Small},
		{Name: "motor", Size: 4},
	}

	for _, v := range vehicle {
		check(v)
	}
}
