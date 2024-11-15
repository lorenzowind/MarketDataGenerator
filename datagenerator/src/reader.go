package src

import (
	"errors"
	"fmt"
	logger "marketmanipulationdetector/logger/src"
)

func getReferenceOffersBook(a_GenerationInfo GenerationInfoType) (FilesInfoType, error) {
	const (
		c_strMethodName = "reader.getReferenceOffersBook"
	)
	var (
		err              error
		strReferencePath string
		strBuyPath       string
		strSellPath      string
		strBenchmarkPath string
		bFileExists      bool
		FilesInfo        FilesInfoType
	)
	strReferencePath = getReferencePath() + "/"

	strBuyPath = strReferencePath + fmt.Sprintf(c_strBuyFile, a_GenerationInfo.dtTickerDate.Day(), a_GenerationInfo.dtTickerDate.Month(), a_GenerationInfo.dtTickerDate.Year(), a_GenerationInfo.strTickerName)
	bFileExists = checkFileExists(strBuyPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Buy file found : strBuyPath="+strBuyPath)
	} else {
		strBuyPath = ""
	}

	strSellPath = strReferencePath + fmt.Sprintf(c_strSellFile, a_GenerationInfo.dtTickerDate.Day(), a_GenerationInfo.dtTickerDate.Month(), a_GenerationInfo.dtTickerDate.Year(), a_GenerationInfo.strTickerName)
	bFileExists = checkFileExists(strSellPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Sell file found : strSellPath="+strSellPath)
	} else {
		strSellPath = ""
	}

	strBenchmarkPath = strReferencePath + c_strBenchmarksFile
	bFileExists = checkFileExists(strBenchmarkPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Benchmarks file found : strBenchmarkPath="+strBenchmarkPath)
	} else {
		strBenchmarkPath = ""
	}

	// Existe os 2 arquivos (compra e venda) ou existe pelo menos o de compra ou venda, alem do arquivo de benchmarks
	if strBuyPath != "" || strSellPath != "" {
		FilesInfo = FilesInfoType{
			GenerationInfo:   a_GenerationInfo,
			strBuyPath:       strBuyPath,
			strSellPath:      strSellPath,
			strBenchmarkPath: strBenchmarkPath,
		}
	} else {
		err = errors.New("file rules existance failed")
	}

	return FilesInfo, err
}
