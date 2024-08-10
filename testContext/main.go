package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	timeNow := time.Now()
	ctx := context.Background()

	val, err := fetchUserData(ctx, 34)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result: ", val)
	fmt.Println("Time: ", time.Since(timeNow))
}

type Response struct {
	num int
	err error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	enlo := ctx.Value("test")
	fmt.Println(enlo)
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	responseChannel := make(chan Response)

	go func() {
		num, err := fetchThatThirdy()

		responseChannel <- Response{num, err}
	}()

	for {
		select {
		case res := <-responseChannel:
			return res.num, res.err
		case <-ctx.Done():
			return 0, errors.New("took too long to load")
		}
	}
}

func fetchThatThirdy() (int, error) {
	time.Sleep(200 * time.Millisecond)

	return 654, nil
}
