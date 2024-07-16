package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// func

func testFile(r bufio.Reader) int {
	sum := 0
	for {
		rd, err := r.ReadString('\n')
		if err == io.EOF {
			return sum
		}
		if err != nil {
			fmt.Println("Error:", err)
		}
		testNum, err := strconv.Atoi(rd[:len(rd)-1])
		if err != nil {
			fmt.Println("Error: ", err, testNum)
		}

	}
}

// func main123123() {
// 	files := []string{"phub.txt", "elgato.txt", "flix.txt", "test.txt"}
// }
