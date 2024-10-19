package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strList) Set(str string) error {
	if len(str) <= 0 {
		return errors.New("str is nil")
	}

	*s = strings.Split(str, ",")

	return nil
}

var (
	ip, port strList
	verbose  bool
)

var wg sync.WaitGroup

func init() {
	flag.Var(&ip, "ip", "ip")
	flag.Var(&port, "port", "port")
	flag.BoolVar(&verbose, "verbose", true, "verbose")
}

func Permutations(ipList strList, portList strList, c chan<- string) {
	defer close(c)
	defer fmt.Println("close")

	for _, ipStr := range ipList {
		for _, portStr := range portList {
			c <- net.JoinHostPort(ipStr, portStr)
		}
	}
}

func hack(url string, verbose bool) {
	defer wg.Done()
	if url != "" {
		fmt.Printf("hacked : %s \n ", url)
	}

	file, _ := os.Create("test.txt")

	fmt.Fprintf(file, "asdf")

	if verbose {
		fmt.Println("potek nayan")
	}
}

func main() {

	flag.Parse()

	c := make(chan string)
	go Permutations(ip, port, c)
	// fmt.Println(flag.NFlag(), "nflag")

	for {
		select {
		case t, ok := <-c:
			if !ok {
				wg.Wait()
				fmt.Println("done")
				return
			}
			wg.Add(1)
			go hack(t, verbose)
		}
	}

}
