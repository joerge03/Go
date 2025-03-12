package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snaplen  = int32(320)
	promisc  = true
	timeout  = pcap.BlockForever
	filter   = "tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18"
	devFound = false
	results  = make(map[string]int)
)

func capture(iface, target string) {
	handler, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panic(err)
	}
	defer handler.Close()

	if err := handler.SetBPFFilter(filter); err != nil {
		log.Panic(err)
	}
	ps := gopacket.NewPacketSource(handler, handler.LinkType())

	for packet := range ps.Packets() {
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			fmt.Println("no network layer")
			continue
		}

		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			fmt.Println("no transport layer")
			continue
		}

		srcHost := networkLayer.NetworkFlow().Src().String()
		srcPort := transportLayer.TransportFlow().Src().String()
		if srcHost != target {
			// fmt.Println("srchost does not match target")
			continue
		}

		results[srcPort] += 1
	}

}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("need more args")
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	iface := os.Args[1]
	ip := os.Args[2]
	for _, d := range devices {
		if d.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panic("no devices found")
	}

	go capture(iface, ip)
	time.Sleep(1 * time.Second)

	ports := explode(os.Args[3])

	// can be improved and use concurrency
	for _, port := range ports {
		target := fmt.Sprintf("%v:%v", ip, port)
		fmt.Printf("Trying... %v with a port of %v\n", target, port)

		c, err := net.DialTimeout("tcp", target, 1*time.Second)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		defer c.Close()
	}
	time.Sleep(5 * time.Second)

	for port, confidence := range results {
		fmt.Printf(" PORT [%v] --- confidence [%v] \n", port, confidence)
	}
}

func explode(s string) []string {
	ret := make([]string, 0)

	ports := strings.Split(s, ",")
	for _, port := range ports {
		ret = append(ret, strings.TrimSpace(port))
	}

	return ret

}
