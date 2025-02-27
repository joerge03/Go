package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func main() {
	var msg dns.Msg

	fqdn := dns.Fqdn("www​.example.com")

	msg.SetQuestion(fqdn, dns.TypeCNAME)
	dnsMsg, err := dns.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		log.Panic(err, "failed dnsx")
	}

	if len(dnsMsg.Answer) < 1 {
		fmt.Println("no records")
		return
	}

	for _, r := range dnsMsg.Answer {
		fmt.Println(r.String())
		if record, ok := r.(*dns.CNAME); ok {

			fmt.Printf("ip %v, header %v\n", record.Target, record.Hdr)
		}
	}
}
