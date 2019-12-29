package structure

var (
	rttradingInstance *RealtimeTrading
)

// RealtimeTrading real time trading information
type RealtimeTrading struct {
	ID          string  //[ETX]0000[ETX]2337
	Timestamp   string  //[STX]0102[ETX]13:22:44 [STX]0602[ETX]13:22:44
	Type        int     //[STX]0119[ETX]1 [STX]0619[ETX]2// 1->green, 2->red
	Deal        float64 //[STX]0104[ETX]36.65, [STX]0604[ETX]36.65
	OrderCount  float64 //[STX]0114[ETX]2397 [STX]0614[ETX]2397
	TotalCount  float64 //[STX]0113[ETX]34033 [STX]0613[ETX]34033
	TotalAmount int     //[STX]0115[ETX]1255601900 [STX]0615[ETX]1255601900
}

// GetRealTimeTradingInstance to return the current real time trading info packets
func GetRealTimeTradingInstance() *RealtimeTrading {
	if rttradingInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		rttradingInstance = &RealtimeTrading{}
	}
	return rttradingInstance
}

// ResetTimeTradingInstance to reset the trading packets
func ResetTimeTradingInstance() *RealtimeTrading {
	rttradingInstance = nil
	return GetRealTimeTradingInstance()
}

// IsFinished check the data is filled or not
func (r *RealtimeTrading) IsFinished() bool {
	if r.ID != "" && r.Timestamp != "" && r.Type != 0 && r.Deal != 0 &&
		r.OrderCount != 0 && r.TotalAmount != 0 && r.TotalCount != 0 {
		return true
	}
	return false
}
