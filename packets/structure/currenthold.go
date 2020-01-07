package structure

import (
	"fmt"
	"math"
	"sort"
)

var (
	holdinstance *CurrentHold
)

// EachHold to save the realtime deal info
type EachHold struct {
	Deal  float64
	Type  int
	Count int
}

// CurrentHold to record the current hold information
type CurrentHold struct {
	Keep    []EachHold
	Surplus float64
}

// GetCurrentDealInstance to get the current deal instance
func GetCurrentDealInstance() *CurrentHold {
	if holdinstance == nil {
		lock.Lock()
		defer lock.Unlock()
		holdinstance = &CurrentHold{}
	}
	return holdinstance
}

func checkAddValid(each EachHold) bool {
	for _, k := range GetCurrentDealInstance().Keep {
		if each.Type != k.Type {
			return false
		}
	}
	return true
}

// SortingHold to sort the hold deal info
func (c *CurrentHold) SortingHold() {
	sort.Slice(c.Keep, func(i, j int) bool {
		return c.Keep[i].Deal < c.Keep[j].Deal
	})

}

// AddHolding to add the trading to current hold portion
func (c *CurrentHold) AddHolding(each EachHold) {
	if checkAddValid(each) {
		c.Keep = append(c.Keep, each)
		c.SortingHold()
	} else {
		// Type is not the same, do the trading analysis
		c.TradingDeal(each)
	}

}

// GetCurrentHoldAverageDeal to get the current hold average amount
func (c *CurrentHold) GetCurrentHoldAverageDeal() (int, float64) {
	sum := 0.0
	count := 0
	if len(c.Keep) > 0 {
		for _, v := range c.Keep {
			sum += v.Deal * float64(v.Count)
			count += v.Count
		}
		return count, sum / float64(len(c.Keep))
	}
	return 0, 0.0
}

// TradingDeal to trade from current hold status
func (c *CurrentHold) TradingDeal(each EachHold) {
	switch each.Type {
	case 1: // Will do the sell action
		count, holdAverage := c.GetCurrentHoldAverageDeal()
		if (float64(each.Deal) - math.Sqrt(holdAverage)/10) >= holdAverage {
			c.Surplus += (each.Deal - holdAverage) * float64(count*1000)
			c.Keep = nil
			// Sell all
		} else if (float64(each.Deal) - math.Sqrt(holdAverage)/15) >= holdAverage {
			tempCount := 0
			for i, v := range c.Keep {
				if float64(tempCount)/float64(count) < 0.5 {
					tempCount += v.Count
					c.Surplus += (each.Deal - v.Deal) * float64(v.Count*1000)
					c.Keep = append(c.Keep[:i], c.Keep[i+1:]...)
				} else {
					break
				}
			}
			// Sell 50% portion
		} else if (float64(each.Deal) - math.Sqrt(holdAverage)/20) >= holdAverage {
			c.Surplus += (each.Deal - c.Keep[0].Deal) * float64(1*1000)
			c.Keep = append(c.Keep[:0], c.Keep[1:]...)
			// Sell 1
		} else if (math.Sqrt(holdAverage)/20 + holdAverage) >= float64(each.Deal) {
			// Sell 1
			c.Surplus += (each.Deal - c.Keep[len(c.Keep)-1].Deal) * float64(1*1000)
			c.Keep = append(c.Keep[:len(c.Keep)-1], c.Keep[len(c.Keep)-1:]...)
		} else if (math.Sqrt(holdAverage)/15 + holdAverage) >= float64(each.Deal) {
			// Sell 50%
			tempCount := 0
			for i := len(c.Keep) - 1; i >= 0; i-- {
				if float64(tempCount)/float64(count) < 0.5 {
					tempCount += c.Keep[i].Count
					c.Surplus += (each.Deal - c.Keep[i].Deal) * float64(c.Keep[i].Count*1000)
					c.Keep = append(c.Keep[:i], c.Keep[i+1:]...)
				} else {
					break
				}
			}
		} else if (math.Sqrt(holdAverage)/10 + holdAverage) >= float64(each.Deal) {
			// Sell all
			c.Surplus += (each.Deal - holdAverage) * float64(count*1000)
			c.Keep = nil
		}

	case 2: // Will do the buy action
		count, holdAverage := c.GetCurrentHoldAverageDeal()
		if (holdAverage - math.Sqrt(holdAverage)/10) <= float64(each.Deal) {
			// return to buy all
			c.Surplus += (holdAverage - each.Deal) * float64(count*1000)
			c.Keep = nil
		} else if (holdAverage - math.Sqrt(holdAverage)/15) <= float64(each.Deal) {
			// return to buy 50%
			tempCount := 0
			for i := len(c.Keep) - 1; i >= 0; i-- {
				if float64(tempCount)/float64(count) < 0.5 {
					tempCount += c.Keep[i].Count
					c.Surplus += (c.Keep[i].Deal - each.Deal) * float64(c.Keep[i].Count*1000)
					c.Keep = append(c.Keep[:i], c.Keep[i+1:]...)
				} else {
					break
				}
			}
		} else if (holdAverage - math.Sqrt(holdAverage)/20) <= float64(each.Deal) {
			// return to buy 1
			c.Surplus += (c.Keep[len(c.Keep)-1].Deal - each.Deal) * float64(1*1000)
			c.Keep = append(c.Keep[:len(c.Keep)-1], c.Keep[len(c.Keep)-1:]...)
		} else if (holdAverage + math.Sqrt(holdAverage)/20) >= float64(each.Deal) {
			// return to buy 1
			c.Surplus += (c.Keep[0].Deal - each.Deal) * float64(1*1000)
			c.Keep = append(c.Keep[:0], c.Keep[1:]...)
		} else if (holdAverage + math.Sqrt(holdAverage)/15) >= float64(each.Deal) {
			// return to buy 50%
			tempCount := 0
			for i, v := range c.Keep {
				if float64(tempCount)/float64(count) < 0.5 {
					tempCount += v.Count
					c.Surplus += (v.Deal - each.Deal) * float64(v.Count*1000)
					c.Keep = append(c.Keep[:i], c.Keep[i+1:]...)
				} else {
					break
				}
			}
		} else if (holdAverage + math.Sqrt(holdAverage)/10) >= float64(each.Deal) {
			// return to buy all
			c.Surplus += (holdAverage - each.Deal) * float64(count*1000)
			c.Keep = nil
		}
	default:
		break
	}
}

// ReleaseAllHolding : when in the last trading, release all holdings
func (c *CurrentHold) ReleaseAllHolding(currentDeal float64) {
	count, holdAverage := c.GetCurrentHoldAverageDeal()
	fmt.Println("Final Deal", currentDeal)
	fmt.Println("Before releasing holdings :", c.Keep)
	fmt.Println("Average :", holdAverage)
	switch c.Keep[0].Type {
	case 1:
		c.Surplus += (holdAverage - currentDeal) * float64(count*1000)
		c.Keep = nil
	case 2:
		c.Surplus += (currentDeal - holdAverage) * float64(count*1000)
		c.Keep = nil
	default:
		break
	}
}
