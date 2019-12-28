package api

import (
	"fmt"
	"reflect"
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

var sPool *structure.PacketsPool

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

func packageQueue(packStr string) { //, pq *s.PQueue) {
	// fmt.Println(packStr)
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
			sPool = structure.GetInstance()
			switch datatype {
			case "id":
				tempID = info
				break
			case "best_five_left":
				if pack == nil {
					p, _ := pack.(structure.RealtimeFive)
					p.ID = tempID
					p.Type = 0
					p.Five = info
					sPool.AddPackets(p)

					for _, pool := range sPool.Packets {
						fmt.Println(reflect.TypeOf(pool))
					}
					//fmt.Println(sPool.Packets)
					pack = nil
				}
				break
			default:
				break
			}

		}
	}

	/*
		for packet := range packetspool {
			fmt.Println(packet)
		}*/

}
