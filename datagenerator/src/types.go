package src

import "time"

type GenerationInfoType struct {
	strTickerName string
	dtTickerDate  time.Time
}

type FilesInfoType struct {
	GenerationInfo   GenerationInfoType
	strBuyPath       string
	strSellPath      string
	strBenchmarkPath string
}
