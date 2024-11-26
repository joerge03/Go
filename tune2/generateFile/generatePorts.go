package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	inputFile,
	address,
	
	fileName string
	from,
	to uint64
)


func init(){	
	flag.StringVar(&fileName, "fn", "", "Name for outfile/input file")
	flag.StringVar(&fileName, "fileName", "", "Name for outfile/input file")

	flag.StringVar(&inputFile, "if", "", "Input file location")
	flag.StringVar(&inputFile, "inputFile", "", "Input file location")

	flag.Uint64Var(&from, "f", 1, "(-f) is port from and can be combined with (-t, -to)")
	flag.Uint64Var(&from, "from", 1, "(-from) is port from and can be combined with (-t, -to)")

	flag.Uint64Var(&to, "t", 65535, "(-t) is port to and can be combined with (-f, -from)")
	flag.Uint64Var(&to, "to", 65535, "(-to) is port to and can be combined with (-f, -from)")

	flag.StringVar(&address, "a", "", "(-a, -address) XD")
	flag.StringVar(&address, "address", "", "(-a, -address) you know this man, cmon")

	flag.Parse()

	if fileName == "" && inputFile == ""{
		log.Panic("Provide atleast a fileName (-fn) or inputFile (-if) if you already have a file")
	}

	if address == "" {
		log.Panic("use (-a, -address) to provide a address")
	}	
}

func main(){
	var w io.Writer
	newFile, err:= os.Create(fmt.Sprintf("%v.txt", fileName))
	
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	// w := io.MultiWriter(providedFile)
	if inputFile != "" {
		inputW, err := os.OpenFile(inputFile,os.O_WRONLY|os.O_APPEND,0666)
		inputW.Write([]byte("\n"))

		if err != nil {
			log.Panic(err)
		}
		w = io.MultiWriter(inputW, newFile, os.Stdout)
	}else{	
		w = io.MultiWriter(newFile, os.Stdout)
	}
	fmt.Println()
	outputLog := log.New(w,"",0)

	for from <= to {
		outputLog.Printf("%v:%v", address, from)
		from++
	}
	
}