package api

import (
	"fmt"
	"strings"

	"../structure"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type task struct {
	ID     string
	Name   string
	Result chan string
}

// PackageDispatch to dispatch the packet to goroutines
func PackageDispatch(packetSource *gopacket.PacketSource, todo interface{}) {
	DEBUG := false
	var currentState = make(map[string]string)
	tasks := make(chan task)
	final := make(chan bool)
	go func() {
		for {
			select {
			case t := <-tasks:
				currentState[t.ID] = t.Name
				//fmt.Println(currentState[t.ID])
				t.Result <- currentState[t.ID]
			}
		}
	}()

	go func() {
		i := 0
		for packet := range packetSource.Packets() {

			ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethernetLayer != nil {
				ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
				if DEBUG {
					fmt.Println("========")
					fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
					fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
					fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
				}
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)
				if DEBUG {
					fmt.Printf("From %s to %s\n ", ip.SrcIP, ip.DstIP)
					fmt.Println("Protocol: ", ip.Protocol)
				}
			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if DEBUG {
					fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
					fmt.Println("Sequence number: ", tcp.Seq)
				}

				// tcp packets analyze
				if tcp.Payload != nil {
					if packStr := ASCIIDecode(tcp.Payload); packStr != nil {
						//fmt.Printf("\033[1;33m%s\033[0m : "+" %s\n", "Pakets", packStr)
					}
				}
			}

			/*
				fmt.Println("All packet layers:")
				for _, layer := range packet.Layers() {
					fmt.Println("- ", layer.LayerType())
				}
			*/

			applicationLayer := packet.ApplicationLayer()
			if DEBUG {
				if applicationLayer != nil {
					fmt.Printf("%s\n", applicationLayer.Payload())
					if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
						fmt.Println("HTTP found!")
					}
					fmt.Println("========")
				}
			}

			currentTask := task{
				ID:     "1233",
				Name:   "123123",
				Result: make(chan string)}
			i++
			tasks <- currentTask
			<-currentTask.Result
		}
		fmt.Println("process", i, "packages")
		final <- true
	}()

	<-final

	for _, value := range GetCurrentPool().Packets {
		switch v := value.(type) {
		case structure.RealtimeFive:
			//fmt.Println(v)
			break
		case structure.RealtimeTrading:
			switch v.Type {
			case 1:
				fmt.Printf("\033[0;32mTime:%s ID:%s Deal:%.2f OrderCount:%d TotalCount: %d\033[0m\n", v.Timestamp, v.ID, v.Deal, int(v.OrderCount), int(v.TotalCount))
			case 2:
				fmt.Printf("\033[0;31mTime:%s ID:%s Deal:%.2f OrderCount:%d TotalCount: %d\033[0m\n", v.Timestamp, v.ID, v.Deal, int(v.OrderCount), int(v.TotalCount))
			}

			break
		}

	}

}
