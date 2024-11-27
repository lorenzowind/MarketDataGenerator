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

func generateFilterTickers(a_GenerationRule GenerationRuleType) []GenerationInfoType {
	var (
		err               error
		arrTickers        []string
		nCurrent          int
		strName           string
		dtReference       time.Time
		GenerationInfo    GenerationInfoType
		arrGenerationInfo []GenerationInfoType
	)
	arrGenerationInfo = make([]GenerationInfoType, 0)

	arrTickers = getFilterTickersByBenchmarkFromFile(getReferenceBenchmark(), a_GenerationRule.Pattern)

	for _, strName = range arrTickers {
		// Tenta data do ticker (na pasta de referencia)
		dtReference, err = getReferenceTickerDate(strName, getReferencePath())
		if err != nil {
			// Tenta data do ticker (na pasta de referencia de compra)
			dtReference, err = getReferenceTickerDate(strName, getInputBuyPath())
			if err != nil {
				// Tenta data do ticker (na pasta de referencia de venda)
				dtReference, err = getReferenceTickerDate(strName, getInputSellPath())
			}
		}

		if err == nil {
			// Incrementa o numero de tickers
			nCurrent++
			// Obtem nome do ticker para geracao
			GenerationInfo.strTickerName = a_GenerationRule.strTickerNameRule + strconv.Itoa(nCurrent)
			// Obtem data do ticker para geracao
			GenerationInfo.dtTickerDate = a_GenerationRule.dtTickerDate
			// Obtem nome do ticker de referencia
			GenerationInfo.strReferenceTickerName = strName
			// Obtem data do ticker de referencia
			GenerationInfo.dtReferenceTickerDate = dtReference

			arrGenerationInfo = append(arrGenerationInfo, GenerationInfo)
		}
	}

	return arrGenerationInfo
}

func getReferenceTickerDate(a_strTicker string, a_strReferenceFolder string) (time.Time, error) {
	var (
		err          error
		strFileName  string
		arrFileInfo  []string
		dtTickerDate time.Time
		dir          fs.DirEntry
		arrDir       []fs.DirEntry
	)
	arrDir, err = os.ReadDir(a_strReferenceFolder)

	if err != nil {
		return time.Time{}, err
	}

	// Itera sobre cada arquivo da pasta input
	for _, dir = range arrDir {
		// So verifica se for um arquivo
		if !dir.IsDir() {
			strFileName = filepath.Base(dir.Name())
			arrFileInfo = strings.Split(strFileName, "_")

			// So verifica arquivo no formato ddmmyyyy_<TICKER>_<compra/venda>.csv
			if len(arrFileInfo) == 3 {
				// So verifica arquivo com o mesmo nome do ticker encontrado no arquivo de benchmark
				if arrFileInfo[1] == a_strTicker {
					dtTickerDate, err = validateDateString(arrFileInfo[0])
					if err == nil {
						return dtTickerDate, err
					}
				}
			}
		}
	}

	return time.Time{}, errors.New("reference ticker not found")
}

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
		strInputBuyPath  string
		strInputSellPath string
		bFileExists      bool
		FilesInfo        FilesInfoType
	)

	strReferencePath = getReferencePath() + "/"
	strInputBuyPath = getInputBuyPath() + "/"
	strInputSellPath = getInputSellPath() + "/"

	strBuyPath = fmt.Sprintf(c_strReferenceBuyFile, a_GenerationInfo.dtReferenceTickerDate.Day(), a_GenerationInfo.dtReferenceTickerDate.Month(), a_GenerationInfo.dtReferenceTickerDate.Year(), a_GenerationInfo.strReferenceTickerName)

	// Tenta arquivo de compra na pasta de referencia
	bFileExists = checkFileExists(strReferencePath + strBuyPath)
	if !bFileExists {
		// Tenta arquivo na pasta de referencia de compra
		bFileExists = checkFileExists(strInputBuyPath + strBuyPath)
		if bFileExists {
			strBuyPath = strInputBuyPath + strBuyPath
		}
	} else {
		strBuyPath = strReferencePath + strBuyPath
	}

	if bFileExists {
		logger.Log(m_LogInfo, "Main", c_strMethodName, "Buy reference file found : strBuyPath="+strBuyPath)
	} else {
		strBuyPath = ""
	}

	strSellPath = fmt.Sprintf(c_strReferenceSellFile, a_GenerationInfo.dtReferenceTickerDate.Day(), a_GenerationInfo.dtReferenceTickerDate.Month(), a_GenerationInfo.dtReferenceTickerDate.Year(), a_GenerationInfo.strReferenceTickerName)

	// Tenta arquivo de venda na pasta de referencia
	bFileExists = checkFileExists(strReferencePath + strSellPath)
	if !bFileExists {
		// Tenta arquivo na pasta de referencia de venda
		bFileExists = checkFileExists(strInputSellPath + strSellPath)
		if bFileExists {
			strSellPath = strInputSellPath + strSellPath
		}
	} else {
		strSellPath = strReferencePath + strSellPath
	}

	if bFileExists {
		logger.Log(m_LogInfo, "Main", c_strMethodName, "Sell reference file found : strSellPath="+strSellPath)
	} else {
		strSellPath = ""
	}

	strBenchmarkPath = getReferenceBenchmark()
	if strBenchmarkPath != "" {
		logger.Log(m_LogInfo, "Main", c_strMethodName, "Benchmark file found : strBenchmarkPath="+strBenchmarkPath)
	}

	// Existe os 2 arquivos (compra e venda) ou existe pelo menos o de compra ou venda
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

func getOffersBook(a_FilesInfo *FilesInfoType) {
	//lint:ignore S1021 Ignore merge variable declaration with assignment on next line
	var (
		strInputPath string
	)
	strInputPath = getInputPath() + "/"

	// Salva nome do arquivo de compra (referencia -> geracao)
	if a_FilesInfo.strReferenceBuyPath != "" {
		a_FilesInfo.strBuyPath = strInputPath + fmt.Sprintf(c_strBuyFile, a_FilesInfo.GenerationInfo.dtTickerDate.Year(), a_FilesInfo.GenerationInfo.dtTickerDate.Month(), a_FilesInfo.GenerationInfo.dtTickerDate.Day(), a_FilesInfo.GenerationInfo.strTickerName)
	}
	// Salva nome do arquivo de venda (referencia -> geracao)
	if a_FilesInfo.strReferenceSellPath != "" {
		a_FilesInfo.strSellPath = strInputPath + fmt.Sprintf(c_strSellFile, a_FilesInfo.GenerationInfo.dtTickerDate.Year(), a_FilesInfo.GenerationInfo.dtTickerDate.Month(), a_FilesInfo.GenerationInfo.dtTickerDate.Day(), a_FilesInfo.GenerationInfo.strTickerName)
	}
	// Salva nome do arquivo de benchmarks (referencia -> geracao)
	if a_FilesInfo.strReferenceBenchmarkPath != "" {
		a_FilesInfo.strBenchmarkPath = strInputPath + c_strBenchmarksFile
	}
}

func loadTickerData(a_FilesInfo FilesInfoType) TickerDataType {
	const (
		c_strMethodName = "reader.loadTickerData"
	)
	var (
		TickerData TickerDataType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	TickerData.FilesInfo = &a_FilesInfo

	TickerData.MaskDataInfo.hshMaskAccount = make(map[string]int)
	TickerData.MaskDataInfo.hshMaskPrimaryID = make(map[int]int)
	TickerData.MaskDataInfo.hshMaskSecondaryID = make(map[int]int)

	// Carrega dados do arquivo de compra
	if a_FilesInfo.strBuyPath != "" {
		loadOfferDataFromFile(a_FilesInfo.strReferenceBuyPath, &TickerData, true)
	}

	// Carrega dados do arquivo de venda
	if a_FilesInfo.strSellPath != "" {
		loadOfferDataFromFile(a_FilesInfo.strReferenceSellPath, &TickerData, false)
	}

	// Carrega dados de benchmark
	if a_FilesInfo.strBenchmarkPath != "" {
		TickerData.BenchmarkData.bHasBenchmarkData = tryLoadBenchmarkFromFile(a_FilesInfo.strReferenceBenchmarkPath, &TickerData)
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "Ticker data loaded successfully : strTicker="+a_FilesInfo.GenerationInfo.strTickerName+" : "+getTickerData(TickerData))
	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")

	return TickerData
}

func getFilterTickersByBenchmarkFromFile(a_strPath string, a_Pattern PatternType) []string {
	const (
		c_strMethodName          = "reader.getFilterTickersByBenchmarkFromFile"
		c_nTickerIndex           = 0
		c_nAvgTradeIntervalIndex = 1
		c_nLastIndex             = 4
	)
	var (
		err                error
		arrRecord          []string
		arrFullRecords     [][]string
		file               *os.File
		reader             *csv.Reader
		dtAvgTradeInterval time.Time
		arrFilterTickers   []string
	)
	arrFilterTickers = make([]string, 0)

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
					logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}

				if a_Pattern.nBenchmarkFilter == bfAvgTradeInterval {
					dtAvgTradeInterval = getTimeFromFile(arrRecord, c_nAvgTradeIntervalIndex)

					// Verifica se benchmark se encontra no intervalo fechado definido
					if a_Pattern.dtMinAvgTradeInterval.Compare(dtAvgTradeInterval) <= 0 && a_Pattern.dtMaxAvgTradeInterval.Compare(dtAvgTradeInterval) >= 0 {
						arrFilterTickers = append(arrFilterTickers, arrRecord[c_nTickerIndex])
					}
				}
			}

			defer file.Close()
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to open the file : "+err.Error())
	}

	return arrFilterTickers
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
					logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica se encontrou benchmark do ticker
				if a_TickerData.FilesInfo.GenerationInfo.strReferenceTickerName == arrRecord[c_nTickerIndex] {
					// Verifica benchmark de intervalo entre negocios
					a_TickerData.BenchmarkData.dtAvgTradeInterval = getTimeFromFile(arrRecord, c_nAvgTradeIntervalIndex)
					// Verifica benchmark da media da quantidade de lotes
					a_TickerData.BenchmarkData.sAvgOfferSize = getAvgOfferSizeFromFile(arrRecord, c_nAvgOfferSizeIndex)
					// Verifica benchmark do desvio padrao da quantidade de lotes (min)
					a_TickerData.BenchmarkData.sSmallerSDOfferSize = getSDOfferSizeFromFile(arrRecord, c_nSmallerSDOfferSizeIndex)
					// Verifica benchmark do desvio padrao da quantidade de lotes (max)
					a_TickerData.BenchmarkData.sBiggerSDOfferSize = getSDOfferSizeFromFile(arrRecord, c_nBiggerSDOfferSizeIndex)

					bFound = true
					break
				}
			}

			if !bFound {
				logger.LogWarning(m_LogInfo, "Main", c_strMethodName, "Benchmark for ticker not found : strTicker="+a_TickerData.FilesInfo.GenerationInfo.strReferenceTickerName)
			}

			defer file.Close()
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to open the file : "+err.Error())
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
		nPrimaryID     int
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
					logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid columns size : "+strconv.Itoa(len(arrRecord))+" : arrRecord="+strings.Join(arrRecord, ", "))
					continue
				}
				// Verifica natureza da operacao
				OfferData.nOperation = getOfferOperationFromFile(arrRecord, c_nOperationIndex)
				// Verifica timestamp da oferta e faz normalizacao
				OfferData.dtTime = normalizeTime(getTimeFromFile(arrRecord, c_nTimeIndex), a_TickerData.FilesInfo.GenerationInfo.dtReferenceTickerDate, a_TickerData.FilesInfo.GenerationInfo.dtTickerDate)
				// Verifica numero de geracao da oferta
				OfferData.nGenerationID = getOfferGenerationIDFromFile(arrRecord, c_nGenerationIDIndex)
				// Verifica conta e faz mascaramento
				OfferData.strAccount = maskStringToIntString(arrRecord[c_nAccountIndex], a_TickerData.MaskDataInfo.hshMaskAccount, &a_TickerData.MaskDataInfo.nCurrentAccount)
				// Verifica numero do negocio relacionado
				OfferData.nTradeID = getTradeIDFromFile(arrRecord, c_nTradeIDIndex)

				// Verifica numero primario da oferta e faz mascaramento
				// Alem disso, verifica regra de excecao onde o evento de criacao aparece com o mesmo numero primario da oferta anterior
				nPrimaryID = maskIntToInt(getOfferPrimaryIDFromFile(arrRecord, c_nPrimaryIDIndex), a_TickerData.MaskDataInfo.hshMaskPrimaryID, &a_TickerData.MaskDataInfo.nCurrentPrimaryID)
				if OfferData.nPrimaryID != 0 && OfferData.nPrimaryID == nPrimaryID && OfferData.nOperation == ofopCreation {
					OfferData.nOperation = ofopEdit
				}
				OfferData.nPrimaryID = nPrimaryID

				// Verifica numero secundario da oferta e faz mascaramento
				OfferData.nSecondaryID = maskIntToInt(getOfferSecondaryIDFromFile(arrRecord, c_nSecondaryIDIndex), a_TickerData.MaskDataInfo.hshMaskSecondaryID, &a_TickerData.MaskDataInfo.nCurrentSecondaryD)
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
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to read the records : "+err.Error())
		}
	} else {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to open the file : "+err.Error())
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
	logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer operation type : "+a_arrRecord[a_nIndex])
	return ofopUnknown
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid timestamp : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid loaded date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem somente a data da referencia
	dtReferenceDate, err = validateDateString(fmt.Sprintf(c_strCustomDateFormat, a_dtReferenceTime.Year(), a_dtReferenceTime.Month(), a_dtReferenceTime.Day()))
	if err != nil {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid reference date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem somente a data de geracao
	dtTime, err = validateDateString(fmt.Sprintf(c_strCustomDateFormat, a_dtTime.Year(), a_dtTime.Month(), a_dtTime.Day()))
	if err != nil {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid date : "+err.Error())
		return a_dtLoadedTime
	}
	// Obtem o resultado da subtracao entre a data lida e a referencia
	dtTime = dtTime.Add(dtLoadedDate.Sub(dtReferenceDate))
	// Obtem a nova data com valor normalizado
	dtTime, err = validateTimestampString(fmt.Sprintf(c_strCustomTimestampFormat, dtTime.Year(), dtTime.Month(), dtTime.Day(), a_dtLoadedTime.Hour(), a_dtLoadedTime.Minute(), a_dtLoadedTime.Second(), a_dtLoadedTime.Nanosecond()))
	if err != nil {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid timestamp : "+err.Error())
	}
	return dtTime
}

func maskIntToInt(a_nData int, a_hshIntData map[int]int, a_nCurrentInt *int) int {
	var (
		nMaskIntData int
		bKeyExists   bool
	)
	// Mascara int e concatena atual
	nMaskIntData, bKeyExists = a_hshIntData[a_nData]
	if !bKeyExists {
		*a_nCurrentInt = *a_nCurrentInt + 1
		nMaskIntData = *a_nCurrentInt

		a_hshIntData[a_nData] = nMaskIntData
	}
	return nMaskIntData
}

func maskStringToIntString(a_strData string, a_hshStringData map[string]int, a_nCurrentInt *int) string {
	var (
		nMaskIntData int
		bKeyExists   bool
	)
	// Mascara string e concatena atual
	nMaskIntData, bKeyExists = a_hshStringData[a_strData]
	if !bKeyExists {
		*a_nCurrentInt = *a_nCurrentInt + 1
		nMaskIntData = *a_nCurrentInt

		a_hshStringData[a_strData] = nMaskIntData
	}
	return strconv.Itoa(nMaskIntData)
}

func getOfferGenerationIDFromFile(a_arrRecord []string, a_nIndex int) int {
	const (
		c_strMethodName = "reader.getOfferGenerationIDFromFile"
	)
	var (
		err                error
		nOfferGenerationID int
	)
	nOfferGenerationID, err = validateIntString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer generation ID : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid trade ID : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer primary ID : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer secondary ID : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid trade quantity : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid price : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid current quantity : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid total quantity : "+err.Error())
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
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer size average : "+err.Error())
	}
	return sAvgOfferSize
}

func getSDOfferSizeFromFile(a_arrRecord []string, a_nIndex int) float64 {
	const (
		c_strMethodName = "reader.getSDOfferSizeFromFile"
	)
	var (
		err          error
		sSDOfferSize float64
	)
	sSDOfferSize, err = validateFloatString(a_arrRecord[a_nIndex])
	if err != nil {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Invalid offer size sd : a_nIndex="+strconv.Itoa(a_nIndex)+" : "+err.Error())
		return 0
	}
	return sSDOfferSize
}
