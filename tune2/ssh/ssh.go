package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func scan(addr string , openPorts *[]string)  {
		defer wg.Done()
			
		if addr == "error"{
			return
		}
		conn, err := net.Dial("tcp",addr)
		if err != nil {
			fmt.Printf("scanned but failed:%v\n", addr)
			return
		}
		conn.Close()
		*openPorts = append(*openPorts,  addr)
	
		
}

func main(){

	file, err := os.Open("../generateFile/localhost.txt")

	if err != nil {
		log.Panic(err)
	}

	

	defer file.Close()

	scanner := bufio.NewScanner(file)
	c := make(chan string, 4096)
	strings := []string{}
	openPorts := []string{}

	for scanner.Scan() {
		c <- scanner.Text()
		test := <-c

		strings = append(strings,test)
	}

	
	sort.Slice(strings, func(i, j int) bool {
		_,portI,err := net.SplitHostPort(strings[i])
		if err != nil {
			return false
		}

		_, portJ, err:= net.SplitHostPort(strings[j])
		if err != nil {
			return false
		}

		port1,_ := strconv.Atoi(portI)
		port2,_ := strconv.Atoi(portJ)
		if port1 <= port2 {
			return true
		}else {
			return false
		}
	})
		// fmt.Printf("%v", strings)
		close(c)
		
		for _,addr := range strings {
			wg.Add(1)
			go scan(addr, &openPorts)
		}
		wg.Wait()
		
		fmt.Println("done")
		fmt.Printf("%v ", openPorts)
		
		
		
	
	// }(scanner,c)

	// wg.Add(1)

	// go func(ch <- chan string ){
		// for {
		// 	fmt.Println("test")
		// 	test := <-c
		// 	if test == "error" {
		// 		break
		// 	}

			
		// }
		// defer wg.Done()
	// }(c)
	// defer close(c)


	
	// wg.Wait()
}