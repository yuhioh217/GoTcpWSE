package structure

// RealtimeFive to record the data that the best buy or sell five
type RealtimeFive struct {
	ID   string
	Type int // 0 -> sell, 1 -> buy
	Five string
}

// SetID to set the packets id
func (r *RealtimeFive) SetID(id string) {
	r.ID = id
}

// SetType to set the packets type
func (r *RealtimeFive) SetType(Type int) {
	r.Type = Type
}

// SetFive to set the packets five info
func (r *RealtimeFive) SetFive(Five string) {
	r.Five = Five
}

// GetDataFinished to check the data is ready to pool
func (r *RealtimeFive) GetDataFinished() bool {
	if r.ID != "" && r.Five != "" {
		return true
	}
	return false
}

var resetData = &RealtimeFive{}

// Reset the struct
func (r *RealtimeFive) Reset() {
	*r = *resetData
}
