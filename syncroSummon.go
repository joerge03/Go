package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SyncData struct {
	inner map[string]int
	mutex sync.Mutex
}

func (s *SyncData) insert(str string, val int) {
	s.mutex.Lock()
	s.inner[str] = val
	defer s.mutex.Unlock()
}

func (s *SyncData) get(str string) int {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	data := s.inner[str]

	return data
}

func main5213() {
	// data := SyncData{inner: make(map[string]int)}
	// data.insert("asdf", 1)
	// data.insert("asdff", 135)

	// fmt.Println(data.get("asdf"))
	// fmt.Println(data.get("asdff"))

	// var waitGroup sync.WaitGroup

	// sum := 0
	// for i := 0; i < 20; i++ {

	// 	waitGroup.Add(1)
	// 	value := i

	// 	fmt.Println(i, "test")
	// 	go func() {
	// 		defer waitGroup.Done()
	// 		sum += value
	// 		fmt.Println(value, "func")
	// 	}()

	// }

	// waitGroup.Wait()

	// fmt.Println(waitGroup)

	var waitGroup sync.WaitGroup
	value := 0

	for i := 0; i < 20; i++ {
		waitGroup.Add(1)
		value++
		go func() {
			fmt.Println(value)
			defer func() {
				waitGroup.Done()
				// fmt.Println(waitGroup)
				value--
				fmt.Println("func")
			}()
			duration := time.Duration(rand.Intn(500) * int(time.Millisecond))
			fmt.Println("duration : ", duration)
			time.Sleep(duration)
		}()
		waitGroup.Wait()
		fmt.Println("testing to wait")
	}
}
