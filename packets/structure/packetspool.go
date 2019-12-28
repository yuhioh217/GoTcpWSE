package structure

import (
	"fmt"
	"sync"
)

var (
	lock     *sync.Mutex = &sync.Mutex{}
	instance *PacketsPool
)

// PacketsPool is the pool that will lot of processing packets struct in it.
type PacketsPool struct {
	Packets []interface{}
}

// GetInstance to create singleton struct if there is no instance
func GetInstance() *PacketsPool {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &PacketsPool{}
			fmt.Println("instance...")
		}
	}
	return instance
}

// ResetPacketsPool the instance to nil and empty the Packets interface in PacketsPool
func ResetPacketsPool() {
	instance = nil
}

// AddPackets will append the packets to slice
func (p *PacketsPool) AddPackets(i interface{}) {
	p.Packets = append(p.Packets, i)
}
