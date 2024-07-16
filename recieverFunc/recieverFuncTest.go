package main

import "fmt"

type Space struct {
	Occupied bool
}

type Lot struct {
	Spaces []Space
}

func occupySpace(lot *Lot, spaceLoc int) {
	lot.Spaces[spaceLoc-1].Occupied = true
}

func (lot *Lot) occupySpace(spaceLoc int) {
	lot.Spaces[spaceLoc-1].Occupied = true
}

func (lot *Lot) vacateSpace(spaceLoc int) {
	lot.Spaces[spaceLoc-1].Occupied = false
}

func RecirecieverFuncTest() {
	spaces := Lot{make([]Space, 20)}

	spaces.occupySpace(4)
	spaces.vacateSpace(4)

	fmt.Println(spaces)
}
