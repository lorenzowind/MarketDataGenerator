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
		err              error
		strInputPath     string
		strBuyPath       string
		strSellPath      string
		strBenchmarkPath string
		bFileExists      bool
		FilesInfo        FilesInfoType
	)
	strInputPath = getInputPath() + "/"

	strBuyPath = strInputPath + fmt.Sprintf(c_strBuyFile, a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.strTickerName)
	bFileExists = checkFileExists(strBuyPath)

	if bFileExists {
		logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, "Buy file found : strBuyPath="+strBuyPath)
	} else {
		strBuyPath = ""
	}

	strSellPath = strInputPath + fmt.Sprintf(c_strSellFile, a_TradeRunInfo.dtTickerDate.Year(), a_TradeRunInfo.dtTickerDate.Month(), a_TradeRunInfo.dtTickerDate.Day(), a_TradeRunInfo.strTickerName)
	bFileExists = checkFileExists(strSellPath)

	if bFileExists {
		logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, "Sell file found : strSellPath="+strSellPath)
	} else {
		strSellPath = ""
	}

	strBenchmarkPath = strInputPath + c_strBenchmarksFile
	bFileExists = checkFileExists(strBenchmarkPath)

	if bFileExists {
		logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, "Benchmarks file found : strBenchmarkPath="+strBenchmarkPath)
	} else {
		strBenchmarkPath = ""
	}

	// Existe os 2 arquivos (compra e venda) ou existe pelo menos o de compra ou venda, alem do arquivo de benchmarks
	if strBuyPath != "" || strSellPath != "" {
		FilesInfo = FilesInfoType{
			TradeRunInfo:     a_TradeRunInfo,
			strBuyPath:       strBuyPath,
			strSellPath:      strSellPath,
			strBenchmarkPath: strBenchmarkPath,
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
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, "Fail to get the directory : "+err.Error())
		return arrTickersInfo
	}

	// Itera sobre cada arquivo da pasta input
	for _, dir = range arrDir {
		// So verifica se for um arquivo
		if !dir.IsDir() {
			strFileName = filepath.Base(dir.Name())
			arrFileInfo = strings.Split(strFileName, "_")

			// So verifica arquivo no formato yyyy-mm-dd_<TICKER>_<BUY/SELL>.csv
			if len(arrFileInfo) == 3 {
				dtTickerDate, err = validateDateString(arrFileInfo[0])
				if err == nil {
					// So verifica se tem a mesma data do que foi passado via parametro e nao foi adicionado ainda
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
	const (
		c_strMethodName = "reader.loadTickerData"
	)
	var (
		TickerData TickerDataType
	)
	logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_FilesInfo.TradeRunInfo)+" : Begin")

	TickerData.FilesInfo = &a_FilesInfo

	TickerData.AuxiliarData.hshFullTrade = make(map[int]*FullTradeType)
	TickerData.AuxiliarData.hshOffersByPrimary = make(map[int][]*OfferDataType)
	TickerData.AuxiliarData.hshOffersBySecondary = make(map[int][]*OfferDataType)
	TickerData.AuxiliarData.hshTradesBySecondary = make(map[int][]*OfferDataType)
	TickerData.AuxiliarData.hshTradesByAccount = make(map[string][]*FullTradeType)

	// Carrega dados do arquivo de compra
	if a_FilesInfo.strBuyPath != "" {
		loadOfferDataFromFile(a_FilesInfo.strBuyPath, &TickerData, true)
	}

	// Carrega dados do arquivo de venda
	if a_FilesInfo.strSellPath != "" {
		loadOfferDataFromFile(a_FilesInfo.strSellPath, &TickerData, false)
	}

	// Carrega dados de benchmark
	if a_FilesInfo.strBenchmarkPath != "" {
		TickerData.AuxiliarData.BenchmarkData.bHasBenchmarkData = tryLoadBenchmarkFromFile(a_FilesInfo.strBenchmarkPath, &TickerData)
		// Verifica se conseguiu encontrar os dados de benchmark do ativo
		if TickerData.AuxiliarData.BenchmarkData.bHasBenchmarkData {
			// Calcula o valor de referencia da oferta expressiva
			TickerData.AuxiliarData.BenchmarkData.sExpressiveOfferSize = calculateExpressiveOfferSize(&TickerData)

			TickerData.TempData.hshTradePrice = make(map[int]TradePriceType)
		}
	}

	logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_FilesInfo.TradeRunInfo)+" : Ticker data loaded successfully : "+getTickerData(&TickerData))
	logger.Log(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_FilesInfo.TradeRunInfo)+" : End")

	return TickerData
}

func tryLoadBenchmarkFromFile(a_strPath string, a_TickerData *TickerDataType) bool {
	const (
		c_strMethodName            = "reader.tryLoadBenchmarkFromFile"
		c_nTickerIndex             = 0
		c_nAvgTradeIntervalIndex   = 1
		c_nAvgOfferSizeIndex       = 2
		c_nSmallerSDOfferSizeIndex = 3
		c_nBiggerSDOfferSizeIndex  = 4
		c_nLastIndex               = 4
	)
	var (
		err            error
		arrRecord      []string
		arrFullRecords [][]string
		file           *os.File
		reader         *csv.Reader
		bFound         bool
	)
	bFound = false

	file, err = os.Open(a_strPath)
	if err == nil {
		reader = csv.NewReader(file)
		reader.Comma = ';'

		arrFullRecords, err = reader.ReadAll()
		if err == nil {
			// Inicia da linha 1 (pula o header)
			for _, arrRecord = range arrFullRecords[1:] {
				// Verifica tamanho da linha
				if len(arrRecord) != c_nLastIndex+1 {
					logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica se encontrou benchmark do ticker
				if a_TickerData.FilesInfo.TradeRunInfo.strTickerName == arrRecord[c_nTickerIndex] {
					// Verifica benchmark de intervalo entre negocios
					a_TickerData.AuxiliarData.BenchmarkData.dtAvgTradeInterval = getDurationFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nAvgTradeIntervalIndex)
					// Verifica benchmark da media da quantidade de lotes
					a_TickerData.AuxiliarData.BenchmarkData.sAvgOfferSize = getAvgOfferSizeFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nAvgOfferSizeIndex)
					// Verifica benchmark do desvio padrao da quantidade de lotes
					a_TickerData.AuxiliarData.BenchmarkData.sSDOfferSize = getSDOfferSizeFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nSmallerSDOfferSizeIndex, c_nBiggerSDOfferSizeIndex)

					bFound = true
					break
				}
			}

			if !bFound {
				logger.LogWarning(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Benchmark for ticker not found : strTicker="+a_TickerData.FilesInfo.TradeRunInfo.strTickerName)
			}

			defer file.Close()
		} else {
			logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Fail to open the file : "+err.Error())
	}

	return bFound
}

func loadOfferDataFromFile(a_strPath string, a_TickerData *TickerDataType, bBuy bool) {
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
		lstData        *list.List
		arrRecord      []string
		arrFullRecords [][]string
		file           *os.File
		reader         *csv.Reader
		OfferData      OfferDataType
	)
	file, err = os.Open(a_strPath)
	if err == nil {
		reader = csv.NewReader(file)
		reader.Comma = ';'

		arrFullRecords, err = reader.ReadAll()
		if err == nil {
			if bBuy {
				lstData = &a_TickerData.lstBuy
			} else {
				lstData = &a_TickerData.lstSell
			}

			// Inicia da linha 1 (pula o header)
			for _, arrRecord = range arrFullRecords[1:] {
				// Verifica tamanho da linha
				if len(arrRecord) != c_nLastIndex+1 {
					logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica natureza da operacao
				OfferData.nOperation = getOfferOperationFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nOperationIndex)
				// Verifica timestamp da oferta
				OfferData.dtTime = getTimeFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nTimeIndex)
				// Verifica numero de geracao da oferta
				OfferData.nGenerationID = getOfferGenerationIDFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nGenerationIDIndex)
				// Verifica conta
				OfferData.strAccount = arrRecord[c_nAccountIndex]
				// Verifica numero do negocio relacionado
				OfferData.nTradeID = getTradeIDFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nTradeIDIndex)
				// Verifica numero primario da oferta
				OfferData.nPrimaryID = getOfferPrimaryIDFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nPrimaryIDIndex)
				// Verifica numero secundario da oferta
				OfferData.nSecondaryID = getOfferSecondaryIDFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nSecondaryIDIndex)
				// Verifica quantidade restante
				OfferData.nCurrentQuantity = getCurrentQuantityFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nCurrentQuantityIndex)
				// Verifica quantidade negociada ate o momento
				OfferData.nTradeQuantity = getTradeQuantityFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nTradeQuantityIndex)
				// Verifica quantidade total
				OfferData.nTotalQuantity = getTotalQuantityFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nTotalQuantityIndex)
				// Verifica preco
				OfferData.sPrice = getPriceFromFile(a_TickerData.FilesInfo.TradeRunInfo, arrRecord, c_nPriceIndex)

				lstData.PushBack(OfferData)

				relateOfferIntoAuxiliarData(a_TickerData, OfferData, bBuy)
			}

			defer file.Close()
		} else {
			logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Fail to open the file : "+err.Error())
	}
}

func relateOfferIntoAuxiliarData(a_TickerData *TickerDataType, a_OfferData OfferDataType, bBuy bool) {
	var (
		FullTradeAux FullTradeType
		FullTrade    *FullTradeType
		lstFullTrade []*FullTradeType
		lstOfferData []*OfferDataType
		bKeyExists   bool
	)
	if a_OfferData.nOperation == ofopTrade {
		// Relaciona evento da oferta referente a ocorrencia de um negocio
		FullTrade, bKeyExists = a_TickerData.AuxiliarData.hshFullTrade[a_OfferData.nTradeID]
		if !bKeyExists {
			FullTrade = &FullTradeAux
			a_TickerData.AuxiliarData.hshFullTrade[a_OfferData.nTradeID] = FullTrade

			// Relaciona conta do investidor referente a um evento da oferta
			lstFullTrade, bKeyExists = a_TickerData.AuxiliarData.hshTradesByAccount[a_OfferData.strAccount]
			if !bKeyExists {
				lstFullTrade = make([]*FullTradeType, 0)
				a_TickerData.AuxiliarData.hshTradesByAccount[a_OfferData.strAccount] = lstFullTrade
			}
			a_TickerData.AuxiliarData.hshTradesByAccount[a_OfferData.strAccount] = append(lstFullTrade, FullTrade)
		}

		if bBuy {
			FullTrade.BuyOfferTrade = &a_OfferData
		} else {
			FullTrade.SellOfferTrade = &a_OfferData
		}

		// Relaciona ID secundário com evento de trade
		lstOfferData, bKeyExists = a_TickerData.AuxiliarData.hshTradesBySecondary[a_OfferData.nSecondaryID]
		if !bKeyExists {
			lstOfferData = make([]*OfferDataType, 0)
		}
		a_TickerData.AuxiliarData.hshTradesBySecondary[a_OfferData.nSecondaryID] = append(lstOfferData, &a_OfferData)
	}
	// Relaciona ID primário do evento da oferta
	lstOfferData, bKeyExists = a_TickerData.AuxiliarData.hshOffersByPrimary[a_OfferData.nPrimaryID]
	if !bKeyExists {
		lstOfferData = make([]*OfferDataType, 0)
	}
	a_TickerData.AuxiliarData.hshOffersByPrimary[a_OfferData.nPrimaryID] = append(lstOfferData, &a_OfferData)
	// Relaciona ID secundário do evento da oferta
	lstOfferData, bKeyExists = a_TickerData.AuxiliarData.hshOffersBySecondary[a_OfferData.nSecondaryID]
	if !bKeyExists {
		lstOfferData = make([]*OfferDataType, 0)
	}
	a_TickerData.AuxiliarData.hshOffersBySecondary[a_OfferData.nSecondaryID] = append(lstOfferData, &a_OfferData)
}

func getOfferOperationFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) OfferOperationType {
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
	logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer operation type : "+a_arrRecord[a_nIndex])
	return ofopUnknown
}

func getDurationFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) time.Duration {
	return getTimeFromFile(a_TradeRunInfo, a_arrRecord, a_nIndex).Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
}

func getTimeFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) time.Time {
	const (
		c_strMethodName = "reader.getTimeFromFile"
	)
	var (
		err    error
		dtTime time.Time
	)
	dtTime, err = validateTimestampString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid timestamp : "+err.Error())
	}
	return dtTime
}

func getOfferGenerationIDFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferGenerationIDFromFile"
	)
	var (
		err                error
		nOfferGenerationID int
	)
	nOfferGenerationID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer generation ID : "+err.Error())
	}
	return nOfferGenerationID
}

func getTradeIDFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTradeIDFromFile"
	)
	var (
		err error
		nID int
	)
	nID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid trade ID : "+err.Error())
	}
	return nID
}

func getOfferPrimaryIDFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferPrimaryIDFromFile"
	)
	var (
		err             error
		nOfferPrimaryID int
	)
	nOfferPrimaryID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer primary ID : "+err.Error())
	}
	return nOfferPrimaryID
}

func getOfferSecondaryIDFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferSecondaryIDFromFile"
	)
	var (
		err               error
		nOfferSecondaryID int
	)
	nOfferSecondaryID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer secondary ID : "+err.Error())
	}
	return nOfferSecondaryID
}

func getTradeQuantityFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTradeQuantityFromFile"
	)
	var (
		err       error
		nQuantity int
	)
	nQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid trade quantity : "+err.Error())
	}
	return nQuantity
}

func getPriceFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) float64 {
	const (
		c_strMethodName = "reader.getPriceFromFile"
	)
	var (
		err    error
		sPrice float64
	)
	sPrice, err = validateFloatString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid price : "+err.Error())
	}
	return sPrice
}

func getCurrentQuantityFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getCurrentQuantityFromFile"
	)
	var (
		err              error
		nCurrentQuantity int
	)
	nCurrentQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid current quantity : "+err.Error())
	}
	return nCurrentQuantity
}

func getTotalQuantityFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getTotalQuantityFromFile"
	)
	var (
		err            error
		nTotalQuantity int
	)
	nTotalQuantity, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid total quantity : "+err.Error())
	}
	return nTotalQuantity
}

func getAvgOfferSizeFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nIndex int) float64 {
	const (
		c_strMethodName = "reader.getAvgOfferSizeFromFile"
	)
	var (
		err           error
		sAvgOfferSize float64
	)
	sAvgOfferSize, err = validateFloatString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer size average : "+err.Error())
	}
	return sAvgOfferSize
}

func getSDOfferSizeFromFile(a_TradeRunInfo TradeRunInfoType, a_arrRecord []string, a_nSmallerIndex int, a_nBiggerIndex int) float64 {
	const (
		c_strMethodName = "reader.getSDOfferSizeFromFile"
	)
	var (
		err                 error
		sSmallerSDOfferSize float64
		sBiggerSDOfferSize  float64
	)
	sSmallerSDOfferSize, err = validateFloatString(a_arrRecord[a_nSmallerIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer size smaller sd : "+err.Error())
		return 0
	}
	sBiggerSDOfferSize, err = validateFloatString(a_arrRecord[a_nBiggerIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Ticker-External-Data", c_strMethodName, getHeaderRun(a_TradeRunInfo)+" : Invalid offer size smaller sd : "+err.Error())
		return 0
	}

	return (sSmallerSDOfferSize + sBiggerSDOfferSize) / 2
}

func calculateExpressiveOfferSize(a_TickerData *TickerDataType) float64 {
	const (
		c_nMultiplier = 3
	)
	return a_TickerData.AuxiliarData.BenchmarkData.sAvgOfferSize + (c_nMultiplier * a_TickerData.AuxiliarData.BenchmarkData.sSDOfferSize)
}

func getOffersByPrimaryID(a_TickerData *TickerDataType, a_nPrimaryID int) []*OfferDataType {
	var (
		lstOfferData []*OfferDataType
		bKeyExists   bool
	)
	lstOfferData, bKeyExists = a_TickerData.AuxiliarData.hshOffersByPrimary[a_nPrimaryID]
	if bKeyExists {
		return lstOfferData
	}
	return make([]*OfferDataType, 0)
}

func getTradesBySecondaryID(a_TickerData *TickerDataType, a_nSecondaryID int) []*OfferDataType {
	var (
		lstOfferData []*OfferDataType
		bKeyExists   bool
	)
	lstOfferData, bKeyExists = a_TickerData.AuxiliarData.hshTradesBySecondary[a_nSecondaryID]
	if bKeyExists {
		return lstOfferData
	}
	return make([]*OfferDataType, 0)
}

func getTradesByAccount(a_TickerData *TickerDataType, a_strAccount string) []*FullTradeType {
	var (
		lstFullTrade []*FullTradeType
		bKeyExists   bool
	)
	lstFullTrade, bKeyExists = a_TickerData.AuxiliarData.hshTradesByAccount[a_strAccount]
	if bKeyExists {
		return lstFullTrade
	}
	return make([]*FullTradeType, 0)
}

func getTradePrice(a_TickerData *TickerDataType, a_bBuy bool, a_nTradeID int) float64 {
	var (
		TradePrice TradePriceType
		bKeyExists bool
	)
	TradePrice, bKeyExists = a_TickerData.TempData.hshTradePrice[a_nTradeID]
	if bKeyExists {
		if a_bBuy {
			return TradePrice.sTopBuyPriceLevel
		}
		return TradePrice.sTopSellPriceLevel
	}
	return 0
}

func getFullTrade(a_TickerData *TickerDataType, a_nTradeID int) *FullTradeType {
	var (
		FullTrade  *FullTradeType
		bKeyExists bool
	)
	FullTrade, bKeyExists = a_TickerData.AuxiliarData.hshFullTrade[a_nTradeID]
	if bKeyExists {
		return FullTrade
	}
	return nil
}
