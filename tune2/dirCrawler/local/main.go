package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
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

func walkF(path string, f os.DirEntry, errx error) error {

	info, err := f.Info()
	if err != nil {
		return err
	}
	threeDaysAgo := time.Now().Add(-3 * 24 * time.Hour)

	for _, r := range regexed {
		if r.MatchString(f.Name()) && info.ModTime().After(threeDaysAgo) {
			fmt.Printf("[HIT] time: %v: %v \n", time.Now().Sub(info.ModTime()), path)
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
