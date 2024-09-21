package src

import (
	"fmt"
	"time"
)

func findTradeFiles(a_TradeRunInfo TradeRunInfoType) bool {
	const (
		c_strMethodName = "reader.findTradeFiles"
	)
	var (
		strInputPath     string
		strFullPath      string
		bTradeFileExists bool
		bBuyFileExists   bool
		bSellFileExists  bool
	)
	strInputPath = getInputPath() + "/"

	strFullPath = strInputPath + fmt.Sprintf(c_strTradeFile, a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.strTickerName)
	bTradeFileExists = checkFileExists(strFullPath)

	strFullPath = strInputPath + fmt.Sprintf(c_strBuyFile, a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.strTickerName)
	bBuyFileExists = checkFileExists(strFullPath)

	strFullPath = strInputPath + fmt.Sprintf(c_strSellFile, a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.strTickerName)
	bSellFileExists = checkFileExists(strFullPath)

	if bTradeFileExists {
		return bBuyFileExists && bSellFileExists
	}

	return bBuyFileExists || bBuyFileExists
}

func getValidTickers(a_dtTickerDate time.Time) []string {

}

func loadTradeDataFromFile(a_TradeRunInfo TradeRunInfoType) {

}
