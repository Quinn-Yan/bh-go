package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func printDevices() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		fmt.Println(device.Name)
		for _, addr := range device.Addresses {
			ones, _ := addr.Netmask.Size()
			fmt.Printf("\tIP:	%s/%v\n", addr.IP, ones)
		}
	}
}

var (
	snaplen  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	devFound = false
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <interface> <bpf filter, e.g. \"tcp and port 80\">", os.Args[0])
		os.Exit(0)
	}
	iface := os.Args[1]
	filter := os.Args[2]

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device '%s' not found\n", iface)
	}
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		//fmt.Println(packet)
		appLayer := packet.ApplicationLayer()
		if appLayer == nil {
			continue
		}
		payload := appLayer.Payload()

		var regexes = []*regexp.Regexp{
			regexp.MustCompile(`^(?i)user`),
			regexp.MustCompile(`^(?i)pass`),
		}
		for _, r := range regexes {
			if r.MatchString(string(payload)) {
				fmt.Printf(string(payload))
			}
		}

		//if bytes.Contains(payload, []byte("USER")) || bytes.Contains(payload, []byte("PASS")) {
		//	fmt.Print(string(payload))
		//}
	}
}
