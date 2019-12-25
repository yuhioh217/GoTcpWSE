package structure

import "errors"

// Packets the packet that with the immediate trading info
type Packets struct {
	ID         string
	Name       string
	Timestamp  string
	Action     float64 // 2: Buy(red), 1: Sell(green)
	Dealprice  float64
	Ordercount float64
	TotalCount float64
	Result     string
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

func (p *Packets) setID(id string) {
	p.ID = id
}

func (p *Packets) setName(name string) {
	p.Name = name
}

func (p *Packets) setTimestamp(time string) {
	if p.Timestamp != time {
		p.Result = "package timeout"
	}
	p.Timestamp = time
}

func (p *Packets) setDealprice(dealprice float64) {
	p.Dealprice = dealprice
}

func (p *Packets) setOrdercount(ordercount float64) {
	p.Ordercount = ordercount
}

func (p *Packets) setTotalCount(totalcount float64) {
	p.TotalCount = totalcount
}

func (p *Packets) isReady() (bool, string) {
	if p.ID != "" && p.Name != "" && p.Timestamp != "" && p.Dealprice != 0 && p.Ordercount != 0 && p.TotalCount != 0 {
		return true, "Package collect finished"
	}
	if p.Result != "" {
		return true, p.Result
	}

	return false, "Package collect not finished"
}
