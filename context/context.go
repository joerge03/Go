package main

import (
	"context"
	"fmt"
	"time"
)

func sampleOperation(ctx context.Context, str string, delay time.Duration) <-chan string {
	out := make(chan string)

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
	context1 := context.Background()

	ctx, cancel := context.WithCancel(context1)

	webServer := sampleOperation(ctx, "webServer", 500)

	microServices := sampleOperation(ctx, "microServices", 700)
	database := sampleOperation(ctx, "database", 800)

	go func() {
		time.Sleep(900 * time.Millisecond)
		cancel()
		fmt.Println("testing")
	}()

Mainloop:
	for {
		select {
		case test := <-webServer:
			fmt.Println(test)
		case test := <-microServices:
			fmt.Println(test, "canceled")
		case test := <-database:
			fmt.Println(test)
			fmt.Println(<-database)
		case <-ctx.Done():
			fmt.Println("test")
			cancel()
			break Mainloop

			// case test := <-database:
			// 	fmt.Println(test)
		}
	}

	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// fmt.Println(<-database)
	// _, open := <-database
	// if open {
	// }
}
