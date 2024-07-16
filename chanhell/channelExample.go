package main

import (
	"fmt"
	"math/rand"
	"time"
)

type jobTest int

func calcLong(val jobTest) int {
	duration := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(duration)

	fmt.Println("done in: ", duration)

	return int(val * 30)
}

func makeJobs() []jobTest {
	jobs := make([]jobTest, 0, 30)

	for i := 0; i < 30; i++ {
		jobs = append(jobs, jobTest(i))
	}
	return jobs
}

func runJobs(resultChan chan<- int, i jobTest) {
	resultChan <- calcLong(i)
}

type Test12 struct {
	test map[int]int
}

func main() {
	jobs := makeJobs()

	resultChan := make(chan int, 10)

	for i := range jobs {
		go runJobs(resultChan, jobs[i])
		fmt.Println(i)
	}

	resultCount := 0

	sum := 0

	for {
		result := <-resultChan
		sum += result

		resultCount++
		if resultCount == len(jobs) {
			break
		}
	}
	test := make([]Test12, 1)

	fmt.Println(test[0])

	fmt.Println(sum, " - ", resultCount)
}
