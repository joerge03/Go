package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func main() {
	var msg dns.Msg

	fqdn := dns.Fqdn("augo.pw")

	msg.SetQuestion(fqdn, dns.TypeA)
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
		if record, ok := r.(*dns.A); ok {

			fmt.Printf("ip %v, header %v\n", record.A, record.Hdr)
		}
	}
}
