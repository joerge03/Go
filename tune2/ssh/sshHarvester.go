package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%v", s)
}

func (s *strList) Set(str string) error {
	if str == "" {
		return fmt.Errorf("the arguments are empty")
	}

	argumentsData := strings.Split(str, ",")
	*s = append(*s, argumentsData...)
	return nil
}

var (
	in, out, address string
	logSSH           log.Logger
	targets          strList
)

func errorExit(s string, err error) {
	if err != nil {
		logSSH.Panicf("%v, %v \n", s, err)
	}
	log.Panicf("%v\n", s)
}

func init() {
	flag.StringVar(&in, "i", "", "-i for input file")
	flag.StringVar(&in, "input", "", "-input for input file")

	flag.StringVar(&out, "o", "", "-o for output file")
	flag.StringVar(&out, "out", "", "-out for output file")

	flag.StringVar(&address, "a", "", "-a address with a format of ip:port or [ipv6:4f::]:23")
	flag.StringVar(&address, "address", "-address", "-out address with a format of ip:port or [ipv6:4f::]:23")

	flag.Var(&targets, "t", "-t for targets separated with a comma (,)")
	flag.Var(&targets, "target", "-target for targets separated with a comma (,)")

	flag.Parse()

	logSSH = *log.New(os.Stdout, "[]", log.Ltime)

	if in == "" && len(targets) == 0 {
		errorExit("-i (input) and -t (targets) needs to be populated", nil)
	}
}

func main() {

}
