package api

import (
	"strconv"
	"strings"

	// to have
	"../structure"
)

type tempBuffer struct {
	ID   string
	Five string
	Type int
}

func (t *tempBuffer) setFive(Five string) {
	t.Five = Five
}

func (t *tempBuffer) setType(Type int) {
	t.Type = Type
}

// ASCIIDecode to decode the STX adn ETX in packets
func ASCIIDecode(ascii []uint8) interface{} {
	str := ""
	//fmt.Println(ascii)
	for _, v := range ascii {
		if v == 0x02 {
			str += "[STX]"
		} else if v == 0x03 {
			str += "[ETX]"
		} else {
			str += string(v)
		}
	}
	return PacketsFilter(str)
}

// PacketsFilter to filter the packets that have the STX and ETX
func PacketsFilter(packStr string) interface{} {
	if strings.Contains(packStr, "[STX]") {
		packageQueue(packStr)
		return packStr
	}
	return nil
}

/** Pool Process **/

// GetCurrentPool return the current sPool
func GetCurrentPool() *structure.PacketsPool {
	return structure.GetInstance()
}

// ResetsPool to reset the spool data to empty
func ResetsPool() *structure.PacketsPool {
	structure.ResetPacketsPool()
	return structure.GetInstance()
}

// GetRealtimeTradingInstance to get the current realtime trading instance
func GetRealtimeTradingInstance() *structure.RealtimeTrading {
	return structure.GetRealTimeTradingInstance()
}

// ResetRealtimeTradingInstance to reset realtime trading instance
func ResetRealtimeTradingInstance() *structure.RealtimeTrading {
	structure.ResetTimeTradingInstance()
	return structure.GetRealTimeTradingInstance()
}

var sPool *structure.PacketsPool
var currentType int
var currentDeal float64

func packageQueue(packStr string) { //, pq *s.PQueue) {
	//fmt.Printf(PacketsColor+" :%s\n", "Packets", packStr)
	// process single package
	var pack interface{}
	var tempID string
	if tempID != "" {

	}
	for _, sub := range strings.Split(packStr, "[STX]") {
		//fmt.Println(sub)
		// best five data to struct

		if datatype, info := StringRule(sub); datatype != "" && info != "" {
			//fmt.Printf(FiveColor+": %s\n", datatype, info)
			sPool = GetCurrentPool()
			switch datatype {
			case "id":
				tempID = info
				break
			case "best_five_left":
				if pack == nil {
					p, _ := pack.(structure.RealtimeFive)
					p.SetID(tempID)
					p.SetType(0)
					p.SetFive(info)
					if p.GetDataFinished() {
						sPool.AddPackets(p)
					}

					pack = nil
				}
				break
			case "best_five_right":
				if pack == nil {
					p, _ := pack.(structure.RealtimeFive)
					p.SetID(tempID)
					p.SetType(1)
					p.SetFive(info)
					if p.GetDataFinished() {
						sPool.AddPackets(p)
					}
					pack = nil
				}
				break
			// realtime deal packets
			case "timestamp": // string
				p := GetRealtimeTradingInstance()
				if p.Timestamp != "" {
					// check the packets timestamp
					if p.Timestamp != info {
						p = ResetRealtimeTradingInstance()
					}
				}
				p.Timestamp = info
				p.ID = tempID
				break
			case "type": //int
				p := GetRealtimeTradingInstance()
				i, _ := strconv.Atoi(info)
				currentType = int(i)
				p.Type = currentType
				break
			case "deal_price": // float64
				p := GetRealtimeTradingInstance()
				f, _ := strconv.ParseFloat(info, 64)
				currentDeal = f
				p.Deal = currentDeal
				break
			case "order_count": // float64
				p := GetRealtimeTradingInstance()
				f, _ := strconv.ParseFloat(info, 64)
				p.OrderCount = f
				break
			case "total_count": // float64
				p := GetRealtimeTradingInstance()
				f, _ := strconv.ParseFloat(info, 64)
				p.TotalCount = f
				break
			case "total_amount": // float64
				p := GetRealtimeTradingInstance()
				i, _ := strconv.Atoi(info)
				p.TotalAmount = i
				break
			default:
				break
			}

			if grtInstance := GetRealtimeTradingInstance(); grtInstance.Deal == 0 {
				grtInstance.Type = currentType
			}

			if grtInstance := GetRealtimeTradingInstance(); grtInstance.Deal == 0 {
				grtInstance.Deal = currentDeal
			}

			if GetRealtimeTradingInstance().IsFinished() {
				sPool.AddPackets(*GetRealtimeTradingInstance())
				ResetRealtimeTradingInstance()
			}
		}
	}
}
