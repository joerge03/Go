package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/miekg/dns"
)

var (
	list     string
	domain   string
	srvrAddr string
	worker   int
)

type Result struct {
	ServerAddress string
	Hostname      string
}

func init() {
	flag.StringVar(&list, "l", "wordlist.txt", "-l enter your worldlist location")
	flag.StringVar(&domain, "d", "example.com", "-d enter your domain")
	flag.StringVar(&srvrAddr, "s", "8.8.8.8:53", "-s enter your server address default (8.8.8.8)")
	flag.IntVar(&worker, "w", 5, "enter your desired worker amount, default (5)")

	flag.Parse()
	if len(list) == 0 || len(domain) == 0 {
		log.Panic("please provide a domain or list")
	}
}

func LookupA1(d, addr string) ([]string, error) {
	var m dns.Msg
	var ips []string

	m.SetQuestion(dns.Fqdn(d), dns.TypeA)
	r, err := dns.Exchange(&m, addr)
	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answers found")
	}

	for _, rr := range r.Answer {
		if typea, ok := rr.(*dns.A); ok {
			fmt.Println(typea.String(), "asdfsadf")
			ips = append(ips, typea.String())
		}
	}
	return ips, nil
}

func LookupCname(d, addr string) ([]string, error) {
	var m dns.Msg
	var fqcname []string

	m.SetQuestion(dns.Fqdn(d), dns.TypeCNAME)
	r, err := dns.Exchange(&m, addr)
	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answers found")
	}

	for _, ans := range r.Answer {
		if rr, ok := ans.(*dns.CNAME); ok {
			fqcname = append(fqcname, rr.Target)
		}
	}

	return fqcname, nil

}

func lookup(fqdn, addr string) []Result {

	var results []Result
	cfqdn := fqdn
	for {
		cdn, err := LookupCname(cfqdn, addr)
		if err == nil && len(cdn) >= 1 {
			cfqdn = cdn[0]
			continue
		}
		ips, err := LookupA1(cfqdn, addr)
		if err != nil {
			break
		}
		for _, typea := range ips {
			fmt.Println(typea, "typea")
			results = append(results, Result{ServerAddress: typea, Hostname: fqdn})
		}
		break
	}

	return results

}

var wg sync.WaitGroup

func Worker(addr string, gather chan []Result, fqdn chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for s := range fqdn {
		fmt.Println("run!")
		results := lookup(s, addr)

		if len(results) > 0 {
			fmt.Println("done11111asdfasdfasdf")
			gather <- results
		}
	}
}

func main2() {

	var results []Result

	fqdns := make(chan string, 1024)
	gatherer := make(chan []Result)

	f, err := os.Open(list)
	if err != nil {
		log.Panic(err, "something wrong with list")
	}
	go func() {
		for res := range gatherer {
			results = append(results, res...)
		}
	}()

	defer f.Close()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		fmt.Println("+1")
		go Worker(srvrAddr, gatherer, fqdns, &wg)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), domain)
		if !scanner.Scan() {
			close(fqdns)
		}
	}

	wg.Wait()
	close(gatherer)

	fmt.Printf("results : %+v\n ", results)
}
