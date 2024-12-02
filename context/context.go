package main

import (
	"context"
	"fmt"
	"time"
)

func sampleOperation(ctx context.Context, str string, delay time.Duration) <-chan string {
	out := make(chan string)

	// wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(delay * time.Millisecond):
				out <- fmt.Sprintf("message: %v", str)
			case <-ctx.Done():
				out <- fmt.Sprintf("aborted %v", str)
				return
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	webServer := sampleOperation(ctx, "webServer", 500)

	microServices := sampleOperation(ctx, "microServices", 500)
	database := sampleOperation(ctx, "database", 500)

	go func() {
		fmt.Println("go func cancel ")
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	// Mainloop:
	for {
		select {
		case sdf := <-webServer:
			fmt.Println(sdf)
		case test := <-microServices:
			fmt.Println(test)
		case test := <-database:
			fmt.Println(test)
			// fmt.Println("asdf")
			// case test := <-database:
			// 	fmt.Println(test)
		}
	}

	// wg.Wait()

	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// _, open := <-database
	// if open {
	// }
}
