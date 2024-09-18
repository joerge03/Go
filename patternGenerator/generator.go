package main

import (
	"fmt"
	"math/rand"
)

func generate(min, max int) <-chan int {
	numbers := make(chan int, 5)

	go func() {
		for {
			numbers <- rand.Intn(max-min) + min
		}
	}()
	return numbers
}

func generateNumbers(amount, max, min int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i <= amount; i++ {
			out <- (rand.Intn(max-min+1) + min)
		}
		close(out)
	}()

	return out
}

func main() {
	numbers := generateNumbers(25, 20, 1)

	for {
		number1, open := <-numbers

		if !open {
			break
		}
		fmt.Println(number1)
	}

	fmt.Println(<-numbers)
}
