package main

import (
	"fmt"
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

func (sb *SBratio) AddTotal(t int) {
	sb.Total = sb.Total + t
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
	//timeArr := df.Col("Time").String()
	da := df.Col("Deal").Float()    // deal array
	ta, err := df.Col("Type").Int() // type array
	checkError(err)
	oca, err := df.Col("OrderCount").Int() // order count array
	checkError(err)
	//tca, err := df.Col("TotalCount").Int() // total count array
	//checkError(err)

	switch ta[len(ta)-1] { // use the latest packets type
	case 1: // sell deal
		sbRatio.AddSell(oca[len(oca)-1])
		sbRatio.AddTotal(oca[len(oca)-1])
		break
	case 2: // buy deal
		sbRatio.AddBuy(oca[len(oca)-1])
		sbRatio.AddTotal(oca[len(oca)-1])
		break
	default:
		break
	}
	fmt.Println(sbRatio)
	if bratio, sratio := sbRatio.GetBratio(), sbRatio.GetSratio(); bratio != 0 && sratio != 0 {
		if len(oca) > 5 { // do nothing before 9:00:00 - 9:00:30 > about 5 packets data
			//if (sratio / bratio) < 0.67 { // 40%/60% -> go to sell
			lastAmount := int(da[len(da)-1] * 1000 * float64(oca[len(oca)-1]))
			fmt.Println("lastAmount : ", lastAmount)
			// if current trading amount is over pre-ten trading average amount
			if lastAmount >= GetPreTenAverageAmount(da, oca)*2 {
				fmt.Println("Strong buy")
			}
			//}
		}
	}
}

func GetContinuousTenSlope() float64 {
	return 0.0
}

func GetPreTenAverageAmount(dealArr []float64, ocaArr []int) int {
	fmt.Print("Get the pre-ten deals average amount : ")
	sum := 0.0
	for i := len(dealArr) - 2; i > len(dealArr)-10 && i >= 0; i-- {
		// fmt.Printf("%.2f * %d \n", dealArr[i], ocaArr[i])
		sum += dealArr[i] * 1000 * float64(ocaArr[i])
	}
	fmt.Println(int(sum / float64(len(dealArr))))
	return int(sum / float64(len(dealArr)))
}

func main() {

	rows := []Info{}

	test := []Info{
		{"13:22:24", 36.75, 1, 2, 2},
		{"13:22:29", 36.75, 1, 5, 7},
		{"13:22:34", 36.80, 2, 9, 16},
		{"13:22:39", 36.75, 1, 18, 34},
		{"13:22:39", 36.80, 2, 2, 34},
		{"13:22:39", 36.75, 1, 11, 34},
		{"13:22:39", 36.80, 2, 40, 34},
		{"13:22:39", 36.80, 2, 70, 34},
		{"13:22:39", 36.85, 2, 2, 34},
		{"13:22:39", 36.80, 1, 12, 34},
	}

	for _, t := range test {
		rows = append(rows, t)
		NewModel2Dataframe(rows)
	}
}
