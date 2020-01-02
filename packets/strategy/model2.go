package main

import (
	"sync"

	"github.com/go-gota/gota/dataframe"
)

// MPT(Modern Protfolio Theory) strategy

var (
	lock       *sync.Mutex = &sync.Mutex{}
	sbInstance *SBratio
	df         dataframe.DataFrame
)

type SBratio struct {
	Sell  int
	Buy   int
	Total int
	//Sratio float64
	//Bratio float64
}

func GetSBratioInstance() *SBratio {
	if sbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		sbInstance = &SBratio{0, 0, 0} //, 0.0, 0.0}
	}
	return sbInstance
}

func (sb *SBratio) AddSell(s int) {
	sb.Sell = sb.Sell + s
}

func (sb *SBratio) AddBuy(b int) {
	sb.Buy = sb.Buy + b
}

func (sb *SBratio) UpdateTotal(t int) {
	sb.Total = t
}

func (sb *SBratio) GetSratio() float64 {
	return float64(sb.Sell) / float64(sb.Total)
}

func (sb *SBratio) GetBratio() float64 {
	return float64(sb.Buy) / float64(sb.Total)
}

type Info struct {
	Time       string
	Deal       float64
	Type       int
	OrderCount int
	TotalCount int
}

func NewModel2Dataframe(rows []Info) {
	df = dataframe.LoadStructs(rows)
	Model2Calculate()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Model2Calculate() {
	sbRatio := GetSBratioInstance()
	timeArr = df.Col("Time").String()
	//da := df.Col("Deal").Float()    // deal array
	ta, err := df.Col("Type").Int() // type array
	checkError(err)
	oca, err := df.Col("OrderCount").Int() // order count array
	checkError(err)
	tca, err := df.Col("TotalCount").Int() // total count array
	checkError(err)

	switch ta[len(ta)-1] { // use the latest packets type
	case 1: // sell deal
		sbRatio.AddSell(oca[len(oca)-1])
		sbRatio.UpdateTotal(tca[len(tca)-1])
		break
	case 2: // buy deal
		sbRatio.AddBuy(oca[len(oca)-1])
		sbRatio.UpdateTotal(tca[len(tca)-1])
		break
	default:
		break
	}
}

func main() {
	rows := []Info{
		{"13:22:24", 36.75, 1, 2, 2},
		{"13:22:29", 36.75, 1, 5, 7},
		{"13:22:34", 36.80, 1, 9, 16},
		{"13:22:39", 36.75, 2, 18, 34},
	}

	NewModel2Dataframe(rows)
}
