package main

import (
	"errors"
	"fmt"
)

const (
	low = iota
	medium
	high
)

type priorityQueue[p comparable, t any] struct {
	items    map[p][]any
	priority []p
}

func NewPriorityQueue[p comparable, t any](prop []p) priorityQueue[p, t] {
	return priorityQueue[p, t]{items: make(map[p][]any), priority: prop}
}

func (P *priorityQueue[x, y]) add(p x, value y) {
	P.items[p] = append(P.items[p], value)
}

func (P *priorityQueue[p, t]) next() (p, error) {
	for i, item := range P.items {
		if len(item) > 0 {
			fmt.Println("test")
			P.items[i] = item[1:]
			// for y, prioItem := range item {
			// }
			return i, nil
		}
	}
	return P.priority[0], errors.New("No items")
}

func main223() {
	test := NewPriorityQueue[int, int]([]int{low, medium, high})
	test.add(low, 23)
	test.add(low, 52)
	test.add(low, 12)
	test.add(medium, 50)
	test.add(medium, 50)
	test.add(high, 60)
	test.add(high, 70)
	fmt.Println(test)
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test.next())
	fmt.Println(test)
}
