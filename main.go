package main

import (
	"fmt"

	"./packets/api"
	s "./packets/structure"
)

func isError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f := api.TCPPcapng("./pcapng/packet_20191218_2337__00001_20191218091121.pcapng")
	packetSource := api.PCAParse(f)

	packets, err := s.NewPackets("2337", "time")
	isError(err)
	fmt.Println(packets)
	fmt.Println(packetSource)

	/*
		api.PackageDispatch(packetSource, func() {
			fmt.Println("do somthing to package")
		})*/
}
