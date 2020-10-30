package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

//ack and fin: 00010001 (0x11)
//ack: 0010000 (0x10)
//ack and psh: 00011000 (0x18)
//14th byte
//tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18

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
	promisc  = true
	timeout  = pcap.BlockForever
	devFound = false
	filter   = "tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18"
	results  = make(map[string]int)
)

func capture(iface, target string) {
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	fmt.Println("Capturing packets")
	for packet := range source.Packets() {
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue
		}
		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			continue
		}

		srcHost := networkLayer.NetworkFlow().Src().String()
		srcPort := transportLayer.TransportFlow().Src().String()

		if srcHost != target {
			continue
		}
		results[srcPort]++
	}
}

func explode(arg string) ([]string, error) {
	ports := strings.Split(arg, ",")
	for _, port := range ports {
		_, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
	}
	return ports, nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <interface> <target ip> <port1,2,3>", os.Args[0])
	}
	iface := os.Args[1]

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

	ip := os.Args[2]
	go capture(iface, ip)
	time.Sleep(1 * time.Second)

	ports, err := explode(os.Args[3])
	if err != nil {
		log.Panicln(err)
	}

	for _, port := range ports {
		target := fmt.Sprintf("%s:%s", ip, port)
		fmt.Println("Trying", target)
		c, err := net.DialTimeout("tcp", target, 1000*time.Millisecond)
		if err != nil {
			continue
		}
		c.Close()
	}
	time.Sleep(5 * time.Second)

	for port, confidence := range results {
		if confidence > 1 {
			fmt.Printf("Port %s open (confidence: %d)\n", port, confidence)
		}
	}
}
