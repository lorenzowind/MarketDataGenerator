package src

import (
	"container/list"
	logger "marketmanipulationdetector/logger/src"
	"strconv"
)

func processEvents(a_TickerData TickerDataType) {
	const (
		c_strMethodName = "manager.processEvents"
	)
	var (
		NextBuy   *list.Element
		NextSell  *list.Element
		BuyData   OfferDataType
		SellData  OfferDataType
		OfferData OfferDataType
		bProcess  bool
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	NextBuy = a_TickerData.lstBuy.Front()
	NextSell = a_TickerData.lstSell.Front()

	for {
		bProcess = true
		if NextBuy != nil && NextSell != nil {
			BuyData = NextBuy.Value.(OfferDataType)
			SellData = NextSell.Value.(OfferDataType)
			// Verifica ID de geracao para verificar qual evento ocorreu primeiro
			if BuyData.nGenerationID < SellData.nGenerationID {
				OfferData = BuyData
				// Obtem o proximo evento de oferta de compra
				NextBuy = NextBuy.Next()
			} else if SellData.nGenerationID < BuyData.nGenerationID {
				OfferData = SellData
				// Obtem o proximo evento de oferta de venda
				NextSell = NextSell.Next()
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "Generation ID is equal for buy and sell offer : nGenerationID="+strconv.Itoa(BuyData.nGenerationID))
			}
		} else if NextBuy != nil {
			BuyData = NextBuy.Value.(OfferDataType)
			OfferData = BuyData
			// Obtem o proximo evento de oferta de compra
			NextBuy = NextBuy.Next()
		} else if NextSell != nil {
			SellData = NextSell.Value.(OfferDataType)
			OfferData = SellData
			// Obtem o proximo evento de oferta de venda
			NextSell = NextSell.Next()
		} else {
			bProcess = false
			logger.LogError(m_strLogFile, c_strMethodName, "NextBuy and NextSell are nil")
		}
		// Processa evento da oferta de compra ou venda
		if bProcess {
			processOffer(OfferData)
		}
		// Condicao de parada -> os eventos foram processados
		if NextBuy == a_TickerData.lstBuy.Front() && NextSell == a_TickerData.lstSell.Front() {
			break
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func processOffer(a_OfferData OfferDataType) {

}

func exportResults(a_TickerData TickerDataType) {

}
