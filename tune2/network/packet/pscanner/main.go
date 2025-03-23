package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snaplen = int32(320)
	promisc = true
	timeout = pcap.BlockForever
	filter  = "tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18"
	results = make(map[string]int)
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

func Dial(ip, port string, wg *sync.WaitGroup) {
	target := fmt.Sprintf("%v:%v", ip, port)
	defer wg.Done()
	fmt.Printf("Trying... %v with a port of %v\n", target, port)
	c, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer c.Close()
}

type ArgVar struct {
	ip    string
	port  string
	iface string
}

var argVar ArgVar

func init() {
	flag.StringVar(&argVar.ip, "ip", "", "enter IP seperated by comma or a file location (seperate it with newline)")
	flag.StringVar(&argVar.iface, "iface", "", "name of your network interface")
	flag.StringVar(&argVar.port, "p", "80", "enter PORTS seperated by comma or a file location (seperate it with newline)")

	flag.Parse()
	if len(argVar.ip) < 1 {
		log.Panic("Missing IP arg, fill out -ip")
	}
	if len(argVar.iface) < 1 {
		log.Panic("Missing Interface, fill out -iface")
	}
}

func isValidLocation(s string) bool {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		return false
	}
	return true
}

func ValidateAndHandleFile(loc string) []string {
	var ipFormatted []string
	if !isValidLocation(loc) {
		ipFormatted = explode(loc)
	} else {
		os, err := os.Open(loc)
		if err != nil {
			log.Panic(err)
		}
		buf := bufio.NewScanner(os)
		for buf.Scan() {
			ipFormatted = append(ipFormatted, buf.Text())
		}
	}
	return ipFormatted
}

func PortsWork(ips, ports string, wg *sync.WaitGroup) {
	formattedIps := ValidateAndHandleFile(ips)
	formattedPorts := ValidateAndHandleFile(ports)

	for _, ip := range formattedIps {
		for _, port := range formattedPorts {
			wg.Add(1)
			go Dial(ip, port, wg)
		}
	}
}

func isDeviceExist(dname string) bool {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	// ip := os.Args[2]
	for _, d := range devices {
		if d.Name == dname {
			return true
		}
	}
	return false
}

func main() {
	iface := argVar.iface
	ip := argVar.ip
	port := argVar.port
	wg := sync.WaitGroup{}

	if !isDeviceExist(iface) {
		log.Panicf("No device detected with %v name", iface)
	}

	go capture(iface, ip)
	time.Sleep(1 * time.Second)

	PortsWork(ip, port, &wg)

	wg.Wait()
	time.Sleep(5 * time.Second)

	for port, confidence := range results {
		fmt.Printf(" PORT [%v] --- confidence [%v] \n", port, confidence)
	}
}

func explode(s string) []string {
	var ret []string

	ports := strings.Split(s, ",")
	for _, port := range ports {
		ret = append(ret, strings.TrimSpace(port))
	}
	return ret

}
