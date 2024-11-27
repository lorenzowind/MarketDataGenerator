package src

import (
	"container/list"
	"encoding/csv"
	logger "marketmanipulationdetector/logger/src"
	"os"
	"strconv"
)

func saveOffersBook(a_TickerData TickerDataType) {
	const (
		c_strMethodName = "writer.saveOffersBook"
	)
	// Salva dados do arquivo de compra
	if a_TickerData.FilesInfo.strBuyPath != "" {
		if saveOfferDataFromFile(a_TickerData.FilesInfo.strBuyPath, &a_TickerData, true) {
			logger.Log(m_LogInfo, "Main", c_strMethodName, "Buy file generated successfully : strBuyPath="+a_TickerData.FilesInfo.strBuyPath)
		}
	}

	// Salva dados do arquivo de venda
	if a_TickerData.FilesInfo.strSellPath != "" {
		if saveOfferDataFromFile(a_TickerData.FilesInfo.strSellPath, &a_TickerData, false) {
			logger.Log(m_LogInfo, "Main", c_strMethodName, "Sell file generated successfully : strSellPath="+a_TickerData.FilesInfo.strSellPath)
		}
	}

	// Carrega dados de benchmark
	if a_TickerData.FilesInfo.strBenchmarkPath != "" {
		if saveBenchmarkFromFile(a_TickerData.FilesInfo.strBenchmarkPath, &a_TickerData) {
			logger.Log(m_LogInfo, "Main", c_strMethodName, "Benchmark file generated successfully : strBenchmarkPath="+a_TickerData.FilesInfo.strBenchmarkPath)
		}
	}
}

func saveOfferDataFromFile(a_strPath string, a_TickerData *TickerDataType, bBuy bool) bool {
	const (
		c_strMethodName            = "writer.saveOfferDataFromFile"
		c_strOperationHeader       = "cod_evento_oferta"
		c_strTickerHeader          = "cod_simbolo_instrumento_negociacao"
		c_strTimeHeader            = "dthr_inclusao_oferta"
		c_strGenerationIDHeader    = "num_geracao_oferta"
		c_strAccountHeader         = "num_identificacao_conta"
		c_strTradeIDHeader         = "num_negocio"
		c_strPrimaryIDHeader       = "num_sequencia_oferta"
		c_strSecondaryIDHeader     = "num_sequencia_oferta_secundaria"
		c_strCurrentQuantityHeader = "qte_divulgada_oferta"
		c_strTradeQuantityHeader   = "qte_negociada"
		c_strTotalQuantityHeader   = "qte_total_oferta"
		c_strPriceHeader           = "val_preco_oferta"
	)
	var (
		err        error
		bFullWrite bool
		lstData    *list.List
		file       *os.File
		writer     *csv.Writer
		OfferData  OfferDataType
		Temp       *list.Element
	)
	bFullWrite = true

	file, err = os.OpenFile(a_strPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err == nil {
		defer file.Close()

		writer = csv.NewWriter(file)
		writer.Comma = ';'

		defer writer.Flush()

		// Escreve os dados do cabecalho (11 colunas)
		err = writer.Write([]string{
			c_strOperationHeader,
			c_strTickerHeader,
			c_strTimeHeader,
			c_strGenerationIDHeader,
			c_strAccountHeader,
			c_strTradeIDHeader,
			c_strPrimaryIDHeader,
			c_strSecondaryIDHeader,
			c_strCurrentQuantityHeader,
			c_strTradeQuantityHeader,
			c_strTotalQuantityHeader,
			c_strPriceHeader,
		})

		if err == nil {
			if bBuy {
				lstData = &a_TickerData.lstBuy
			} else {
				lstData = &a_TickerData.lstSell
			}

			if lstData.Front() != nil {
				Temp = lstData.Front()
				// Itera sobre cada item da lista encadeada
				for Temp != nil {
					OfferData = Temp.Value.(OfferDataType)
					// Escreve os dados da oferta
					err = writer.Write([]string{
						string(OfferData.nOperation),                        // cod_evento_oferta
						a_TickerData.FilesInfo.GenerationInfo.strTickerName, // cod_simbolo_instrumento_negociacao (geracao - normalizado)
						getTimeAsCustomTimestamp(OfferData.dtTime),          // dthr_inclusao_oferta (geracao - normalizado)
						strconv.Itoa(OfferData.nGenerationID),               // num_geracao_oferta
						OfferData.strAccount,                                // num_identificacao_conta
						strconv.Itoa(OfferData.nTradeID),                    // num_negocio
						strconv.Itoa(OfferData.nPrimaryID),                  // num_sequencia_oferta
						strconv.Itoa(OfferData.nSecondaryID),                // num_sequencia_oferta_secundaria
						strconv.Itoa(OfferData.nCurrentQuantity),            // qte_divulgada_oferta
						strconv.Itoa(OfferData.nTradeQuantity),              // qte_negociada
						strconv.Itoa(OfferData.nTotalQuantity),              // qte_total_oferta
						strconv.FormatFloat(OfferData.sPrice, 'f', -1, 64),  // val_preco_oferta
					})
					if err != nil {
						logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to write record on the file : "+err.Error())
						bFullWrite = false
						break
					}
					// Obtem o proximo item
					Temp = Temp.Next()
				}
			}
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to write header on the file : "+err.Error())
			bFullWrite = false
		}
	} else {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to open the file : "+err.Error())
		bFullWrite = false
	}

	return bFullWrite
}

func saveBenchmarkFromFile(a_strPath string, a_TickerData *TickerDataType) bool {
	const (
		c_strMethodName               = "writer.saveBenchmarkFromFile"
		c_strTickerHeader             = "cod_simbolo_instrumento_negociacao"
		c_strAvgTradeIntervalHeader   = "media_intervalo_negocios"
		c_strAvgOfferSizeHeader       = "media_qtd_ofertas"
		c_strBiggerSDOfferSizeHeader  = "menor_dp_qtd_ofertas"
		c_strSmallerSDOfferSizeHeader = "maior_dp_qtd_ofertas"
	)
	var (
		err        error
		bFullWrite bool
		bAppend    bool
		file       *os.File
		writer     *csv.Writer
	)
	bFullWrite = true

	_, err = os.Stat(a_strPath)
	if err == nil {
		file, err = os.OpenFile(a_strPath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		bAppend = true
	} else {
		file, err = os.OpenFile(a_strPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		bAppend = false
	}

	if err == nil {
		defer file.Close()

		writer = csv.NewWriter(file)
		writer.Comma = ';'

		defer writer.Flush()

		if !bAppend {
			// Escreve os dados do cabecalho (5 colunas)
			err = writer.Write([]string{
				c_strTickerHeader,
				c_strAvgTradeIntervalHeader,
				c_strAvgOfferSizeHeader,
				c_strBiggerSDOfferSizeHeader,
				c_strSmallerSDOfferSizeHeader,
			})
		}

		if err == nil || bAppend {
			// Escreve os dados da oferta
			err = writer.Write([]string{
				a_TickerData.FilesInfo.GenerationInfo.strTickerName,                              // cod_simbolo_instrumento_negociacao (geracao - normalizado)
				getTimeAsCustomDuration(a_TickerData.BenchmarkData.dtAvgTradeInterval),           // media_intervalo_negocios
				strconv.FormatFloat(a_TickerData.BenchmarkData.sAvgOfferSize, 'f', -1, 64),       // media_qtd_ofertas
				strconv.FormatFloat(a_TickerData.BenchmarkData.sSmallerSDOfferSize, 'f', -1, 64), // menor_dp_qtd_ofertas
				strconv.FormatFloat(a_TickerData.BenchmarkData.sBiggerSDOfferSize, 'f', -1, 64),  // maior_dp_qtd_ofertas
			})
			if err != nil {
				logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to write record on the file : "+err.Error())
				bFullWrite = false
			}
		} else {
			logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to write header on the file : "+err.Error())
			bFullWrite = false
		}
	} else {
		logger.LogError(m_LogInfo, "Main", c_strMethodName, "Fail to open the file : "+err.Error())
		bFullWrite = false
	}

	return bFullWrite
}
