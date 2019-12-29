package structure

var (
	dealinstance *CurrentDeal
)

// CurrentDeal to record the current deal point information
type CurrentDeal struct {
	Deal float64 // current deal point
	Type int     // 1: green, 2: red
}

func getCurrentDealInstance() *CurrentDeal {
	if dealinstance == nil {
		lock.Lock()
		defer lock.Unlock()
		dealinstance = &CurrentDeal{}
	}
	return dealinstance
}
