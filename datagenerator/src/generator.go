package src

import (
	logger "marketmanipulationdetector/logger/src"
	"time"
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
		c_strMethodName = "generator.beginLogInfo"
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
}

func endLogInfo() {
	const (
		c_strMethodName = "generator.endLogInfo"
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func startMenu() {
	const (
		c_strMethodName = "generator.startMenu"
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
				startGenerationForUniqueOffersBook()
			} else if nOption == 2 {
				startGenerationForAllOffersBook(bfAvgTradeInterval)
			} else if nOption == 0 {
				break
			}
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func startGenerationForUniqueOffersBook() {
	const (
		c_strMethodName = "generator.startGenerationForUniqueOffersBook"
	)
	var (
		err            error
		FilesInfo      FilesInfoType
		GenerationInfo GenerationInfoType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	GenerationInfo, err = readGenerationInput("Main")
	if err == nil {
		logger.Log(m_LogInfo, "Main", c_strMethodName, "strReferenceTickerName="+GenerationInfo.strReferenceTickerName+" : dtReferenceTickerDate="+GenerationInfo.dtReferenceTickerDate.Format(time.DateOnly)+" : strTickerName="+GenerationInfo.strTickerName+" : dtTickerDate="+GenerationInfo.dtTickerDate.Format(time.DateOnly))

		FilesInfo, err = getReferenceOffersBook(GenerationInfo)
		// Verifica se arquivos (compra e venda) de referencia existem conforme ticker e data informado
		if err == nil {
			// Inicia geracao do livro
			generateOffersBook(FilesInfo)
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Ticker file not found")
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func startGenerationForAllOffersBook(a_nBenchmarkFilter BenchmarkFilterType) {
	const (
		c_strMethodName = "generator.startGenerationForAllOffersBook"
	)
	var (
		err               error
		FilesInfo         FilesInfoType
		GenerationRule    GenerationRuleType
		GenerationInfo    GenerationInfoType
		arrGenerationInfo []GenerationInfoType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin")

	GenerationRule, err = readGenerationRule("Main", a_nBenchmarkFilter)
	if err == nil {
		// Loga o filtro de benchmark considerando o intervalo entre trades
		if a_nBenchmarkFilter == bfAvgTradeInterval {
			logger.Log(m_LogInfo, "Main", c_strMethodName, "strTickerNameRule="+GenerationRule.strTickerNameRule+" : dtTickerDate="+GenerationRule.dtTickerDate.Format(time.DateOnly)+" : dtMinAvgTradeInterval="+getTimeAsCustomDuration(GenerationRule.Pattern.dtMinAvgTradeInterval)+" : dtMaxAvgTradeInterval="+getTimeAsCustomDuration(GenerationRule.Pattern.dtMaxAvgTradeInterval))
		}
		// Gera tickers com base no benchmark informado (ja filtrados)
		arrGenerationInfo = generateFilterTickers(GenerationRule)

		if len(arrGenerationInfo) > 0 {
			// Itera sobre tickers disponiveis e processa cada um
			for _, GenerationInfo = range arrGenerationInfo {
				FilesInfo, err = getReferenceOffersBook(GenerationInfo)
				// Verifica se arquivos (compra e venda) de referencia existem conforme ticker e data informado
				if err == nil {
					// Inicia geracao do livro
					generateOffersBook(FilesInfo)
				} else {
					logger.LogError(m_LogInfo, "Main", c_strMethodName, "Ticker file not found")
				}
			}
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Any ticker files not found")
		}
	}

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End")
}

func generateOffersBook(a_FilesInfo FilesInfoType) {
	const (
		c_strMethodName = "generator.generateOffersBook"
	)
	var (
		TickerData TickerDataType
	)
	logger.Log(m_LogInfo, "Main", c_strMethodName, "Begin : strTicker="+a_FilesInfo.GenerationInfo.strTickerName)

	// 1 - Salva nome dos arquivos a serem gerados na pasta input na data e ticker escolhido
	getOffersBook(&a_FilesInfo)

	// 2 - Carrega os dados a partir dos arquivos e armazena tudo em memoria (ja normalizados e mascarados)
	TickerData = loadTickerData(a_FilesInfo)

	// 3 - Salva valores nos arquivos finais
	saveOffersBook(TickerData)

	logger.Log(m_LogInfo, "Main", c_strMethodName, "End : strTicker="+a_FilesInfo.GenerationInfo.strTickerName)
}
