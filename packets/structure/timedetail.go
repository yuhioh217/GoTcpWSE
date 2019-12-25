package structure

import "errors"

// Packets the packet that with the immediate trading info
type Packets struct {
	ID         string
	Name       string
	Timestamp  string
	Dealprice  float64
	Ordercount float64
	TotalCount float64
	Result     chan string
}

// NewPackets : new a Packets object
func NewPackets(id string, timeStamp string) (*Packets, error) {
	if id == "" {
		return nil, errors.New("id is null")
	}

	if timeStamp == "" {
		return nil, errors.New("timeStamp is null")
	}

	return &Packets{
		ID:        id,
		Timestamp: timeStamp}, nil
}
