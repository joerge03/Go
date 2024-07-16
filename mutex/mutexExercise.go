package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"unicode"
)

type str struct {
	count int
	sync.Mutex
}

func getWords(s string) []string {
	return strings.Split(s, " ")
}

func getLettersCount(arrS string) int {
	count := 0

	for _, rune1 := range arrS {
		if unicode.IsLetter(rune1) {
			count++
		}
	}
	return count
}

func (s *str) serve(swg *sync.WaitGroup, fun func()) {
	defer s.Unlock()
	defer swg.Done()
	fun()
	s.Lock()
}

func main2332() {
	test := str{}

	var swg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			text := scanner.Text()
			words := getWords(text)
			for _, text1 := range words {
				go test.serve(&swg, func() {
					swg.Add(1)
					test.count += getLettersCount(text1)
				})
			}
			swg.Wait()
		} else {
			break
		}
		if scanner.Err() == io.EOF {
			break
		}
	}

	fmt.Println("count: ", test.count)
}
