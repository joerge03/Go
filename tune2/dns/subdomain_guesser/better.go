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

type Config struct {
	Wordlist   string
	Domain     string
	ServerAddr string
	Workers    int
}

type Result1 struct {
	ServerAddress string
	Hostname      string
	IP            string // More explicit than storing dns.A string
}

func main1() {
	cfg := parseFlags()
	validateConfig(cfg)

	var results []Result1
	result1s := make(chan Result1)
	fqdns := make(chan string)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start worker1 pool
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go worker1(ctx, &wg, cfg.ServerAddr, fqdns, result1s)
	}

	// Start result1 processor
	result1WG := sync.WaitGroup{}
	result1WG.Add(1)
	go func() {
		defer result1WG.Done()
		for res := range result1s {
			results = append(results, res)
		}
	}()

	// Feed workers with FQDNs
	if err := feedWorkers(cfg.Domain, cfg.Wordlist, fqdns); err != nil {
		log.Fatalf("Error feeding workers: %v", err)
	}
	// Wait for all workers to finish
	wg.Wait()

	// Close result1s channel after all workers are done
	close(result1s)

	// Wait for result1 processor to finish
	result1WG.Wait()
	fmt.Printf("%+v\n", results)
}

func parseFlags() Config {
	var cfg Config
	flag.StringVar(&cfg.Wordlist, "l", "wordlist.txt", "Wordlist file path")
	flag.StringVar(&cfg.Domain, "d", "example.com", "Target domain")
	flag.StringVar(&cfg.ServerAddr, "s", "8.8.8.8:53", "DNS server address")
	flag.IntVar(&cfg.Workers, "w", 5, "Number of workers")
	flag.Parse()
	return cfg
}

func validateConfig(cfg Config) {
	if _, err := os.Stat(cfg.Wordlist); os.IsNotExist(err) {
		log.Fatalf("Wordlist file not found: %s", cfg.Wordlist)
	}
}

func feedWorkers(domain, wordlist string, fqdns chan<- string) error {
	defer close(fqdns) // Important: close channel when done

	file, err := os.Open(wordlist)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := scanner.Text()
		fqdns <- fmt.Sprintf("%s.%s", subdomain, domain)
	}
	return scanner.Err()
}

func worker1(ctx context.Context, wg *sync.WaitGroup, serverAddr string, fqdns <-chan string, result1s chan<- Result1) {
	defer wg.Done()

	client := &dns.Client{
		Timeout: 5 * time.Second, // Add timeout
	}

	for fqdn := range fqdns {
		select {
		case <-ctx.Done():
			return // Graceful shutdown
		default:
			if ips, err := resolve(fqdn, serverAddr, client); err == nil {
				for _, ip := range ips {
					result1s <- Result1{
						Hostname:      fqdn,
						ServerAddress: serverAddr,
						IP:            ip,
					}
				}
			}
		}
	}
}

func resolve(fqdn, serverAddr string, client *dns.Client) ([]string, error) {
	var ips []string

	// Follow CNAME chain
	for {
		cname, err := lookupCNAME(fqdn, serverAddr, client)
		if err != nil || cname == "" {
			break
		}
		fqdn = cname
	}

	// Get A records
	records, err := lookupA(fqdn, serverAddr, client)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		ips = append(ips, record.A.String())
	}
	return ips, nil
}

func lookupA(fqdn, serverAddr string, client *dns.Client) ([]*dns.A, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)

	r, _, err := client.Exchange(m, serverAddr)
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS error: %s", dns.RcodeToString[r.Rcode])
	}

	var records []*dns.A
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			records = append(records, a)
		}
	}

	if len(records) == 0 {
		return nil, errors.New("no A records found")
	}
	return records, nil
}

func lookupCNAME(fqdn, serverAddr string, client *dns.Client) (string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)

	r, _, err := client.Exchange(m, serverAddr)
	if err != nil {
		return "", err
	}

	if r.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("DNS error: %s", dns.RcodeToString[r.Rcode])
	}

	for _, ans := range r.Answer {
		if cname, ok := ans.(*dns.CNAME); ok {
			return cname.Target, nil
		}
	}
	return "", errors.New("no CNAME records found")
}
