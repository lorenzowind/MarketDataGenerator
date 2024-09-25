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
		bBuyEnd   bool
		bSellEnd  bool
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	NextBuy = a_TickerData.lstBuy.Front()
	NextSell = a_TickerData.lstSell.Front()

	bBuyEnd = false
	bSellEnd = false
	for {
		bProcess = true
		if NextBuy != nil && NextSell != nil && !bBuyEnd && !bSellEnd {
			BuyData = NextBuy.Value.(OfferDataType)
			SellData = NextSell.Value.(OfferDataType)
			// Verifica ID de geracao para verificar qual evento ocorreu primeiro
			if BuyData.nGenerationID < SellData.nGenerationID {
				OfferData = BuyData
				// Obtem o proximo evento de oferta de compra
				NextBuy = NextBuy.Next()
				if NextBuy == nil {
					bBuyEnd = true
				}
			} else if SellData.nGenerationID < BuyData.nGenerationID {
				OfferData = SellData
				// Obtem o proximo evento de oferta de venda
				NextSell = NextSell.Next()
				if NextSell == nil {
					bSellEnd = true
				}
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "Generation ID is equal for buy and sell offer : nGenerationID="+strconv.Itoa(BuyData.nGenerationID))
			}
		} else if NextBuy != nil && !bBuyEnd {
			BuyData = NextBuy.Value.(OfferDataType)
			OfferData = BuyData
			// Obtem o proximo evento de oferta de compra
			NextBuy = NextBuy.Next()
			if NextBuy == nil {
				bBuyEnd = true
			}
		} else if NextSell != nil && !bSellEnd {
			SellData = NextSell.Value.(OfferDataType)
			OfferData = SellData
			// Obtem o proximo evento de oferta de venda
			NextSell = NextSell.Next()
			if NextSell == nil {
				bSellEnd = true
			}
		} else {
			bProcess = false
			logger.LogError(m_strLogFile, c_strMethodName, "NextBuy and NextSell are nil")
		}
		// Processa evento da oferta de compra ou venda
		if bProcess {
			processOffer(OfferData)
		}
		// Condicao de parada -> os eventos foram processados
		if bBuyEnd && bSellEnd {
			break
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func processOffer(a_OfferData OfferDataType) {
	const (
		c_strMethodName = "manager.processOffer"
	)
	switch a_OfferData.chOperation {
	case ofopCreation:
		processEventCreation(a_OfferData)
	case ofopCancel:
		processEventCancel(a_OfferData)
	case ofopEdit:
		processEventEdit(a_OfferData)
	case ofopExpired:
		processEventExpired(a_OfferData)
	case ofopReafirmed:
		processEventReafirmed(a_OfferData)
	case ofopTrade:
		processEventTrade(a_OfferData)
	case ofopUnknown:
		logger.LogError(m_strLogFile, c_strMethodName, "Unknown offer operation : chOperation="+string(a_OfferData.chOperation))
	default:
		logger.LogError(m_strLogFile, c_strMethodName, "Default offer operation : chOperation="+string(a_OfferData.chOperation))
	}
}

func processEventCreation(a_OfferData OfferDataType) {

}

func processEventCancel(a_OfferData OfferDataType) {

}

func processEventEdit(a_OfferData OfferDataType) {

}

func processEventExpired(a_OfferData OfferDataType) {

}

func processEventReafirmed(a_OfferData OfferDataType) {

}

func processEventTrade(a_OfferData OfferDataType) {

}

func exportResults(a_TickerData TickerDataType) {

}
