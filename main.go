package main

import (
	"fmt"
	"time"

	"./packets/api"
)

func isError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	startTime := time.Now()
	f := api.TCPPcapng("./pcapng/packet_20191224_3019_00001_20191224090005.pcapng")
	packetSource := api.PCAParse(f)
	api.PackageDispatch(packetSource, func() {
		fmt.Println("After PackageDispatching")
	})
	uptime := time.Since(startTime)
	fmt.Println("Process time : ", uptime)
}
