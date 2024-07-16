package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line, path string, linNum int) Result {
	return Result{line, linNum, path}
}

func FindInFile(path string, find string) *Results {
	results := Results{make([]Result, 0)}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	scanner := bufio.NewScanner(file)

	lineNum := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), find) {
			result := NewResult(scanner.Text(), path, lineNum)
			results.Inner = append(results.Inner, result)
		}
		lineNum++
	}
	if len(results.Inner) != 0 {
		return &results
	}
	return nil
}
