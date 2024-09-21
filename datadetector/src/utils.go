package src

import (
	"bufio"
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
	if len(a_strDate) == len(time.DateOnly) {
		dtDate, err = time.Parse(time.DateOnly, a_strDate[:10])
	} else {
		err = errors.New("ticker date size is invalid")
	}

	return dtDate, err
}

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
