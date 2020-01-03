package main

import (
	"fmt"
	"math"
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
	Sell       int
	Buy        int
	Total      int
	Accumulate []int // Current accumulate buy amount
	//Sratio float64
	//Bratio float64
}

func GetSBratioInstance() *SBratio {
	if sbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		sbInstance = &SBratio{0, 0, 0, []int{}} //, 0.0, 0.0}
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

func (sb *SBratio) UpdateAccumulate(t int, a int) {
	switch t {
	case 1:
		if len(sb.Accumulate) == 0 {
			sb.Accumulate = append(sb.Accumulate, 0-a)
		} else {
			sb.Accumulate = append(sb.Accumulate, sb.Accumulate[len(sb.Accumulate)-1]-a)
		}
		break
	case 2:
		if len(sb.Accumulate) == 0 {
			sb.Accumulate = append(sb.Accumulate, a)
		} else {
			sb.Accumulate = append(sb.Accumulate, sb.Accumulate[len(sb.Accumulate)-1]+a)
		}
		break
	default:
		break
	}

}

func (sb *SBratio) GetSratio() float64 {
	return float64(sb.Sell) / float64(sb.Total)
}

func (sb *SBratio) GetBratio() float64 {
	return float64(sb.Buy) / float64(sb.Total)
}

type Coordinate struct {
	PacketNo int // coordinate x
	Amount   int // coordinate y
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
		sbRatio.UpdateAccumulate(1, int(da[len(oca)-1]*1000*float64(oca[len(oca)-1])))
		break
	case 2: // buy deal
		sbRatio.AddBuy(oca[len(oca)-1])
		sbRatio.AddTotal(oca[len(oca)-1])
		sbRatio.UpdateAccumulate(2, int(da[len(oca)-1]*1000*float64(oca[len(oca)-1])))
		break
	default:
		break
	}

	//fmt.Println(sbRatio)
	if bratio, sratio := sbRatio.GetBratio(), sbRatio.GetSratio(); bratio != 0 && sratio != 0 {
		if len(oca) > 9 { // do nothing before 9:00:00 - 9:00:30 > about 5 packets data
			//if (sratio / bratio) < 0.67 { // 40%/60% -> go to sell
			GetContinuousTenSlope(da, oca)
			/*lastAmount := int(da[len(da)-1] * 1000 * float64(oca[len(oca)-1]))
			fmt.Println("lastAmount : ", lastAmount)
			// if current trading amount is over pre-ten trading average amount
			if lastAmount >= GetPreTenAverageAmount(da, oca)*2 {
				fmt.Println("Strong buy")
			}*/
			//}
		}
	}
}

func GetContinuousTenSlope(dealArr []float64, ocaArr []int) (int, float64) {
	tempCoordinate := []Coordinate{}
	direction := 0 // 1 -> down , 2 -> up
	// fmt.Println(len(dealArr))
	// Configure the amount and packets No. to coordinate struct
	for i := len(dealArr) - 10; i < len(dealArr) && i >= 0; i++ {
		// fmt.Printf("%.2f * %d \n", dealArr[i], ocaArr[i])
		temp := Coordinate{i, GetSBratioInstance().Accumulate[i]}
		tempCoordinate = append(tempCoordinate, temp)
		//fmt.Println(tempCoordinate)
	}
	c1 := tempCoordinate[len(tempCoordinate)-10]
	c2 := tempCoordinate[len(tempCoordinate)-1]

	fmt.Printf("Coordinate : (%d, %d)-(%d, %d) \n", c1.PacketNo, c1.Amount, c2.PacketNo, c2.Amount)
	//fmt.Println(c2.PacketNo - c1.PacketNo)
	divideResult := float64(c2.Amount-c1.Amount) / float64(c2.PacketNo-c1.PacketNo)
	if divideResult >= 0 {
		direction = 2
		fmt.Printf("\033[0;31mSlope(M) : %.2f ↗ \033[0m \n", math.Abs(divideResult))
	} else {
		direction = 1
		fmt.Printf("\033[0;32mSlope(M) : %.2f ↘	\033[0m \n", math.Abs(divideResult))
	}
	return direction, math.Abs(divideResult)
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
		{"13:22:39", 36.85, 2, 11, 34},
		{"13:22:39", 36.80, 1, 50, 34},
		{"13:22:39", 36.85, 2, 12, 34},
		{"13:22:39", 36.80, 1, 3, 34},
		{"13:22:39", 36.80, 1, 1, 34},
		{"13:22:39", 36.80, 1, 20, 34},
		{"13:22:39", 36.85, 2, 1, 34},
		{"13:22:39", 36.85, 2, 3, 34},
		{"13:22:39", 36.85, 2, 31, 34},
		{"13:22:39", 36.80, 1, 9, 34},
		{"13:22:39", 36.80, 1, 5, 34},
	}

	for _, t := range test {
		rows = append(rows, t)
		NewModel2Dataframe(rows)
	}
}
