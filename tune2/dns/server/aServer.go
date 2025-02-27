package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

func main() {
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		var res dns.Msg
		res.SetReply(r)
		for _, q := range r.Question {
			fmt.Printf("%+v\n", q)
			a := &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP("127.0.0.1").To4(),
			}
			res.Answer = append(res.Answer, a)
		}
		w.WriteMsg(&res)
	})
	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
