package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func delay() {
	duration := time.Duration(rand.Intn(1000) * int(time.Millisecond))
	time.Sleep(duration)
}

type test struct {
	count int
	sync.Mutex
}

func serve(s *sync.WaitGroup, t *test, iteration int) {
	// delay()
	defer t.Unlock()
	defer s.Done()

	t.count += 1
	fmt.Println("done the: ", iteration)
	t.Lock()
}

func main() {
	testValue := test{}

	var swg sync.WaitGroup

	for i := 0; i < 50; i++ {
		go serve(&swg, &testValue, i)
		swg.Add(1)
	}

	fmt.Println("Waiting to load")
	swg.Wait()

	testValue.Lock()
	totalHits := testValue.count
	testValue.Unlock()

	fmt.Println("hits ", totalHits)
}
