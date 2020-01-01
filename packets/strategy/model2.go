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
	Sell   int
	Buy    int
	Total  int
	Sratio float64
	Bratio float64
}

func GetSBratioInstance() *SBratio {
	if sbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		sbInstance = &SBratio{0, 0, 0, 0.0, 0.0}
	}
	return sbInstance
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

func Model2Calculate() {
	for _, d := range df.Col("Deal").Float() {
		fmt.Println(d)
	}
}

func main() {
	rows := []Info{
		{"13:22:24", 36.75, 1, 2, 33337},
		{"13:22:29", 36.75, 1, 5, 33342},
		{"13:22:34", 36.80, 1, 9, 33351},
		{"13:22:39", 36.75, 2, 18, 33369},
	}

	NewModel2Dataframe(rows)
}
