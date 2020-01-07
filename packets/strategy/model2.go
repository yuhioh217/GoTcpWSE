package strategy

import (
	"fmt"
	"math"
	"sync"

	"../structure"
	"github.com/go-gota/gota/dataframe"
)

// MPT(Modern Protfolio Theory) strategy

var (
	lock       *sync.Mutex = &sync.Mutex{}
	sbInstance *SBratio
	df         dataframe.DataFrame
	slopeArr   []SlopeM
	i          int
)

// SlopeM struct to No./amount coordinate data save
type SlopeM struct {
	Direction int
	Result    float64
}

// SBratio struct to save the sell/buy ratio data
type SBratio struct {
	Sell       int
	Buy        int
	Total      int
	Accumulate []int // Current accumulate buy amount
	//Sratio float64
	//Bratio float64
}

// GetSBratioInstance : create or get the current SBratio instance (only the single instance)
func GetSBratioInstance() *SBratio {
	if sbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		sbInstance = &SBratio{0, 0, 0, []int{}} //, 0.0, 0.0}
	}
	return sbInstance
}

// AddSell : add sell count
func (sb *SBratio) AddSell(s int) {
	sb.Sell = sb.Sell + s
}

// AddBuy : add buy count
func (sb *SBratio) AddBuy(b int) {
	sb.Buy = sb.Buy + b
}

// AddTotal to add the total trading count
func (sb *SBratio) AddTotal(t int) {
	sb.Total = sb.Total + t
}

// UpdateAccumulate to update the current accumulative amount
func (sb *SBratio) UpdateAccumulate(t int, a int) {
	switch t {
	case 1:
		if len(sb.Accumulate) == 0 {
			sb.Accumulate = append(sb.Accumulate, 0-a)
		} else {
			sb.Accumulate = append(sb.Accumulate, sb.Accumulate[len(sb.Accumulate)-1]-a)
		}
	case 2:
		if len(sb.Accumulate) == 0 {
			sb.Accumulate = append(sb.Accumulate, a)
		} else {
			sb.Accumulate = append(sb.Accumulate, sb.Accumulate[len(sb.Accumulate)-1]+a)
		}
	default:
		break
	}
}

// GetSratio get the Sell - Total ratio
func (sb *SBratio) GetSratio() float64 {
	return float64(sb.Sell) / float64(sb.Total)
}

// GetBratio get the Buy - Total ratio
func (sb *SBratio) GetBratio() float64 {
	return float64(sb.Buy) / float64(sb.Total)
}

// Coordinate to save the each packet amount data to (x, y) -> (No. , amount)
type Coordinate struct {
	PacketNo int // coordinate x
	Amount   int // coordinate y
}

// Info struct to save the each deal info
type Info struct {
	Time       string
	Deal       float64
	Type       int
	OrderCount int
	TotalCount int
}

// NewModel2Dataframe to new the packets dataframe
func NewModel2Dataframe(rows []Info) {
	df = dataframe.LoadStructs(rows)
	Model2Calculate()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Model2Calculate is the entry to run the model2
func Model2Calculate() {
	sbRatio := GetSBratioInstance()
	timeArr := df.Col("Time").Records()
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
		if len(oca) > 11 { // do nothing before 9:00:00 - 9:00:30 > about 5 packets data
			//if (sratio / bratio) < 0.67 { // 40%/60% -> go to sell
			lastAmount := int(da[len(da)-1] * 1000 * float64(oca[len(oca)-1]))
			d, result := GetContinuousTenSlope(da, oca)
			slopeArr = append(slopeArr, SlopeM{d, result})
			fmt.Println("LastAmount :", lastAmount, ", Pre-10 average :", GetPreTenAverageAmount(da, oca))

			// if current trading amount is over pre-ten trading average amount
			fmt.Println((sratio / bratio))
			fmt.Println((bratio / sratio))
			if lastAmount >= int(float64(GetPreTenAverageAmount(da, oca))*6) &&
				//ta[len(ta)-1] == slopeArr[len(slopeArr)-1].Direction &&
				slopeArr[len(slopeArr)-1].Result/slopeArr[len(slopeArr)-2].Result > 1.3 {

				if slopeArr[len(slopeArr)-1].Direction == 2 && (sratio/bratio) < 1 {
					each := structure.EachHold{Deal: da[len(da)-1], Type: 2, Count: 1}
					structure.GetCurrentDealInstance().AddHolding(each)
					fmt.Printf("\033[0;34mCurrent Packets Deal Amount is bigger than the pre-ten deals average amount. (Strong buy)\033[0m \n")
				} else if slopeArr[len(slopeArr)-1].Direction == 1 && (bratio/sratio) < 1 {
					each := structure.EachHold{Deal: da[len(da)-1], Type: 1, Count: 1}
					structure.GetCurrentDealInstance().AddHolding(each)
					fmt.Printf("\033[0;34mCurrent Packets Deal Amount is bigger than the pre-ten deals average amount. (Strong Sell)\033[0m \n")
				}

			}
			//}
		}
	}
	if len(structure.GetCurrentDealInstance().Keep) != 0 && timeArr[len(timeArr)-1] == "13:30:00" {
		structure.GetCurrentDealInstance().ReleaseAllHolding(da[len(da)-1])
	}
}

// GetContinuousTenSlope to get the continuous ten (first and tenth point's slope(M))
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

// GetPreTenAverageAmount to get the pre ten deal average amount
func GetPreTenAverageAmount(dealArr []float64, ocaArr []int) int {
	// fmt.Print("Get the pre-ten deals average amount : ")
	sum := 0.0
	for i := len(dealArr) - 2; i > len(dealArr)-10 && i >= 0; i-- {
		// fmt.Printf("%.2f * %d \n", dealArr[i], ocaArr[i])
		sum += dealArr[i] * 1000 * float64(ocaArr[i])
	}
	//fmt.Println(int(sum / float64(len(dealArr))))
	return int(sum / 10)
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
		{"13:22:39", 36.85, 1, 31, 34},
		{"13:22:39", 36.80, 1, 9, 34},
		{"13:22:39", 36.80, 1, 5, 34},
	}

	for _, t := range test {
		rows = append(rows, t)
		// ordering the packet info
		NewModel2Dataframe(rows)
	}
}
