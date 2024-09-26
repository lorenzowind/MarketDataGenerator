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
	"strconv"
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

func getAllTickersFiles() []FilesInfoType {
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
					// if checkIfHasSameDate(dtTickerDate, a_dtTickerDate) && !checkIfContains(arrFileInfo[1], arrTickersInfo) {
					if !checkIfContains(arrFileInfo[1], arrTickersInfo) {
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
	//if a_FilesInfo.strTradePath != "" {
	//	TickerData.lstTrade = loadTradeDataFromFile(a_FilesInfo.strTradePath, a_FilesInfo.TradeRunInfo.strTickerName)
	//}

	// Carrega dados do arquivo de compra
	if a_FilesInfo.strBuyPath != "" {
		TickerData.lstBuy = loadOfferDataFromFile(a_FilesInfo.strBuyPath, a_FilesInfo.TradeRunInfo.strTickerName)
	}

	// Carrega dados do arquivo de venda
	if a_FilesInfo.strSellPath != "" {
		TickerData.lstSell = loadOfferDataFromFile(a_FilesInfo.strSellPath, a_FilesInfo.TradeRunInfo.strTickerName)
	}

	TickerData.FilesInfo = a_FilesInfo
	return TickerData
}

//lint:ignore U1000 Ignore unused function
func loadTradeDataFromFile(a_strPath, a_strTicker string) list.List {
	const (
		c_strMethodName           = "reader.loadTradeDataFromFile"
		c_nOperationIndex         = 0
		c_nTickerIndex            = 1
		c_nTimeIndex              = 2
		c_nOfferGenerationIDIndex = 3
		c_nAccountIndex           = 4
		c_nIDIndex                = 5
		c_nOfferPrimaryIDIndex    = 6
		c_nOfferSecondaryIDIndex  = 7
		c_nQuantityIndex          = 8
		c_nPriceIndex             = 9
		c_nLastIndex              = 9
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
	reader.Comma = '|'

	arrFullRecords, err = reader.ReadAll()
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to read the records : "+err.Error())
		return lstData
	}

	// Inicia da linha 1 (pula o header)
	for _, arrRecord = range arrFullRecords[1:] {
		// Verifica tamanho da linha
		if len(arrRecord) != c_nLastIndex+1 {
			logger.LogError(m_strLogFile, c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
			continue
		}
		// Verifica nome do ticker
		if arrRecord[c_nTickerIndex] != a_strTicker {
			logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker : "+arrRecord[c_nTickerIndex])
			continue
		}
		// Verifica natureza da operacao
		TradeData.chOperation = getTradeOperationFromFile(arrRecord, c_nOperationIndex)
		// Verifica timestamp do negocio
		TradeData.dtTime = getTimeFromFile(arrRecord, c_nTimeIndex)
		// Verifica numero de geracao da oferta
		TradeData.nOfferGenerationID = getOfferGenerationFromFile(arrRecord, c_nOfferGenerationIDIndex)
		// Verifica conta
		TradeData.strAccount = arrRecord[c_nAccountIndex]
		// Verifica numero do negocio
		TradeData.nID = getTradeIDFromFile(arrRecord, c_nIDIndex)
		// Verifica numero primario da oferta relacionada
		TradeData.nOfferPrimaryID = getOfferPrimaryIDFromFile(arrRecord, c_nOfferPrimaryIDIndex)
		// Verifica numero secundario da oferta relacionada
		TradeData.nOfferSecondaryID = getOfferSecondaryIDFromFile(arrRecord, c_nOfferSecondaryIDIndex)
		// Verifica quantidade
		TradeData.nQuantity = getTradeQuantityFromFile(arrRecord, c_nQuantityIndex)
		// Verifica preco
		TradeData.sPrice = getPriceFromFile(arrRecord, c_nPriceIndex)

		lstData.PushBack(TradeData)
	}

	defer file.Close()

	return lstData
}

func getTimeFromFile(a_arrRecord []string, a_nIndex int) time.Time {
	const (
		c_strMethodName = "reader.getTimeFromFile"
	)
	var (
		err    error
		dtTime time.Time
	)
	dtTime, err = validateTimestampString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid timestamp : "+err.Error())
	}
	return dtTime
}

func getTradeOperationFromFile(a_arrRecord []string, a_nIndex int) TradeOperationType {
	const (
		c_strMethodName = "reader.getTradeOperationFromFile"
	)
	if a_arrRecord[a_nIndex][0] == 'C' {
		return tropBuy
	} else if a_arrRecord[a_nIndex][0] == 'V' {
		return tropSell
	}
	logger.LogError(m_strLogFile, c_strMethodName, "Invalid trade operation type : "+a_arrRecord[a_nIndex])
	return tropUnknown
}

func getOfferGenerationFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferGenerationFromFile"
	)
	var (
		err                error
		nOfferGenerationID int
	)
	nOfferGenerationID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer generation ID : "+err.Error())
	}
	return nOfferGenerationID
}

func getTradeIDFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTradeIDFromFile"
	)
	var (
		err error
		nID int
	)
	nID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid trade ID : "+err.Error())
	}
	return nID
}

func getOfferPrimaryIDFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferPrimaryIDFromFile"
	)
	var (
		err             error
		nOfferPrimaryID int
	)
	nOfferPrimaryID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer primary ID : "+err.Error())
	}
	return nOfferPrimaryID
}

func getOfferSecondaryIDFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferSecondaryIDFromFile"
	)
	var (
		err               error
		nOfferSecondaryID int
	)
	nOfferSecondaryID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer secondary ID : "+err.Error())
	}
	return nOfferSecondaryID
}

func getTradeQuantityFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTradeQuantityFromFile"
	)
	var (
		err       error
		nQuantity int
	)
	nQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid trade quantity : "+err.Error())
	}
	return nQuantity
}

func getPriceFromFile(a_arrRecord []string, a_nIndex int) float64 {
	const (
		c_strMethodName = "reader.getPriceFromFile"
	)
	var (
		err    error
		sPrice float64
	)
	sPrice, err = validateFloatString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid price : "+err.Error())
	}
	return sPrice
}

func loadOfferDataFromFile(a_strPath, a_strTicker string) list.List {
	const (
		c_strMethodName         = "reader.loadOfferDataFromFile"
		c_nOperationIndex       = 0
		c_nTickerIndex          = 1
		c_nTimeIndex            = 2
		c_nGenerationIDIndex    = 3
		c_nAccountIndex         = 4
		c_nTradeIDIndex         = 5
		c_nPrimaryIDIndex       = 6
		c_nSecondaryIDIndex     = 7
		c_nCurrentQuantityIndex = 8
		c_nTradeQuantityIndex   = 9
		c_nTotalQuantityIndex   = 10
		c_nPriceIndex           = 11
		c_nLastIndex            = 11
	)
	var (
		err            error
		lstData        list.List
		arrRecord      []string
		arrFullRecords [][]string
		file           *os.File
		reader         *csv.Reader
		OfferData      OfferDataType
	)
	file, err = os.Open(a_strPath)
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to open the file : "+err.Error())
		return lstData
	}

	reader = csv.NewReader(file)
	reader.Comma = ';'

	arrFullRecords, err = reader.ReadAll()
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to read the records : "+err.Error())
		return lstData
	}

	// Inicia da linha 1 (pula o header)
	for _, arrRecord = range arrFullRecords[1:] {
		// Verifica tamanho da linha
		if len(arrRecord) != c_nLastIndex+1 {
			logger.LogError(m_strLogFile, c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
			continue
		}
		// Verifica nome do ticker
		if arrRecord[c_nTickerIndex] != a_strTicker {
			logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker : "+arrRecord[c_nTickerIndex])
			continue
		}
		// Verifica natureza da operacao
		OfferData.chOperation = getOfferOperationFromFile(arrRecord, c_nOperationIndex)
		// Verifica timestamp da oferta
		OfferData.dtTime = getTimeFromFile(arrRecord, c_nTimeIndex)
		// Verifica numero de geracao da oferta
		OfferData.nGenerationID = getOfferGenerationFromFile(arrRecord, c_nGenerationIDIndex)
		// Verifica conta
		OfferData.strAccount = arrRecord[c_nAccountIndex]
		// Verifica numero do negocio relacionado
		OfferData.nTradeID = getTradeIDFromFile(arrRecord, c_nTradeIDIndex)
		// Verifica numero primario da oferta
		OfferData.nPrimaryID = getOfferPrimaryIDFromFile(arrRecord, c_nPrimaryIDIndex)
		// Verifica numero secundario da oferta
		OfferData.nSecondaryID = getOfferSecondaryIDFromFile(arrRecord, c_nSecondaryIDIndex)
		// Verifica quantidade restante
		OfferData.nCurrentQuantity = getCurrentQuantityFromFile(arrRecord, c_nCurrentQuantityIndex)
		// Verifica quantidade negociada ate o momento
		OfferData.nTradeQuantity = getTradeQuantityFromFile(arrRecord, c_nTradeQuantityIndex)
		// Verifica quantidade total
		OfferData.nTotalQuantity = getTotalQuantityFromFile(arrRecord, c_nTotalQuantityIndex)
		// Verifica preco
		OfferData.sPrice = getPriceFromFile(arrRecord, c_nPriceIndex)

		lstData.PushBack(OfferData)
	}

	defer file.Close()

	return lstData
}

func getOfferOperationFromFile(a_arrRecord []string, a_nIndex int) OfferOperationType {
	const (
		c_strMethodName = "reader.getOfferOperationFromFile"
	)
	if a_arrRecord[a_nIndex][0] == '0' {
		return ofopCreation
	} else if a_arrRecord[a_nIndex][0] == '4' {
		return ofopCancel
	} else if a_arrRecord[a_nIndex][0] == '5' {
		return ofopEdit
	} else if a_arrRecord[a_nIndex][0] == 'F' {
		return ofopTrade
	} else if a_arrRecord[a_nIndex][0] == 'C' {
		return ofopExpired
	} else if a_arrRecord[a_nIndex][0] == 'D' {
		return ofopReafirmed
	}
	logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer operation type : "+a_arrRecord[a_nIndex])
	return ofopUnknown
}

func getCurrentQuantityFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getCurrentQuantityFromFile"
	)
	var (
		err              error
		nCurrentQuantity int
	)
	nCurrentQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid current quantity : "+err.Error())
	}
	return nCurrentQuantity
}

func getTotalQuantityFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTotalQuantityFromFile"
	)
	var (
		err            error
		nTotalQuantity int
	)
	nTotalQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid total quantity : "+err.Error())
	}
	return nTotalQuantity
}
