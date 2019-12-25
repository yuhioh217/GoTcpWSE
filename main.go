package main

import (
	"fmt"

	"./packets/api"
)

func isError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f := api.TCPPcapng("./pcapng/packet_20191218_2337__00001_20191218091121.pcapng")
	packetSource := api.PCAParse(f)
	api.PackageDispatch(packetSource, func() {
		fmt.Println("After PackageDispatching")
	})
}
