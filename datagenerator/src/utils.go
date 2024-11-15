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

func getReferencePath() string {
	return getDataPath() + c_strReferenceFolder
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
	strOptions += "\t1 - Generate unique offers book (buy and sell data)\n"

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

func getIntegerFromInput() int {
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

	logger.Log(m_strLogFile, c_strMethodName, "Read integer successfully : nResult="+strconv.Itoa(nResult))

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

func validateGenerationInput(a_strTickerName, a_strTickerDate string) (GenerationInfoType, error) {
	const (
		c_strMethodName = "utils.validateGenerationInput"
	)
	var (
		err          error
		dtTickerDate time.Time
	)

	// Valida ticker informado no terminal
	if a_strTickerName == "" || strings.Contains(a_strTickerName, " ") {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker name")
		return GenerationInfoType{}, errors.New("ticker name validation failure")
	}

	// Valida data informada no terminal e converte para um tipo data
	dtTickerDate, err = validateDateString(a_strTickerDate)
	if err != nil {
		logger.LogError(m_strLogFile, c_strMethodName, "Invalid ticker date : "+err.Error())
		return GenerationInfoType{}, errors.New("ticker date validation failure")
	}

	return GenerationInfoType{
		strTickerName: a_strTickerName,
		dtTickerDate:  dtTickerDate,
	}, nil
}

func readGenerationInput() (GenerationInfoType, error) {
	const (
		c_strMethodName = "utils.readGenerationInput"
	)
	var (
		strTickerName string
		strTickerDate string
	)

	logger.Log(m_strLogFile, c_strMethodName, "Write the ticker name on terminal")
	strTickerName = getStringFromInput()

	logger.Log(m_strLogFile, c_strMethodName, "Write the trade date on terminal (format yyyy-mm-dd)")
	strTickerDate = getStringFromInput()

	return validateGenerationInput(strTickerName, strTickerDate)
}

func checkFileExists(a_strFullPath string) bool {
	var (
		err error
	)
	_, err = os.Stat(a_strFullPath)
	return err == nil
}
