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
	"text/tabwriter"

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

	for range cfg.Worker {
		fmt.Println("+1")
		wg.Add(1)
		go worker2(fqdns, cfg.Server, gatherer, &wg, ctx)
	}

	// go resultsGatherer(&results, gatherer)
	gwg := sync.WaitGroup{}
	mut := &sync.Mutex{}
	gwg.Add(1)
	go func() {
		gwg.Done()
		for g := range gatherer {
			mut.Lock()
			results = append(results, g...)
			mut.Unlock()
		}
	}()

	err := feedWorker(fqdns, cfg.Domain, cfg.fileLoc)
	if err != nil {
		log.Panic(err)
	}

	wg.Wait()
	fmt.Println("test")
	close(gatherer)
	gwg.Wait()

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)

	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.ServerAddress)
	}

	w.Flush()

}

func LookupC(fqdn, addr string) ([]string, error) {
	var res []string
	msg := dns.Msg{}

	// client := dns.Client{
	// 	Net: "udp",
	// }
	// client := &dns.Client{
	// 	Net: "udp", // Use "tcp" for TCP queries
	// 	Dialer: &net.Dialer{
	// 		Timeout: 5 * time.Second, // Set a timeout for the dialing operation
	// 		// Optional: Specify a local address
	// 		LocalAddr: &net.UDPAddr{
	// 			IP:   net.ParseIP("192.168.1.10"), // Replace with your local IP address
	// 			Port: 0,                           // Let the system choose an available port
	// 		},
	// 	},
	// }
	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	r, err := dns.Exchange(&msg, addr)
	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answer found")
	}

	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.CNAME); ok {
			res = append(res, a.Target)
		}
	}
	return res, nil
}

func LookupA(fqdn, addr string) ([]string, error) {
	var res []string
	msg := dns.Msg{}
	// test := &dns.Client{
	// 	Net: "udp", // Use "tcp" for TCP queries
	// 	Dialer: &net.Dialer{
	// 		Timeout: 5 * time.Second, // Set a timeout for the dialing operation
	// 		// Optional: Specify a local address
	// 		LocalAddr: &net.UDPAddr{
	// 			IP:   net.ParseIP("192.168.1.10"), // Replace with your local IP address
	// 			Port: 0,                           // Let the system choose an available port
	// 		},
	// 	},
	// }

	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	r, err := dns.Exchange(&msg, addr)

	if err != nil {
		return nil, err
	}

	if len(r.Answer) < 1 {
		return nil, errors.New("no answers found")
	}

	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.String(), "a string")
			res = append(res, a.String())
		}
	}

	return res, nil
}

func lookup2(fqdn, addr string) []Result3 {
	var res []Result3
	fqdns := fqdn
	fmt.Println("lookup2")

	for {
		cnames, err := LookupC(fqdns, addr)
		if err == nil || len(cnames) >= 1 {
			fqdns = cnames[0]
			continue
		}
		a, err := LookupA(fqdns, addr)
		if err != nil {
			fmt.Printf("lookup err: %v", err)
			break
		}
		for _, answers := range a {
			res = append(res, Result3{ServerAddress: answers, Hostname: fqdn})
		}
		break
	}
	// fmt.Printf("%+v\n", res)
	return res
}

func worker2(fqdns <-chan string, address string, gatherer chan<- []Result3, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer fmt.Println("worker done!")
	// client := &dns.Client{Net: "tcp"}

	fmt.Println("2")

	for fqdn := range fqdns {
		fmt.Println(fqdn)
		select {
		case <-ctx.Done():
			fmt.Println("ctx <-done")
			return
		default:
			fmt.Println("worker run")
			results := lookup2(fqdn, address)
			gatherer <- results
		}
	}
}

func feedWorker(fqdns chan string, domain string, fileLoc string) error {
	defer close(fqdns)
	f, err := os.Open(fileLoc)
	if err != nil {
		log.Panic(err, "unable to open the file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println("scanner run!")
		fqdns <- fmt.Sprintf("%v.%v", scanner.Text(), domain)
	}

	return scanner.Err()
}

func validate(c cfg) {
	if _, err := os.Stat(c.fileLoc); os.IsNotExist(err) {
		log.Fatalf("Provided file does not exist")
	}
}
