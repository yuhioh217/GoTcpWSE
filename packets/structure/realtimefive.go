package structure

// RealtimeFive to record the data that the best buy or sell five
type RealtimeFive struct {
	ID   string
	Type int // 0 -> sell, 1 -> buy
	Five string
}

func (r *RealtimeFive) setID(id string) {
	r.ID = id
}

func (r *RealtimeFive) setType(Type int) {
	r.Type = Type
}

func (r *RealtimeFive) setFive(Five string) {
	r.Five = Five
}

func (r *RealtimeFive) getDataFinished() bool {
	if r.ID != "" && r.Type != 0 && r.Five != "" {
		return true
	}
	return false
}

var resetData = &RealtimeFive{}

// Reset the struct
func (r *RealtimeFive) Reset() {
	*r = *resetData
}
