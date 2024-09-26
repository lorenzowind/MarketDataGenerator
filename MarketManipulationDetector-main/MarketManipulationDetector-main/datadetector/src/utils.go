package src

import (
	"bufio"
	"container/list"
	"errors"
	logger "marketmanipulationdetector/logger/src"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func getDataPath() string {
	var (
		strFile     string
		strRootPath string
	)

	_, strFile, _, _ = runtime.Caller(0)

	strRootPath = filepath.Join(filepath.Dir(strFile), "..")

	return strRootPath + c_strDataFolder
}

func getLogsPath() string {
	return getDataPath() + c_strLogsFolder
}

func getInputPath() string {
	return getDataPath() + c_strInputFolder
}

func printMainMenuOptions() {
	const (
		c_strMethodName = "utils.printMainMenuOptions"
	)
	var (
		strOptions string
	)

	strOptions = "\n\n"
	strOptions += "\t0 - Exit\n"
	strOptions += "\t1 - Run unique sequentially\n"
	strOptions += "\t2 - Run all sequentially\n"
	strOptions += "\t3 - Run all with parallelism between tickers\n"
	strOptions += "\t4 - [IN ANALYSIS] Run unique with parallelism between blocks\n"
	strOptions += "\t5 - [IN ANALYSIS] Run all with parallelism between blocks\n"
	strOptions += "\t6 - [IN ANALYSIS] Run all with full parallelism\n"

	logger.Log(m_strLogFile, c_strMethodName, strOptions)
	logger.Log(m_strLogFile, c_strMethodName, "Write an option on terminal")
}

func validateMainMenuOption(a_nOption int) bool {
	const (
		c_strMethodName = "utils.validateMainMenuOption"
	)
	if a_nOption < 0 && a_nOption > 6 {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid option")
		return false
	}

	logger.Log(m_strLogFile, c_strMethodName, "Valid option")
	return true
}

func getOption() int {
	const (
		c_strMethodName = "utils.getOption"
	)
	var (
		nResult     int
		strRead     string
		err         error
		InputReader *bufio.Reader
	)
	InputReader = bufio.NewReader(os.Stdin)

	// Obtem opcao escrita no terminal
	strRead, err = InputReader.ReadString('\n')
	if err != nil {
		logger.LogException(m_strLogFile, c_strMethodName, err.Error())
		return -1
	}

	// Remove o \n do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\n")

	// Remove o \r do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\r")

	// Converte opcao lida do terminal
	nResult, err = strconv.Atoi(strRead)
	if err != nil {
		logger.LogException(m_strLogFile, c_strMethodName, err.Error())
		return -1
	}

	logger.Log(m_strLogFile, c_strMethodName, "Read option successfully : nResult="+strconv.Itoa(nResult))

	return nResult
}

func getStringFromInput() string {
	const (
		c_strMethodName = "utils.getStringFromInput"
	)
	var (
		strRead     string
		err         error
		InputReader *bufio.Reader
	)
	InputReader = bufio.NewReader(os.Stdin)

	// Obtem string escrita no terminal
	strRead, err = InputReader.ReadString('\n')
	if err != nil {
		logger.LogException(m_strLogFile, c_strMethodName, err.Error())
		return ""
	}

	// Remove o \n do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\n")

	logger.Log(m_strLogFile, c_strMethodName, "Read string successfully : strRead="+strRead)

	return strRead
}

func readTradeRunInput(a_bReadTickerName bool) (TradeRunInfoType, error) {
	const (
		c_strMethodName = "utils.readTradeRunInput"
	)
	var (
		strTickerName string
		strTickerDate string
	)

	if a_bReadTickerName {
		logger.Log(m_strLogFile, c_strMethodName, "Write the ticker name on terminal")
		strTickerName = getStringFromInput()
	} else {
		strTickerName = ""
	}

	logger.Log(m_strLogFile, c_strMethodName, "Write the trade date on terminal (format yyyy-mm-dd)")
	strTickerDate = getStringFromInput()

	return validateTradeRunInput(strTickerName, strTickerDate, a_bReadTickerName)
}

func validateTradeRunInput(a_strTickerName, a_strTickerDate string, a_bReadTickerName bool) (TradeRunInfoType, error) {
	const (
		c_strMethodName = "utils.validateTradeRunInput"
	)
	var (
		err          error
		dtTickerDate time.Time
	)

	// Valida ticker informado no terminal
	if a_bReadTickerName {
		if a_strTickerName == "" || strings.Contains(a_strTickerName, " ") {
			logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker name")
			return TradeRunInfoType{}, errors.New("ticker name validation failure")
		}
	}

	// Valida data informada no terminal e converte para um tipo data
	dtTickerDate, err = validateDateString(a_strTickerDate)
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker date : "+err.Error())
		return TradeRunInfoType{}, errors.New("ticker date validation failure")
	}

	return TradeRunInfoType{
		strTickerName: a_strTickerName,
		dtTickerDate:  dtTickerDate,
	}, nil
}

func validateDateString(a_strDate string) (time.Time, error) {
	var (
		err    error
		dtDate time.Time
	)
	if len(a_strDate) > len(time.DateOnly) {
		a_strDate = a_strDate[:len(time.DateOnly)]
	}

	dtDate, err = time.Parse(time.DateOnly, a_strDate)
	return dtDate, err
}

func validateTimestampString(a_strTimestamp string) (time.Time, error) {
	var (
		err    error
		dtDate time.Time
	)
	if len(a_strTimestamp) > len(c_strCustomTimestampLayout) {
		a_strTimestamp = a_strTimestamp[:len(c_strCustomTimestampLayout)]
	}

	if strings.Contains(a_strTimestamp, "T") {
		dtDate, err = time.Parse(c_strCustomTimestampLayout, a_strTimestamp)
	} else {
		dtDate, err = time.Parse(c_strCustomTimestampLayout2, a_strTimestamp)
	}
	return dtDate, err
}

func validateIntString(a_strValue string) (int, error) {
	var (
		err    error
		nValue int
	)
	nValue, err = strconv.Atoi(a_strValue)
	return nValue, err
}

func validateFloatString(a_strValue string) (float64, error) {
	var (
		err    error
		sValue float64
	)
	a_strValue = strings.Replace(a_strValue, ",", ".", -1)
	sValue, err = strconv.ParseFloat(a_strValue, 64)
	return sValue, err
}

//lint:ignore U1000 Ignore unused function
func checkIfHasSameDate(a_dtLeft, a_dtRight time.Time) bool {
	return a_dtLeft.Format(time.DateOnly) == a_dtRight.Format(time.DateOnly)
}

func checkIfContains(a_strItem string, a_arrList []FilesInfoType) bool {
	var (
		FilesInfo FilesInfoType
	)
	for _, FilesInfo = range a_arrList {
		if FilesInfo.TradeRunInfo.strTickerName == a_strItem {
			return true
		}
	}
	return false
}

func checkFileExists(a_strFullPath string) bool {
	var (
		err error
	)
	_, err = os.Stat(a_strFullPath)
	return err == nil
}

//lint:ignore U1000 Ignore unused function
func printListTrades(a_lstData list.List) {
	const (
		c_strMethodName = "utils.printListTrades"
	)
	var (
		TradeData TradeDataType
		Temp      *list.Element
	)
	if a_lstData.Front() == nil {
		logger.Log(m_strLogFile, c_strMethodName, "List of trades is empty")
	} else {
		Temp = a_lstData.Front()
		// Itera sobre cada item da lista encadeada
		for Temp != nil {
			TradeData = Temp.Value.(TradeDataType)
			// Loga os dados de negocio
			printTradeData(TradeData)
			// Obtem o proximo item
			Temp = Temp.Next()
		}
	}
}

func printTradeData(a_TradeData TradeDataType) {
	const (
		c_strMethodName = "utils.printTradeData"
	)
	var (
		strResult string
	)
	strResult = "chOperation=" + string(a_TradeData.chOperation)
	strResult = strResult + " : dtTime=" + a_TradeData.dtTime.String()
	strResult = strResult + " : nID=" + strconv.Itoa(a_TradeData.nID)
	strResult = strResult + " : nOfferGenerationID=" + strconv.Itoa(a_TradeData.nOfferGenerationID)
	strResult = strResult + " : nOfferPrimaryID=" + strconv.Itoa(a_TradeData.nOfferPrimaryID)
	strResult = strResult + " : nOfferSecondaryID=" + strconv.Itoa(a_TradeData.nOfferSecondaryID)
	strResult = strResult + " : strAccount=" + a_TradeData.strAccount
	strResult = strResult + " : nQuantity=" + strconv.Itoa(a_TradeData.nQuantity)
	strResult = strResult + " : sPrice=" + strconv.FormatFloat(a_TradeData.sPrice, 'f', -1, 64)

	logger.Log(m_strLogFile, c_strMethodName, strResult)
}

//lint:ignore U1000 Ignore unused function
func printListOffers(a_lstData list.List) {
	const (
		c_strMethodName = "utils.printListOffers"
	)
	var (
		OfferData OfferDataType
		Temp      *list.Element
	)
	if a_lstData.Front() == nil {
		logger.Log(m_strLogFile, c_strMethodName, "List of offers is empty")
	} else {
		Temp = a_lstData.Front()
		// Itera sobre cada item da lista encadeada
		for Temp != nil {
			OfferData = Temp.Value.(OfferDataType)
			// Loga os dados da oferta
			printOfferData(OfferData)
			// Obtem o proximo item
			Temp = Temp.Next()
		}
	}
}

func printOfferData(a_OfferData OfferDataType) {
	const (
		c_strMethodName = "utils.printOfferData"
	)
	var (
		strResult string
	)
	strResult = "chOperation=" + string(a_OfferData.chOperation)
	strResult = strResult + " : dtTime=" + a_OfferData.dtTime.String()
	strResult = strResult + " : nGenerationID=" + strconv.Itoa(a_OfferData.nGenerationID)
	strResult = strResult + " : nPrimaryID=" + strconv.Itoa(a_OfferData.nPrimaryID)
	strResult = strResult + " : nSecondaryID=" + strconv.Itoa(a_OfferData.nSecondaryID)
	strResult = strResult + " : nTradeID=" + strconv.Itoa(a_OfferData.nTradeID)
	strResult = strResult + " : strAccount=" + a_OfferData.strAccount
	strResult = strResult + " : nCurrentQuantity=" + strconv.Itoa(a_OfferData.nCurrentQuantity)
	strResult = strResult + " : nTradeQuantity=" + strconv.Itoa(a_OfferData.nTradeQuantity)
	strResult = strResult + " : nTotalQuantity=" + strconv.Itoa(a_OfferData.nTotalQuantity)
	strResult = strResult + " : sPrice=" + strconv.FormatFloat(a_OfferData.sPrice, 'f', -1, 64)

	logger.Log(m_strLogFile, c_strMethodName, strResult)
}

//lint:ignore U1000 Ignore unused function
func printListBookOffers(a_lstData list.List) {
	const (
		c_strMethodName = "utils.printListBookOffers"
	)
	var (
		BookOffer BookOfferType
		Temp      *list.Element
	)
	if a_lstData.Front() == nil {
		logger.Log(m_strLogFile, c_strMethodName, "List of book offers is empty")
	} else {
		Temp = a_lstData.Front()
		// Itera sobre cada item da lista encadeada
		for Temp != nil {
			BookOffer = Temp.Value.(BookOfferType)
			// Loga os dados da oferta
			printBookOffer(BookOffer)
			// Obtem o proximo item
			Temp = Temp.Next()
		}
	}
}

func printBookOffer(a_BookOffer BookOfferType) {
	const (
		c_strMethodName = "utils.printBookOffer"
	)
	var (
		strResult string
	)
	strResult = "nGenerationID=" + strconv.Itoa(a_BookOffer.nGenerationID)
	strResult = strResult + " : nQuantity=" + strconv.Itoa(a_BookOffer.nQuantity)
	strResult = strResult + " : nSecondaryID=" + strconv.Itoa(a_BookOffer.nSecondaryID)
	strResult = strResult + " : strAccount=" + a_BookOffer.strAccount
	strResult = strResult + " : sPrice=" + strconv.FormatFloat(a_BookOffer.sPrice, 'f', -1, 64)

	logger.Log(m_strLogFile, c_strMethodName, strResult)
}

//lint:ignore U1000 Ignore unused function
func printListBookPrices(a_lstData list.List) {
	const (
		c_strMethodName = "utils.printListBookPrices"
	)
	var (
		BookPrice BookPriceType
		Temp      *list.Element
	)
	if a_lstData.Front() == nil {
		logger.Log(m_strLogFile, c_strMethodName, "List of book prices is empty")
	} else {
		Temp = a_lstData.Front()
		// Itera sobre cada item da lista encadeada
		for Temp != nil {
			BookPrice = Temp.Value.(BookPriceType)
			// Loga os dados da oferta
			printBookPrice(BookPrice)
			// Obtem o proximo item
			Temp = Temp.Next()
		}
	}
}

func printBookPrice(a_BookPrice BookPriceType) {
	const (
		c_strMethodName = "utils.printBookOffer"
	)
	var (
		strResult string
	)
	strResult = "sPrice=" + strconv.FormatFloat(a_BookPrice.sPrice, 'f', -1, 64)
	strResult = strResult + " : nCount=" + strconv.Itoa(a_BookPrice.nCount)
	strResult = strResult + " : nQuantity=" + strconv.Itoa(a_BookPrice.nQuantity)

	logger.Log(m_strLogFile, c_strMethodName, strResult)
}
