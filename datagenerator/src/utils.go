package src

import (
	"bufio"
	"errors"
	"fmt"
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

func getReferencePath() string {
	return getDataPath() + c_strReferenceFolder
}

func printMainMenuOptions(a_strParentLog string) {
	const (
		c_strMethodName = "utils.printMainMenuOptions"
	)
	var (
		strOptions string
	)

	strOptions = "\n\n"
	strOptions += "\t0 - Exit\n"
	strOptions += "\t1 - Generate unique offers book (buy and sell data)\n"

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, strOptions)
	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Write an option on terminal")
}

func validateMainMenuOption(a_strParentLog string, a_nOption int) bool {
	const (
		c_strMethodName = "utils.validateMainMenuOption"
	)
	if a_nOption < 0 && a_nOption > 6 {
		logger.LogError(m_LogInfo, a_strParentLog, c_strMethodName, "Invalid option")
		return false
	}

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Valid option")
	return true
}

func getIntegerFromInput(a_strParentLog string) int {
	const (
		c_strMethodName = "utils.getIntegerFromInput"
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
		logger.LogException(m_LogInfo, a_strParentLog, c_strMethodName, err.Error())
		return -1
	}

	// Remove o \n do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\n")

	// Remove o \r do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\r")

	// Converte opcao lida do terminal
	nResult, err = strconv.Atoi(strRead)
	if err != nil {
		logger.LogException(m_LogInfo, a_strParentLog, c_strMethodName, err.Error())
		return -1
	}

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Read integer successfully : nResult="+strconv.Itoa(nResult))

	return nResult
}

func getStringFromInput(a_strParentLog string) string {
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
		logger.LogException(m_LogInfo, a_strParentLog, c_strMethodName, err.Error())
		return ""
	}

	// Remove o \n do conteudo lido
	strRead = strings.TrimSuffix(strRead, "\n")

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Read string successfully : strRead="+strRead)

	return strRead
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

func validateGenerationInput(a_strParentLog, a_strReferenceTickerName, a_strReferenceTickerDate, a_strTickerName, a_strTickerDate string) (GenerationInfoType, error) {
	const (
		c_strMethodName = "utils.validateGenerationInput"
	)
	var (
		err                   error
		dtReferenceTickerDate time.Time
		dtTickerDate          time.Time
	)
	// Valida ticker de referencia informado no terminal
	if a_strReferenceTickerName == "" || strings.Contains(a_strReferenceTickerName, " ") {
		logger.LogError(m_LogInfo, a_strParentLog, c_strMethodName, "Invalid reference ticker name")
		return GenerationInfoType{}, errors.New("reference ticker name validation failure")
	}
	// Valida data de referencia informada no terminal e converte para um tipo data
	dtReferenceTickerDate, err = validateDateString(a_strReferenceTickerDate)
	if err != nil {
		logger.LogError(m_LogInfo, a_strParentLog, c_strMethodName, "Invalid reference ticker date : "+err.Error())
		return GenerationInfoType{}, errors.New("reference ticker date validation failure")
	}
	// Valida ticker informado no terminal
	if a_strTickerName == "" || strings.Contains(a_strTickerName, " ") {
		logger.LogError(m_LogInfo, a_strParentLog, c_strMethodName, "Invalid ticker name")
		return GenerationInfoType{}, errors.New("ticker name validation failure")
	}
	// Valida data informada no terminal e converte para um tipo data
	dtTickerDate, err = validateDateString(a_strTickerDate)
	if err != nil {
		logger.LogError(m_LogInfo, a_strParentLog, c_strMethodName, "Invalid ticker date : "+err.Error())
		return GenerationInfoType{}, errors.New("ticker date validation failure")
	}

	return GenerationInfoType{
		strTickerName:          a_strTickerName,
		dtTickerDate:           dtTickerDate,
		strReferenceTickerName: a_strReferenceTickerName,
		dtReferenceTickerDate:  dtReferenceTickerDate,
	}, nil
}

func readGenerationInput(a_strParentLog string) (GenerationInfoType, error) {
	const (
		c_strMethodName = "utils.readGenerationInput"
	)
	var (
		strReferenceTickerName string
		strReferenceTickerDate string
		strTickerName          string
		strTickerDate          string
	)

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Write the reference ticker name on terminal")
	strReferenceTickerName = getStringFromInput(a_strParentLog)

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Write the reference trade date on terminal (format yyyy-mm-dd)")
	strReferenceTickerDate = getStringFromInput(a_strParentLog)

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Write the generation ticker name on terminal")
	strTickerName = getStringFromInput(a_strParentLog)

	logger.Log(m_LogInfo, a_strParentLog, c_strMethodName, "Write the generation trade date on terminal (format yyyy-mm-dd)")
	strTickerDate = getStringFromInput(a_strParentLog)

	return validateGenerationInput(a_strParentLog, strReferenceTickerName, strReferenceTickerDate, strTickerName, strTickerDate)
}

func checkFileExists(a_strFullPath string) bool {
	var (
		err error
	)
	_, err = os.Stat(a_strFullPath)
	return err == nil
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
	} else if strings.Contains(a_strTimestamp, " ") {
		dtDate, err = time.Parse(c_strCustomTimestampLayout2, a_strTimestamp)
	} else {
		dtDate, err = time.Parse(c_strCustomTimestampLayout3, a_strTimestamp)
	}

	return dtDate, err
}

//lint:ignore U1000 Ignore unused function
func getOfferData(a_OfferData OfferDataType) string {
	var (
		strResult string
	)
	strResult = "nOperation=" + string(a_OfferData.nOperation)
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

	return strResult
}

func getTickerData(a_TickerData TickerDataType) string {
	var (
		strResult string
	)
	strResult = "Buy=" + strconv.Itoa(a_TickerData.lstBuy.Len())
	strResult = strResult + " : Sell=" + strconv.Itoa(a_TickerData.lstSell.Len())
	strResult = strResult + " : HasBenchmarkData=" + strconv.FormatBool(a_TickerData.BenchmarkData.bHasBenchmarkData)

	// So exibe valores de benchmark caso tenha o encontrado
	if a_TickerData.BenchmarkData.bHasBenchmarkData {
		strResult = strResult + " : AvgTrade=" + a_TickerData.BenchmarkData.dtAvgTradeInterval.String()
		strResult = strResult + " : AvgOfferSize=" + strconv.FormatFloat(a_TickerData.BenchmarkData.sAvgOfferSize, 'f', -1, 64)
		strResult = strResult + " : SmallerSDOfferSize=" + strconv.FormatFloat(a_TickerData.BenchmarkData.sSmallerSDOfferSize, 'f', -1, 64)
		strResult = strResult + " : BiggerSDOfferSize=" + strconv.FormatFloat(a_TickerData.BenchmarkData.sBiggerSDOfferSize, 'f', -1, 64)
	}

	return strResult
}

func getTimeAsCustomTimestamp(a_dtTime time.Time) string {
	return fmt.Sprintf(c_strCustomTimestampFormat, a_dtTime.Year(), a_dtTime.Month(), a_dtTime.Day(), a_dtTime.Hour(), a_dtTime.Minute(), a_dtTime.Second(), a_dtTime.Nanosecond())
}

func getTimeAsCustomDuration(a_dtTime time.Time) string {
	return fmt.Sprintf(c_strCustomDurationFormat, a_dtTime.Hour(), a_dtTime.Minute(), a_dtTime.Second(), a_dtTime.Nanosecond())
}
