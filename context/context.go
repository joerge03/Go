package main

import (
	"context"
	"fmt"
	"time"
)

func sampleOperation(ctx context.Context, str string, delay time.Duration, done chan<- string) <-chan string {
	out := make(chan string)

	// wg.Add(1)
	go func() {
		// defer close(out)
		for {
			select {
			case <-time.After(delay * time.Millisecond):
				out <- fmt.Sprintf("message: %v", str)
			case <-ctx.Done():
				out <- fmt.Sprintf("aborted %v", str)
				done <- "done"
				return
			}
		}
	}()
	return out
}

func main() {
	done := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())

	webServer := sampleOperation(ctx, "webServer", 500, done)

	microServices := sampleOperation(ctx, "microServices", 500, done)
	database := sampleOperation(ctx, "database", 500, done)

	go func() {
		fmt.Println("go func cancel")
		time.Sleep(1000 * time.Millisecond)
		cancel()
	}()

Mainloop:
	for {
		// fmt.Println(<-webServer, <-microServices, <-database)
		select {
		case sdf := <-webServer:
			fmt.Println(sdf)
		case test := <-microServices:
			fmt.Println(test)
		case test := <-database:
			fmt.Println(test)
		case <-ctx.Done():
			// fmt.Println("done")
			// if d == "done" {
			break Mainloop
			// }
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
