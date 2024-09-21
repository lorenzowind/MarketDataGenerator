package src

import (
	logger "marketmanipulationdetector/logger/src"
	"time"
)

var (
	m_strLogFolder string
	m_strLogFile   string
)

func Start() {
	const (
		c_strMethodName = "detector.Start"
	)
	var (
		err error
	)

	m_strLogFolder, err = logger.StartAppLog(getLogsPath())
	if err != nil {
		panic("log folder can not be created")
	}

	m_strLogFile, err = logger.CreateLog(m_strLogFolder, "Main")
	if err != nil {
		panic("log file can not be created")
	}

	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	startMenu()

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func startMenu() {
	const (
		c_strMethodName = "detector.startMenu"
	)
	var (
		nOption int
	)

	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	for {
		printMainMenuOptions()
		nOption = getOption()

		if validateMainMenuOption(nOption) {
			if nOption == 1 {
				startTradeRunForUniqueTicker(false)
			} else if nOption == 2 {
				startTradeRunForAllTickers(false)
			} else if nOption == 3 {
				startTradeRunForAllTickers(true)
			} else if nOption == 4 || nOption == 5 || nOption == 6 {
				logger.Log(m_strLogFile, c_strMethodName, "This option is under analysis if will be implemented")
			} else if nOption == 0 {
				break
			}
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func startTradeRunForUniqueTicker(a_bParallelRun bool) {
	const (
		c_strMethodName = "detector.startTradeRunForUniqueTicker"
	)
	var (
		err          error
		TradeRunInfo TradeRunInfoType
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	// Verifica se deve ser aplicado paralelismo
	if a_bParallelRun {
		logger.Log(m_strLogFile, c_strMethodName, "Not implemented yet")
	} else {
		TradeRunInfo, err = readTradeRunInput(true)
		if err == nil {
			logger.Log(m_strLogFile, c_strMethodName, "strTickerName="+TradeRunInfo.strTickerName+" : dtTickerDate="+TradeRunInfo.dtTickerDate.String())

			// Verifica se arquivos (compra, venda e negocio) existem conforme ticker e data informado
			if findTradeFiles(TradeRunInfo) {
				// Inicia enriquecimento
				runUniqueTicker(a_bParallelRun, TradeRunInfo)
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "Files not found")
			}
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func startTradeRunForAllTickers(a_bParallelRun bool) {
	const (
		c_strMethodName = "detector.startTradeRunForAllTickers"
	)
	var (
		err          error
		TradeRunInfo TradeRunInfoType
		lstTickers   []string
		dtTickerDate time.Time
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	// Verifica se deve ser aplicado paralelismo
	if a_bParallelRun {
		logger.Log(m_strLogFile, c_strMethodName, "Not implemented yet")
	} else {
		TradeRunInfo, err = readTradeRunInput(false)
		if err == nil {
			dtTickerDate = TradeRunInfo.dtTickerDate
			logger.Log(m_strLogFile, c_strMethodName, "dtTickerDate="+dtTickerDate.String())

			// Verifica se arquivos (compra, venda e negocio) existem para cada ticker conforme data informada
			lstTickers = getValidTickers(TradeRunInfo.dtTickerDate)

			if len(lstTickers) > 0 {
				// Itera sobre tickers disponiveis e processa cada um
				for _, strTicker := range lstTickers {
					TradeRunInfo = TradeRunInfoType{
						strTickerName: strTicker,
						dtTickerDate:  dtTickerDate,
					}

					// Inicia enriquecimento
					runUniqueTicker(a_bParallelRun, TradeRunInfo)
				}
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "No files found")
			}
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func runUniqueTicker(a_bParallelRun bool, a_TradeRunInfo TradeRunInfoType) {
	const (
		c_strMethodName = "detector.runUniqueTicker"
	)
	if a_bParallelRun {
		logger.Log(m_strLogFile, c_strMethodName, "Not implemented yet")
	} else {
		processTradeData(a_TradeRunInfo)
	}
}

func processTradeData(a_TradeRunInfo TradeRunInfoType) {
	const (
		c_strMethodName = "detector.processTradeData"
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	// 1 - Carrega os dados a partir dos arquivos e armazena tudo em memoria
	loadTradeDataFromFile(a_TradeRunInfo)

	// 2 - Inicia o processamento dos dados (um por um)
	processEvents(a_TradeRunInfo)

	// 3 - Exporta resultados da detecção
	exportResults(a_TradeRunInfo)

	logger.Log(m_strLogFile, c_strMethodName, "End")
}
