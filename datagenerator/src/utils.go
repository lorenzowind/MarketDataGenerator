package src

import (
	"bufio"
	logger "marketmanipulationdetector/logger/src"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func getLogsPath() string {
	const (
		c_strDataFolder = "/data"
		c_strLogsFolder = "/logs"
	)
	var (
		strFile     string
		strRootPath string
	)

	_, strFile, _, _ = runtime.Caller(0)

	strRootPath = filepath.Join(filepath.Dir(strFile), "..")

	return strRootPath + c_strDataFolder + c_strLogsFolder
}

func printMainMenuOptions() {
	const (
		c_strMethodName = "utils.printMenuOptions"
	)
	var (
		strOptions string
	)

	strOptions = "\n\n"
	strOptions += "\t0 - Exit\n"
	strOptions += "\t1 - Generate data\n"

	logger.Log(m_strLogFile, c_strMethodName, strOptions)
	logger.Log(m_strLogFile, c_strMethodName, "Write an option on terminal")
}

func validateMainMenuOption(a_nOption int) bool {
	const (
		c_strMethodName = "utils.validateOption"
	)
	if a_nOption != 0 && a_nOption != 1 {
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

	// Converte opcao lida do terminal
	nResult, err = strconv.Atoi(strRead)
	if err != nil {
		logger.LogException(m_strLogFile, c_strMethodName, err.Error())
		return -1
	}

	logger.Log(m_strLogFile, c_strMethodName, "Read option successfully : nResult="+strRead)

	return nResult
}
