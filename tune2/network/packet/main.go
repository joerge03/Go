package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	iface    = "enp6s0"
	snaplen  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	filter   = "tcp and port 80"
	devFound = false
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, d := range devices {
		if d.Name == iface {
			devFound = true
		}
	}

	if !devFound {
		log.Panic("Device does not exist")
	}

	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panic(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panic(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("running!")
	for packet := range source.Packets() {
		fmt.Println("running!")

		fmt.Printf("DATA: %+s\n PACKET: %+s\n", hex.Dump(packet.Data()), packet)
		// fmt.Printf("DATA: %+s\n PACKET: %+s\n", packet.)

	}

	// var (
	// 	eth     layers.Ethernet
	// 	ip4     layers.IPv4
	// 	ip6     layers.IPv6
	// 	tcp     layers.TCP
	// 	payload gopacket.Payload
	// 	tls     layers.TLS
	// 	udp     layers.UDP
	// )
	// parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &tcp, &payload, &tls, &ip6, &udp)

	// fmt.Println("test")
	// for {
	// 	data, _, err := handle.ZeroCopyReadPacketData()
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	decodedLayers := []gopacket.LayerType{}

	// 	if err := parser.DecodeLayers(data, &decodedLayers); err != nil {
	// 		log.Panic(err)
	// 	}

	// 	for _, layerType := range decodedLayers {
	// 		switch layerType {
	// 		case layers.LayerTypeEthernet:
	// 			fmt.Printf("src MAC : %v, dst MAC : %v\n", eth.SrcMAC, eth.DstMAC)
	// 		case layers.LayerTypeIPv4:
	// 			fmt.Printf("src IPv4 : %v, dst IPv4 : %v\n", ip4.SrcIP, ip4.DstIP)
	// 		case layers.LayerTypeTCP:
	// 			fmt.Printf("TCP:%+v\n", tcp.Contents)
	// 		case gopacket.LayerTypePayload:
	// 			fmt.Printf("%+v\n", payload)
	// 		default:
	// 			fmt.Printf("%+v\n", payload)
	// 		}
	// 	}
	// }

	// if err := handle.

}
