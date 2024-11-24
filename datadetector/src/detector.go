package src

import (
	logger "marketmanipulationdetector/logger/src"
	"runtime"
	"strconv"
	"sync"
)

var (
	m_LogInfo logger.LogInfoType
)

func Start() {
	beginLogInfo()
	startMenu()
	endLogInfo()
}

func beginLogInfo() {
	const (
		c_strMethodName = "detector.beginLogInfo"
	)
	var (
		err error
	)
	m_LogInfo, err = logger.StartAppLog(getLogsPath())
	if err != nil {
		panic("log folder can not be created")
	}

	m_LogInfo, err = logger.CreateLog(m_LogInfo, "Main")
	if err != nil {
		panic("log file can not be created : Main")
	}
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	m_LogInfo, err = logger.CreateLog(m_LogInfo, "Manipulation-Spoofing")
	if err != nil {
		panic("log file can not be created : Manipulation-Spoofing")
	}
	logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Begin")

	m_LogInfo, err = logger.CreateLog(m_LogInfo, "Manipulation-Layering")
	if err != nil {
		panic("log file can not be created : Manipulation-Layering")
	}
	logger.Log(m_LogInfo, "Manipulation-Layering", c_strMethodName, "Begin")
}

func endLogInfo() {
	const (
		c_strMethodName = "detector.endLogInfo"
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
	logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "End")
	logger.Log(m_LogInfo, "Manipulation-Layering", c_strMethodName, "End")
}

func startMenu() {
	const (
		c_strMethodName = "detector.startMenu"
	)
	var (
		nOption int
	)

	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	for {
		printMainMenuOptions("Main")
		nOption = getIntegerFromInput("Main")

		if validateMainMenuOption("Main", nOption) {
			if nOption == 1 {
				startTradeRunForUniqueTicker()
			} else if nOption == 2 {
				startTradeRunForAllTickers(false)
			} else if nOption == 3 {
				startTradeRunForAllTickers(true)
			} else if nOption == 4 || nOption == 5 || nOption == 6 {
				logger.Log(m_LogInfo, "Main", c_strMethodName, "This option is under analysis if will be implemented")
			} else if nOption == 0 {
				break
			}
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func startTradeRunForUniqueTicker() {
	const (
		c_strMethodName = "detector.startTradeRunForUniqueTicker"
	)
	var (
		err          error
		FilesInfo    FilesInfoType
		TradeRunInfo TradeRunInfoType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	TradeRunInfo, err = readTradeRunInput("Main")
	if err == nil {
		logger.Log(m_LogInfo, "Main", c_strMethodName, "strTickerName="+TradeRunInfo.strTickerName+" : dtTickerDate="+TradeRunInfo.dtTickerDate.String())

		FilesInfo, err = getUniqueTickerFiles(TradeRunInfo)
		// Verifica se arquivos (compra e venda) existem conforme ticker e data informado
		if err == nil {
			// Inicia enriquecimento
			runUniqueTicker(false, FilesInfo, nil)
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Ticker file not found")
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func startTradeRunForAllTickers(a_bParallelRun bool) {
	const (
		c_strMethodName = "detector.startTradeRunForAllTickers"
	)
	var (
		err               error
		InfoForAllTickers InfoForAllTickersType
		FilesInfo         FilesInfoType
		arrFilesInfo      []FilesInfoType
		WaitGroup         sync.WaitGroup
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	// Verifica se deve ser aplicado paralelismo
	if a_bParallelRun {
		InfoForAllTickers, err = readInputRunForAllTickers("Main")
		if err == nil {
			logger.Log(m_LogInfo, "Main", c_strMethodName, "nProcessors="+strconv.Itoa(InfoForAllTickers.nProcessors))
			// Seta o numero de processadores a ser utilizado (worker pool)
			runtime.GOMAXPROCS(InfoForAllTickers.nProcessors)

			// Verifica se arquivos (compra e venda) existem para cada ticker conforme data informada
			arrFilesInfo = getAllTickersFiles()

			// Seta o numero de goroutines a serem executadas
			WaitGroup.Add(len(arrFilesInfo))
			logger.Log(m_LogInfo, "Main", c_strMethodName, "Added numbers of routines to be executed : arrFilesInfo="+strconv.Itoa(len(arrFilesInfo)))

			if len(arrFilesInfo) > 0 {
				// Itera sobre tickers disponiveis e processa cada um
				for _, FilesInfo = range arrFilesInfo {
					// Inicia enriquecimento em paralelo com goroutines
					go runUniqueTicker(a_bParallelRun, FilesInfo, &WaitGroup)
				}
			} else {
				logger.LogError(m_LogInfo, "Main", c_strMethodName, "Any ticker files not found")
			}

			// Espera as goroutines finalizarem
			WaitGroup.Wait()
		}
	} else {
		// Verifica se arquivos (compra e venda) existem para cada ticker conforme data informada
		arrFilesInfo = getAllTickersFiles()

		if len(arrFilesInfo) > 0 {
			// Itera sobre tickers disponiveis e processa cada um
			for _, FilesInfo = range arrFilesInfo {
				// Inicia enriquecimento
				runUniqueTicker(a_bParallelRun, FilesInfo, nil)
			}
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Any ticker files not found")
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func runUniqueTicker(a_bParallelRun bool, a_FilesInfo FilesInfoType, a_WaitGroup *sync.WaitGroup) {
	processTradeData(a_FilesInfo)

	// Sinaliza o WaitGroup que a routine finalizou
	if a_bParallelRun {
		defer a_WaitGroup.Done()
	}
}

func processTradeData(a_FilesInfo FilesInfoType) {
	const (
		c_strMethodName = "detector.processTradeData"
	)
	var (
		TickerData TickerDataType
		DataInfo   DataInfoType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin : strTicker="+a_FilesInfo.TradeRunInfo.strTickerName)

	// 1 - Carrega os dados a partir dos arquivos e armazena tudo em memoria
	TickerData = loadTickerData(a_FilesInfo)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Ticker data loaded successfully : strTicker="+a_FilesInfo.TradeRunInfo.strTickerName+" : "+getTickerData(TickerData))

	// 2 - Inicia o processamento dos dados (um por um)
	processEvents(&TickerData, &DataInfo)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Ticker events processed successfully : strTicker="+a_FilesInfo.TradeRunInfo.strTickerName)

	// 3 - Exporta resultados da detecção
	exportResults(&TickerData)

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End : strTicker="+a_FilesInfo.TradeRunInfo.strTickerName)
}
