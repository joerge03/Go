package main

import (
	"flag"
	"fmt"
)

func main1() {

	port := 500
	flag.IntVar(&port, "port", 200, "port")

	ip := flag.String("ip", "192.168.24.2", "ip")

	verbose := flag.Bool("verbose", true, "verbose")

	flag.Parse()
	if *verbose {
		fmt.Println("chooo chooo")
	}
	fmt.Printf(`ip is %v \n`, *ip)
	fmt.Printf(`result :  ip %v , port: %v \n`, *ip, port)
}
