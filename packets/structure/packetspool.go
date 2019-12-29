package structure

import (
	"fmt"
	"sync"
)

var (
	lock         *sync.Mutex = &sync.Mutex{}
	poolInstance *PacketsPool
)

// PacketsPool is the pool that will lot of processing packets struct in it.
type PacketsPool struct {
	Packets []interface{}
}

// GetInstance to create singleton struct if there is no instance
func GetInstance() *PacketsPool {
	if poolInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if poolInstance == nil {
			poolInstance = &PacketsPool{}
			fmt.Println("instance...")
		}
	}
	return poolInstance
}

// ResetPacketsPool the instance to nil and empty the Packets interface in PacketsPool
func ResetPacketsPool() *PacketsPool {
	poolInstance = nil
	return GetInstance()
}

// AddPackets will append the packets to slice
func (p *PacketsPool) AddPackets(i interface{}) {
	p.Packets = append(p.Packets, i)
}
