package queue

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	q := New(5)

	for i := 0; i < 5; i++ {
		if len(q.items) != i {
			t.Errorf("Item length is incorrect, result: %v, Want: %v", len(q.items), i)
		}
		if !q.Add(1) {
			// fmt.Println("value", q.items)
			t.Errorf("capacity is max capacity %v, expected: < %v", 5, q.capacity)
		}
		// fmt.Println(q)
	}
}

func TestNext(t *testing.T) {
	q := New(5)
	q.Add(5)
	for i := 0; i < 5; i++ {
		num, isNext := q.Next()
		fmt.Println(num, isNext)
		if isNext && num > i {
			t.Errorf("Incorrect next value expected %v, %v, but got, %v,%v", i, true, num, isNext)
		}
		if !isNext {
			t.Errorf("expected q value : %v, got %v", len(q.items), num)
		}
		fmt.Println(q, "next test")
	}
}
