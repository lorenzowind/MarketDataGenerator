package src

import (
	"container/list"
	"encoding/csv"
	"errors"
	"fmt"
	logger "marketmanipulationdetector/logger/src"
	"os"
	"strconv"
	"strings"
	"time"
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

	strBuyPath = strReferencePath + fmt.Sprintf(c_strReferenceBuyFile, a_GenerationInfo.dtReferenceTickerDate.Day(), a_GenerationInfo.dtReferenceTickerDate.Month(), a_GenerationInfo.dtReferenceTickerDate.Year(), a_GenerationInfo.strReferenceTickerName)
	bFileExists = checkFileExists(strBuyPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Buy reference file found : strBuyPath="+strBuyPath)
	} else {
		strBuyPath = ""
	}

	strSellPath = strReferencePath + fmt.Sprintf(c_strReferenceSellFile, a_GenerationInfo.dtReferenceTickerDate.Day(), a_GenerationInfo.dtReferenceTickerDate.Month(), a_GenerationInfo.dtReferenceTickerDate.Year(), a_GenerationInfo.strReferenceTickerName)
	bFileExists = checkFileExists(strSellPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Sell reference file found : strSellPath="+strSellPath)
	} else {
		strSellPath = ""
	}

	strBenchmarkPath = strReferencePath + c_strBenchmarksFile
	bFileExists = checkFileExists(strBenchmarkPath)

	if bFileExists {
		logger.Log(m_strLogFile, c_strMethodName, "Benchmarks reference file found : strBenchmarkPath="+strBenchmarkPath)
	} else {
		strBenchmarkPath = ""
	}

	// Existe os 2 arquivos (compra e venda) ou existe pelo menos o de compra ou venda, alem do arquivo de benchmarks
	if strBuyPath != "" || strSellPath != "" {
		FilesInfo = FilesInfoType{
			GenerationInfo: GenerationInfoType{
				strTickerName:          a_GenerationInfo.strTickerName,
				dtTickerDate:           a_GenerationInfo.dtTickerDate,
				strReferenceTickerName: a_GenerationInfo.strReferenceTickerName,
				dtReferenceTickerDate:  a_GenerationInfo.dtReferenceTickerDate,
			},
			strReferenceBuyPath:       strBuyPath,
			strReferenceSellPath:      strSellPath,
			strReferenceBenchmarkPath: strBenchmarkPath,
			strBuyPath:                "", // Vazio pois sera utilizado posteriormente
			strSellPath:               "", // Vazio pois sera utilizado posteriormente
			strBenchmarkPath:          "", // Vazio pois sera utilizado posteriormente
		}
	} else {
		err = errors.New("file rules existance failed")
	}

	return FilesInfo, err
}

func duplicateOffersBook(a_FilesInfo *FilesInfoType) {
	const (
		c_strMethodName = "reader.duplicateOffersBook"
	)
	strInputPath := getInputPath() + "/"
	// Duplica arquivo de compra (referencia -> geracao)
	if a_FilesInfo.strReferenceBuyPath != "" {
		a_FilesInfo.strBuyPath = strInputPath + fmt.Sprintf(c_strBuyFile, a_FilesInfo.GenerationInfo.dtTickerDate.Year(), a_FilesInfo.GenerationInfo.dtTickerDate.Month(), a_FilesInfo.GenerationInfo.dtTickerDate.Day(), a_FilesInfo.GenerationInfo.strTickerName)

		if duplicateFile(a_FilesInfo.strReferenceBuyPath, a_FilesInfo.strBuyPath) {
			logger.Log(m_strLogFile, c_strMethodName, "Buy file duplicated successfully : strBuyPath="+a_FilesInfo.strBuyPath)
		}
	}
	// Duplica arquivo de venda (referencia -> geracao)
	if a_FilesInfo.strReferenceSellPath != "" {
		a_FilesInfo.strSellPath = strInputPath + fmt.Sprintf(c_strSellFile, a_FilesInfo.GenerationInfo.dtTickerDate.Year(), a_FilesInfo.GenerationInfo.dtTickerDate.Month(), a_FilesInfo.GenerationInfo.dtTickerDate.Day(), a_FilesInfo.GenerationInfo.strTickerName)

		if duplicateFile(a_FilesInfo.strReferenceSellPath, a_FilesInfo.strSellPath) {
			logger.Log(m_strLogFile, c_strMethodName, "Sell file duplicated successfully : strSellPath="+a_FilesInfo.strSellPath)
		}
	}
	// Duplica arquivo de benchmarks (referencia -> geracao)
	if a_FilesInfo.strReferenceBenchmarkPath != "" {
		a_FilesInfo.strBenchmarkPath = strInputPath + c_strBenchmarksFile

		if duplicateFile(a_FilesInfo.strReferenceBenchmarkPath, a_FilesInfo.strBenchmarkPath) {
			logger.Log(m_strLogFile, c_strMethodName, "Benchmark file duplicated successfully : strBenchmarkPath="+a_FilesInfo.strBenchmarkPath)
		}
	}
}

func loadTickerData(a_FilesInfo FilesInfoType) TickerDataType {
	var (
		TickerData TickerDataType
	)
	TickerData.FilesInfo = &a_FilesInfo

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
		TickerData.BenchmarkData.bHasBenchmarkData = tryLoadBenchmarkFromFile(a_FilesInfo.strBenchmarkPath, &TickerData)
	}

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
		reader.Comma = ','

		arrFullRecords, err = reader.ReadAll()
		if err == nil {
			// Inicia da linha 1 (pula o header)
			for _, arrRecord = range arrFullRecords[1:] {
				// Verifica tamanho da linha
				if len(arrRecord) != c_nLastIndex+1 {
					logger.LogError(m_strLogFile, c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica se encontrou benchmark do ticker
				if a_TickerData.FilesInfo.GenerationInfo.strReferenceTickerName == arrRecord[c_nTickerIndex] {
					// Verifica benchmark de intervalo entre negocios
					a_TickerData.BenchmarkData.dtAvgTradeInterval = getDurationFromFile(arrRecord, c_nAvgTradeIntervalIndex)
					// Verifica benchmark da media da quantidade de lotes
					a_TickerData.BenchmarkData.sAvgOfferSize = getAvgOfferSizeFromFile(arrRecord, c_nAvgOfferSizeIndex)
					// Verifica benchmark do desvio padrao da quantidade de lotes
					a_TickerData.BenchmarkData.sSDOfferSize = getSDOfferSizeFromFile(arrRecord, c_nSmallerSDOfferSizeIndex, c_nBiggerSDOfferSizeIndex)

					bFound = true
					break
				}
			}

			if !bFound {
				logger.LogWarning(m_strLogFile, c_strMethodName, "Benchmark for ticker not found : strTicker="+a_TickerData.FilesInfo.GenerationInfo.strReferenceTickerName)
			}

			defer file.Close()
		} else {
			logger.LogError(m_strLogFile, c_strMethodName, "Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to open the file : "+err.Error())
	}

	return bFound
}

func loadOfferDataFromFile(a_strPath string, a_TickerData *TickerDataType, bBuy bool) {
	const (
		c_strMethodName         = "reader.loadOfferDataFromFile"
		c_nOperationIndex       = 0
		c_nTickerIndex          = 3
		c_nTimeIndex            = 6
		c_nGenerationIDIndex    = 7
		c_nAccountIndex         = 8
		c_nTradeIDIndex         = 9
		c_nPrimaryIDIndex       = 10
		c_nSecondaryIDIndex     = 11
		c_nCurrentQuantityIndex = 12
		c_nTradeQuantityIndex   = 13
		c_nTotalQuantityIndex   = 14
		c_nPriceIndex           = 15
		c_nLastIndex            = 15
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
		reader.Comma = '|'

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
					logger.LogError(m_strLogFile, c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica natureza da operacao
				OfferData.chOperation = getOfferOperationFromFile(arrRecord, c_nOperationIndex)
				// Normaliza timestamp da oferta
				OfferData.dtTime = normalizeTime(getTimeFromFile(arrRecord, c_nTimeIndex), a_TickerData.FilesInfo.GenerationInfo.dtReferenceTickerDate, a_TickerData.FilesInfo.GenerationInfo.dtTickerDate)
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
		} else {
			logger.LogError(m_strLogFile, c_strMethodName, "Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_strLogFile, c_strMethodName, "Fail to open the file : "+err.Error())
	}
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

func getDurationFromFile(a_arrRecord []string, a_nIndex int) time.Duration {
	return getTimeFromFile(a_arrRecord, a_nIndex).Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
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

func normalizeTime(a_dtLoadedTime, a_dtReferenceTime, a_dtTime time.Time) time.Time {
	const (
		c_strMethodName = "reader.normalizeTime"
	)
	var (
		err             error
		dtLoadedDate    time.Time
		dtReferenceDate time.Time
		dtTime          time.Time
	)
	// Obtem somente a data do timestamp lido
	dtLoadedDate, err = validateDateString(fmt.Sprintf(c_strCustomDateFormat, a_dtLoadedTime.Year(), a_dtLoadedTime.Month(), a_dtLoadedTime.Day()))
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid loaded date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem somente a data da referencia
	dtReferenceDate, err = validateDateString(fmt.Sprintf(c_strCustomDateFormat, a_dtReferenceTime.Year(), a_dtReferenceTime.Month(), a_dtReferenceTime.Day()))
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid reference date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem somente a data de geracao
	dtTime, err = validateDateString(fmt.Sprintf(c_strCustomDateFormat, a_dtTime.Year(), a_dtTime.Month(), a_dtTime.Day()))
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem o resultado da subtracao entre a data lida e a referencia
	dtTime = dtTime.Add(dtLoadedDate.Sub(dtReferenceDate))
	// Obtem a nova data com valor normalizado
	dtTime, err = validateTimestampString(fmt.Sprintf(c_strCustomTimestampFormat, dtTime.Year(), dtTime.Month(), dtTime.Day(), a_dtLoadedTime.Hour(), a_dtLoadedTime.Minute(), a_dtLoadedTime.Second(), a_dtLoadedTime.Nanosecond()))
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid timestamp : "+err.Error())
	}
	return dtTime
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

func getAvgOfferSizeFromFile(a_arrRecord []string, a_nIndex int) float64 {
	const (
		c_strMethodName = "reader.getAvgOfferSizeFromFile"
	)
	var (
		err           error
		sAvgOfferSize float64
	)
	sAvgOfferSize, err = validateFloatString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer size average : "+err.Error())
	}
	return sAvgOfferSize
}

func getSDOfferSizeFromFile(a_arrRecord []string, a_nSmallerIndex int, a_nBiggerIndex int) float64 {
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
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer size smaller sd : "+err.Error())
		return 0
	}
	sBiggerSDOfferSize, err = validateFloatString(a_arrRecord[a_nBiggerIndex])
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid offer size smaller sd : "+err.Error())
		return 0
	}

	return (sSmallerSDOfferSize + sBiggerSDOfferSize) / 2
}

func saveOffersBook(a_TickerData TickerDataType) {

}
