package queue

import "fmt"

type Queue struct {
	items    []int
	capacity int
}

func (q *Queue) Add(count int) bool {
	if len((*q).items)+count > (*q).capacity {
		fmt.Println("Can't add, full capacity")
		return false
	}
	(*q).items = append((*q).items, make([]int, count)...)
	(*q).capacity += count
	return true
}

func New(val int) Queue {
	return Queue{make([]int, 0, val), val}
}

func (q *Queue) Next() (int, bool) {
	if len(q.items) > 0 {
		item := q.items[0]
		q.items = q.items[1:]

		return item, true
	} else {
		return 0, false
	}
}

// func New(queue *Queue, count int) []Queue {
// 	if len(queue.items)+count > queue.capacity {
// 		println("can't add, max cap")
// 	}
// 	queue.
// 	return newQueue
// }
