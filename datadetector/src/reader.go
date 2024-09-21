package src

import (
	"container/list"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	logger "marketmanipulationdetector/logger/src"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getUniqueTickerFiles(a_TradeRunInfo TradeRunInfoType) (FilesInfoType, error) {
	const (
		c_strMethodName = "reader.getUniqueTickerFiles"
	)
	var (
		err          error
		strInputPath string
		strTradePath string
		strBuyPath   string
		strSellPath  string
		bFileExists  bool
		FilesInfo    FilesInfoType
	)
	strInputPath = getInputPath() + "/"

	strTradePath = strInputPath + fmt.Sprintf(c_strTradeFile, a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.strTickerName)
	bFileExists = checkFileExists(strTradePath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Trade file found : strTradePath="+strTradePath)
	} else {
		strTradePath = ""
	}

	strBuyPath = strInputPath + fmt.Sprintf(c_strBuyFile, a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.strTickerName)
	bFileExists = checkFileExists(strBuyPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Buy file found : strBuyPath="+strBuyPath)
	} else {
		strBuyPath = ""
	}

	strSellPath = strInputPath + fmt.Sprintf(c_strSellFile, a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.strTickerName)
	bFileExists = checkFileExists(strSellPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Sell file found : strSellPath="+strSellPath)
	} else {
		strSellPath = ""
	}

	// Existe os 3 arquivos (compra, venda e negocio) ou existe pelo menos o de compra ou venda
	if (strTradePath != "" && strBuyPath != "" && strSellPath != "") || (strTradePath == "" && (strBuyPath != "" || strSellPath != "")) {
		FilesInfo = FilesInfoType{
			TradeRunInfo: a_TradeRunInfo,
			strTradePath: strTradePath,
			strBuyPath:   strBuyPath,
			strSellPath:  strSellPath,
		}
	} else {
		err = errors.New("file rules existance failed")
	}

	return FilesInfo, err
}

func getAllTickersFiles(a_dtTickerDate time.Time) []FilesInfoType {
	const (
		c_strMethodName = "reader.getAllTickersFiles"
	)
	var (
		err            error
		strFileName    string
		strInputPath   string
		arrFileInfo    []string
		FilesInfo      FilesInfoType
		arrTickersInfo []FilesInfoType
		dtTickerDate   time.Time
		dir            fs.DirEntry
		arrDir         []fs.DirEntry
		TradeRunInfo   TradeRunInfoType
	)
	strInputPath = getInputPath() + "/"

	arrDir, err = os.ReadDir(strInputPath)

	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to get the directory : "+err.Error())
		return arrTickersInfo
	}

	// Itera sobre cada arquivo da pasta input
	for _, dir = range arrDir {
		// So verifica se for um arquivo
		if !dir.IsDir() {
			strFileName = filepath.Base(dir.Name())
			arrFileInfo = strings.Split(strFileName, "_")

			// So verifica arquivo no formato yyyy-mm-dd_<TICKER>_<TRADE/BUY/SELL>.csv
			if len(arrFileInfo) == 3 {
				dtTickerDate, err = validateDateString(arrFileInfo[0])
				if err == nil {
					// So verifica se tem a mesma data do que foi passado via parametro e nao foi adicionado ainda
					if checkIfHasSameDate(dtTickerDate, a_dtTickerDate) && !checkIfContains(arrFileInfo[1], arrTickersInfo) {
						// So verifica se segue as regras de existencia dos arquivos (compra, venda e negocio)
						TradeRunInfo = TradeRunInfoType{
							strTickerName: arrFileInfo[1],
							dtTickerDate:  dtTickerDate,
						}
						FilesInfo, err = getUniqueTickerFiles(TradeRunInfo)
						if err == nil {
							arrTickersInfo = append(arrTickersInfo, FilesInfo)
						}
					}
				}
			}
		}
	}

	return arrTickersInfo
}

func loadTickerData(a_FilesInfo FilesInfoType) TickerDataType {
	var (
		TickerData TickerDataType
	)
	// Carrega dados do arquivo de trade
	if a_FilesInfo.strTradePath != "" {
		TickerData.lstTrade = loadTradeDataFromFile(a_FilesInfo.strTradePath, a_FilesInfo.TradeRunInfo.strTickerName)
	}

	// Carrega dados do arquivo de compra
	if a_FilesInfo.strBuyPath != "" {
		TickerData.lstBuy = loadOfferDataFromFile(a_FilesInfo.strBuyPath)
	}

	// Carrega dados do arquivo de venda
	if a_FilesInfo.strSellPath != "" {
		TickerData.lstSell = loadOfferDataFromFile(a_FilesInfo.strSellPath)
	}

	TickerData.FilesInfo = a_FilesInfo
	return TickerData
}

func loadTradeDataFromFile(a_strPath, a_strTicker string) list.List {
	const (
		c_strMethodName          = "reader.loadTradeDataFromFile"
		c_nOperationIndex        = 0
		c_nTickerIndex           = 1
		c_nTimeIndex             = 2
		c_nOfferGenerationIndex  = 3
		c_nAccountIndex          = 4
		c_nIDIndex               = 5
		c_nOfferPrimaryIDIndex   = 6
		c_nOfferSecondaryIDIndex = 7
		c_nQuantityIndex         = 8
		c_nPriceIndex            = 9
		c_nLastIndex             = 9
	)
	var (
		err            error
		lstData        list.List
		arrRecord      []string
		arrFullRecords [][]string
		file           *os.File
		reader         *csv.Reader
		TradeData      TradeDataType
	)
	file, err = os.Open(a_strPath)
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to open the file : "+err.Error())
		return lstData
	}

	reader = csv.NewReader(file)

	arrFullRecords, err = reader.ReadAll()
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to read the records : "+err.Error())
		return lstData
	}

	// Inicia da linha 1 (pula o header)
	for _, arrRecord = range arrFullRecords[1:] {
		// Verifica tamanho da linha
		if len(arrRecord) == c_nLastIndex+1 {
			// Verifica nome do ticker
			if arrRecord[c_nTickerIndex] == a_strTicker {
				// Verifica natureza da operacao
				if arrRecord[c_nOperationIndex][0] == 'C' {
					TradeData.chOperation = tropBuy
				} else if arrRecord[c_nOperationIndex][0] == 'V' {
					TradeData.chOperation = tropSell
				} else {
					logger.LogError(m_strLogFile, c_strMethodName, "Invalid operation type : "+arrRecord[c_nOperationIndex])
					TradeData.chOperation = tropUnknown
				}
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker : "+arrRecord[c_nTickerIndex])
			}
		}
	}

	defer file.Close()

	return lstData
}

func loadOfferDataFromFile(a_strPath string) list.List {
}
