package main

import (
	"fmt"
	"time"
)

type controlType int

const (
	exit = iota
	exitOk
)

type job struct {
	data int
}

type result struct {
	result int
	job    job
}

func doubler(job <-chan job, res chan<- result, conType chan controlType) {
	for {
		select {
		case msg := <-conType:
			{
				switch msg {
				case exit:
					{
						fmt.Println("exiting")
						conType <- exitOk
						return
					}
				default:
					panic("Error invalid input")
				}
			}
		case j := <-job:
			{
				time.Sleep(499 * time.Millisecond)
				res <- result{result: j.data * 2, job: j}
			}
		}
	}
}

func main5436345() {
	jobs := make(chan job, 101)
	results := make(chan result, 50)
	control := make(chan controlType)
	go doubler(jobs, results, control)

	for i := 0; i <= len(jobs); i++ {
		jobs <- job{i}
	}

	for {
		// control <- exitOk
		select {
		case result := <-results:
			fmt.Println(result, "results")
		case <-time.After(500 * time.Millisecond):
			fmt.Println(results, "result testtststsdt")
			control <- exit
			<-control
			// fmt.Println(<-control)
			return
		}
		// fmt.Println(<-control, "results")
	}

	// channel := make(chan int, 2)

	// func() { channel <- 1 }()
	// func() { channel <- 2 }()
	// go func() { channel <- 3 }()

	// first := <-channel
	// second := <-channel
	// third := <-channel

	// fmt.Println(first, second, third)
}
