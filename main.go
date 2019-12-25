package main

import (
	"fmt"

	"./packets/api"
	s "./packets/api/structure"
)

func main() {
	f := api.TCPPcapng("./pcapng/packet_20191218_2337__00001_20191218091121.pcapng")
	packetSource := api.PCAParse(f)

	packets := s.NewPackets()
	fmt.Println(packets)
	fmt.Println(packetSource)

	/*
		api.PackageDispatch(packetSource, func() {
			fmt.Println("do somthing to package")
		})*/
}
