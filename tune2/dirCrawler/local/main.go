package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var regexed = []*regexp.Regexp{
	regexp.MustCompile(`(?i)user`),
	regexp.MustCompile(`(?i)pass(word)?`),
	regexp.MustCompile(`(?i)kdb`),
	regexp.MustCompile(`(?i)login`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)key`),
}

func walkF(path string, f os.DirEntry, err error) error {
	for _, r := range regexed {
		if r.MatchString(path) {
			fmt.Printf("[HIT] : %v\n", path)
		}
	}
	return nil
}

func main() {
	root := os.Args[1]
	if err := filepath.WalkDir(root, walkF); err != nil {
		log.Panic(err)
	}
}
