package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type cfg struct {
	fileLoc string
	Domain  string
	Server  string
	Worker  int
}

type Result3 struct {
	ServerAddress string
	Hostname      string
}

func Serve() *cfg {
	cfg := new(cfg)
	flag.StringVar(&cfg.fileLoc, "a", "wordlist.txt", "wordlist location, default ./wordlist.txt")
	flag.StringVar(&cfg.Server, "s", "8.8.8.8:53", "server")
	flag.StringVar(&cfg.Domain, "d", "example.com", "server")
	flag.IntVar(&cfg.Worker, "w", 5, "workers")
	flag.Parse()
	return cfg
}

func resultsGatherer(results *[]Result3, gatherer <-chan []Result3) {
	for g := range gatherer {
		*results = append(*results, g...)
	}
}

func main() {
	cfg := Serve()
	validate(*cfg)

	var results []Result3
	fqdns := make(chan string, 1024)
	gatherer := make(chan []Result3)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for range 15 {
		fmt.Println("+1")
		wg.Add(1)
		go worker2(fqdns, cfg.Server, gatherer, &wg, ctx)
	}

	// go resultsGatherer(&results, gatherer)
	go func() {
		for g := range gatherer {
			results = append(results, g...)
		}
	}()

	err := feedWorker(fqdns, cfg.Domain, cfg.fileLoc)
	if err != nil {
		log.Panic(err)
	}

	wg.Wait()
	fmt.Println("test")
	close(gatherer)

	fmt.Printf("%+v\n", results)

}

func LookupC(fqdn, addr string, client *dns.Client) ([]string, error) {
	res := new([]string)
	msg := dns.Msg{}

	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	r, _, err := client.Exchange(&msg, addr)
	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answer found")
	}

	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.CNAME); ok {
			*res = append(*res, a.Target)
		}
	}

	return *res, nil
}

func LookupA(fqdn, addr string, client *dns.Client) ([]string, error) {
	res := new([]string)
	msg := dns.Msg{}

	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	r, _, err := client.Exchange(&msg, addr)

	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answers found")
	}

	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.String())
			*res = append(*res, a.String())
		}
	}

	return *res, nil
}

func lookup2(fqdn, addr string, client *dns.Client) []Result3 {
	res := new([]Result3)
	fqdns := fqdn
	fmt.Println("lookup2")

	for {
		cnames, err := LookupC(fqdns, addr, client)
		if err == nil || len(cnames) >= 1 {
			fqdn = cnames[0]
			continue
		}

		a, err := LookupA(fqdns, addr, client)
		if err != nil {
			fmt.Printf("lookup err: %v", err)
			break
		}
		for _, answers := range a {
			*res = append(*res, Result3{ServerAddress: answers, Hostname: fqdn})
		}
		break
	}
	// fmt.Printf("%+v\n", *res)
	return *res
}

func worker2(fqdns <-chan string, address string, gatherer chan<- []Result3, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer fmt.Println("done2")
	client := &dns.Client{
		Timeout: 30 * time.Second,
	}
	fmt.Println("2")

	for fqdn := range fqdns {
		fmt.Println(fqdn)
		select {
		case <-ctx.Done():
			return
		default:
			results := lookup2(fqdn, address, client)
			gatherer <- results
		}
	}
}

func feedWorker(fqdns chan string, domain string, fileLoc string) error {
	defer fmt.Println("done1")
	defer close(fqdns)
	f, err := os.Open(fileLoc)
	if err != nil {
		log.Panic(err, "unable to open the file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%v.%v", scanner.Text(), domain)
	}

	return scanner.Err()
}

func validate(c cfg) {
	if _, err := os.Stat(c.fileLoc); os.IsNotExist(err) {
		log.Fatalf("Provided file does not exist")
	}
}
