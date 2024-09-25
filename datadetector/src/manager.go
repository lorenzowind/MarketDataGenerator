package src

import (
	"container/list"
	logger "marketmanipulationdetector/logger/src"
	"strconv"
)

func processEvents(a_TickerData TickerDataType, a_DataInfo *DataInfoType) {
	const (
		c_strMethodName = "manager.processEvents"
	)
	var (
		NextBuy   *list.Element
		NextSell  *list.Element
		BuyData   OfferDataType
		SellData  OfferDataType
		OfferData OfferDataType
		EventInfo EventInfoType
	)
	logger.Log(m_strLogFile, c_strMethodName, "Begin")

	NextBuy = a_TickerData.lstBuy.Front()
	NextSell = a_TickerData.lstSell.Front()

	EventInfo.bBuyEventsEnd = false
	EventInfo.bSellEventsEnd = false
	for {
		EventInfo.bProcessEvent = true
		if NextBuy != nil && NextSell != nil && !EventInfo.bBuyEventsEnd && !EventInfo.bSellEventsEnd {
			BuyData = NextBuy.Value.(OfferDataType)
			SellData = NextSell.Value.(OfferDataType)
			// Verifica ID de geracao para verificar qual evento ocorreu primeiro
			if BuyData.nGenerationID < SellData.nGenerationID {
				OfferData = BuyData
				// Especifica que eh um evento de compra
				EventInfo.bBuyEvent = true
				// Obtem o proximo evento de oferta de compra
				NextBuy = NextBuy.Next()
				if NextBuy == nil {
					EventInfo.bBuyEventsEnd = true
				}
			} else if SellData.nGenerationID < BuyData.nGenerationID {
				OfferData = SellData
				// Especifica que eh um evento de venda
				EventInfo.bBuyEvent = false
				// Obtem o proximo evento de oferta de venda
				NextSell = NextSell.Next()
				if NextSell == nil {
					EventInfo.bSellEventsEnd = true
				}
			} else {
				logger.LogError(m_strLogFile, c_strMethodName, "Generation ID is equal for buy and sell offer : nGenerationID="+strconv.Itoa(BuyData.nGenerationID))
			}
		} else if NextBuy != nil && !EventInfo.bBuyEventsEnd {
			BuyData = NextBuy.Value.(OfferDataType)
			OfferData = BuyData
			// Especifica que eh um evento de compra
			EventInfo.bBuyEvent = true
			// Obtem o proximo evento de oferta de compra
			NextBuy = NextBuy.Next()
			if NextBuy == nil {
				EventInfo.bBuyEventsEnd = true
			}
		} else if NextSell != nil && !EventInfo.bSellEventsEnd {
			SellData = NextSell.Value.(OfferDataType)
			OfferData = SellData
			// Especifica que eh um evento de venda
			EventInfo.bBuyEvent = false
			// Obtem o proximo evento de oferta de venda
			NextSell = NextSell.Next()
			if NextSell == nil {
				EventInfo.bSellEventsEnd = true
			}
		} else {
			EventInfo.bProcessEvent = false
			logger.LogError(m_strLogFile, c_strMethodName, "NextBuy and NextSell are nil")
		}
		// Processa evento da oferta de compra ou venda
		if EventInfo.bProcessEvent {
			processOffer(a_DataInfo, OfferData, EventInfo.bBuyEvent)
		}
		// Condicao de parada -> os eventos foram processados
		if EventInfo.bBuyEventsEnd && EventInfo.bSellEventsEnd {
			break
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End")
}

func processOffer(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	const (
		c_strMethodName = "manager.processOffer"
	)
	switch a_OfferData.chOperation {
	case ofopCreation:
		processEventCreation(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopCancel:
		processEventCancel(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopEdit:
		processEventEdit(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopExpired:
		processEventExpired(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopReafirmed:
		processEventReafirmed(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopTrade:
		processEventTrade(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopUnknown:
		logger.LogError(m_strLogFile, c_strMethodName, "Unknown offer operation : chOperation="+string(a_OfferData.chOperation))
	default:
		logger.LogError(m_strLogFile, c_strMethodName, "Default offer operation : chOperation="+string(a_OfferData.chOperation))
	}
}

func processEventCreation(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	var (
		NewBookOffer BookOfferType
		BookOffer    BookOfferType
		BookOfferAux BookOfferType
		NewBookPrice BookPriceType
		BookPrice    BookPriceType
		BookPriceAux BookPriceType
		lstData      *list.List
		TempAux      *list.Element
		Temp         *list.Element
	)
	NewBookOffer.sPrice = a_OfferData.sPrice
	NewBookOffer.nQuantity = a_OfferData.nTotalQuantity
	NewBookOffer.nSecondaryID = a_OfferData.nSecondaryID
	NewBookOffer.nGenerationID = a_OfferData.nGenerationID
	NewBookOffer.strAccount = a_OfferData.strAccount

	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyOffers
	} else {
		lstData = &a_DataInfo.lstSellOffers
	}

	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookOffer = Temp.Value.(BookOfferType)
			if BookOffer.sPrice == NewBookOffer.sPrice {
				if BookOffer.nGenerationID < NewBookOffer.nGenerationID {
					TempAux = Temp.Next()
					if TempAux != nil {
						BookOfferAux = Temp.Value.(BookOfferType)
						if BookOfferAux.sPrice == NewBookOffer.sPrice {
							// 1 - Inserida entre ofertas do mesmo preco, seguindo a ordem crescente do numero de geracao
							if BookOfferAux.nGenerationID > NewBookOffer.nGenerationID {
								lstData.InsertAfter(NewBookOffer, Temp)
								break
							}
						} else {
							// 2 - Inserida depois da ultima oferta do preco em questao
							lstData.InsertAfter(NewBookOffer, Temp)
							break
						}
					} else {
						// 3 - Inserida no final da lista, o numero de geracao eh maior e nao existe proxima oferta
						lstData.PushBack(NewBookOffer)
						break
					}
				} else {
					// 4 - Inserida antes da oferta analisada, pois possui o mesmo preco e numero de geracao menor
					lstData.InsertBefore(NewBookOffer, Temp)
					break
				}
			} else if (BookOffer.sPrice > NewBookOffer.sPrice && a_bBuyEvent) || (BookOffer.sPrice < NewBookOffer.sPrice && !a_bBuyEvent) {
				// 5 - Inserida antes da oferta analisada, pois possui preco maior
				lstData.InsertBefore(NewBookOffer, Temp)
				break
			}

			Temp = Temp.Next()

			if Temp == nil {
				// 6 - Inserida no final da lista
				lstData.PushBack(NewBookOffer)
			}
		}
	} else {
		// 7 - Lista esta vazia, inserida no final
		lstData.PushBack(NewBookOffer)
	}

	NewBookPrice.sPrice = a_OfferData.sPrice
	NewBookPrice.nQuantity = a_OfferData.nTotalQuantity
	NewBookPrice.nCount = 1

	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyBookPrice
	} else {
		lstData = &a_DataInfo.lstSellBookPrice
	}

	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookPrice = Temp.Value.(BookPriceType)
			if BookPrice.sPrice == NewBookPrice.sPrice {
				// 1 - Atualizado grupo de preco existente
				BookPrice.nQuantity += a_OfferData.nTotalQuantity
				BookPrice.nCount++
				break
			} else if (BookOffer.sPrice > NewBookOffer.sPrice && a_bBuyEvent) || (BookOffer.sPrice < NewBookOffer.sPrice && !a_bBuyEvent) {
				// 2 - Inserida antes do grupo analisado, pois possui preco maior
				lstData.InsertBefore(NewBookPrice, Temp)
				break
			} else {
				TempAux = Temp.Next()
				if TempAux != nil {
					BookPriceAux = Temp.Value.(BookPriceType)
					if BookPriceAux.sPrice > NewBookPrice.sPrice {
						// 3 - Inserido entre grupos, seguindo a ordem crescente do preco
						lstData.InsertAfter(NewBookPrice, Temp)
						break
					}
				} else {
					// 4 - Inserido no final da lista, o preco eh maior e nao existe proximo grupo
					lstData.PushBack(NewBookPrice)
					break
				}
			}
			Temp = Temp.Next()
		}
	} else {
		// 5 - Lista esta vazia, inserido no final
		lstData.PushBack(NewBookPrice)
	}
}

func processEventCancel(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventEdit(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventExpired(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventReafirmed(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventTrade(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func exportResults(a_TickerData TickerDataType) {

}
