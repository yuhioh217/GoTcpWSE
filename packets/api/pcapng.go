package api

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func isError(err error) {
	if err != nil {
		panic(err)
	}
}

// TCPPcapng to open the file from path string and get the return *os.File
func TCPPcapng(file string) *pcap.Handle {
	h, err := pcap.OpenOffline(file)
	isError(err)

	return h
}

// PCAParse to parse the data from file and save the data to slice that contain Pack struct data
// Port filter
func PCAParse(h *pcap.Handle) *gopacket.PacketSource {
	//port := uint16(8049)
	//filter := getFilter(port)
	/* Check file is exist or not */
	if h == nil {
		panic("Parse failed, the file is wrong, please provide the correct file")
	}
	if err := h.SetBPFFilter(getFilter(8049)); err != nil {
		panic(err)
	}

	packageSource := gopacket.NewPacketSource(h, h.LinkType())
	return packageSource
}

func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
