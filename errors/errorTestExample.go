package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Time struct {
	hour,
	minute,
	seconds int
}

type TestParser struct {
	Msg string
}

func (t *TestParser) Error() string {
	return fmt.Sprintln("there is a print error")
}

func main() {
	r := bufio.NewReader(os.Stdin)
	sum := 0

	for {
		input, inputError := r.ReadString(' ')
		isEmptyInput := strings.TrimSpace(input) == ""

		if isEmptyInput {
			continue
		}

		num, convErr := strconv.Atoi(strings.TrimSpace(input))
		fmt.Println(strings.TrimSpace(input), "-")
		fmt.Println(input, "--")
		if convErr != nil {
			fmt.Println(convErr)
		} else {
			fmt.Println("asdf")
			sum += num
		}

		if inputError == io.EOF {
			break
		}
		if inputError != nil {
			fmt.Println("error:  Cannot read input |", inputError)
		}
	}
	fmt.Println(sum, "")
}
