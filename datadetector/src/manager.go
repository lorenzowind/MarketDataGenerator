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
	logger.Log(m_strLogFile, c_strMethodName, "Begin : strTicker="+a_TickerData.FilesInfo.TradeRunInfo.strTickerName)

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
			if BuyData.nGenerationID < SellData.nGenerationID || (BuyData.nGenerationID == SellData.nGenerationID && checkIfHasSameDate(BuyData.dtTime, a_TickerData.FilesInfo.TradeRunInfo.dtTickerDate)) {
				OfferData = BuyData
				// Especifica que eh um evento de compra
				EventInfo.bBuyEvent = true
				// Obtem o proximo evento de oferta de compra
				NextBuy = NextBuy.Next()
				if NextBuy == nil {
					EventInfo.bBuyEventsEnd = true
				}
			} else if SellData.nGenerationID < BuyData.nGenerationID || (BuyData.nGenerationID == SellData.nGenerationID && checkIfHasSameDate(SellData.dtTime, a_TickerData.FilesInfo.TradeRunInfo.dtTickerDate)) {
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
			printOfferData(OfferData)
			processOffer(a_DataInfo, OfferData, EventInfo.bBuyEvent)
		}
		// Condicao de parada -> os eventos foram processados
		if EventInfo.bBuyEventsEnd && EventInfo.bSellEventsEnd {
			break
		}
	}

	logger.Log(m_strLogFile, c_strMethodName, "End : strTicker="+a_TickerData.FilesInfo.TradeRunInfo.strTickerName)
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
	NewBookOffer.nPrimaryID = a_OfferData.nPrimaryID
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
				// 1 - Atualizado grupo de preco existente (adiciona atualizado e apaga antigo)
				NewBookPrice = BookPrice
				NewBookPrice.nQuantity += a_OfferData.nTotalQuantity
				NewBookPrice.nCount++
				lstData.InsertAfter(NewBookPrice, Temp)
				lstData.Remove(Temp)
				break
			} else if (BookPrice.sPrice > NewBookPrice.sPrice && a_bBuyEvent) || (BookPrice.sPrice < NewBookPrice.sPrice && !a_bBuyEvent) {
				// 2 - Inserida antes do grupo analisado, pois possui preco maior
				lstData.InsertBefore(NewBookPrice, Temp)
				break
			} else {
				TempAux = Temp.Next()
				if TempAux != nil {
					BookPriceAux = Temp.Value.(BookPriceType)
					if (BookPriceAux.sPrice > NewBookPrice.sPrice && a_bBuyEvent) || (BookPriceAux.sPrice < NewBookPrice.sPrice && !a_bBuyEvent) {
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
	const (
		//lint:ignore U1000 Ignore unused function
		c_strMethodName = "manager.processEventCancel"
	)
	var (
		BookOffer    BookOfferType
		NewBookPrice BookPriceType
		BookPrice    BookPriceType
		lstData      *list.List
		Temp         *list.Element
		bRemoved     bool
	)
	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyOffers
	} else {
		lstData = &a_DataInfo.lstSellOffers
	}

	bRemoved = false
	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookOffer = Temp.Value.(BookOfferType)
			if BookOffer.nSecondaryID == a_OfferData.nSecondaryID {
				lstData.Remove(Temp)
				bRemoved = true
				break
			}

			Temp = Temp.Next()
		}
	}

	if !bRemoved {
		//logger.LogError(m_strLogFile, c_strMethodName, "Offer not found : nGenerationID="+strconv.Itoa(a_OfferData.nGenerationID))
		//printOfferData(a_OfferData)
	} else {
		if a_bBuyEvent {
			lstData = &a_DataInfo.lstBuyBookPrice
		} else {
			lstData = &a_DataInfo.lstSellBookPrice
		}

		bRemoved = false
		Temp = lstData.Front()
		if Temp != nil {
			for Temp != nil {
				BookPrice = Temp.Value.(BookPriceType)
				if BookPrice.sPrice == a_OfferData.sPrice {
					NewBookPrice = BookPrice
					NewBookPrice.nQuantity -= a_OfferData.nCurrentQuantity
					NewBookPrice.nCount--
					if NewBookPrice.nCount > 0 {
						lstData.InsertAfter(NewBookPrice, Temp)
					}
					lstData.Remove(Temp)
					bRemoved = true
					break
				}
				Temp = Temp.Next()
			}
		}

		if !bRemoved {
			//logger.LogError(m_strLogFile, c_strMethodName, "Price not found : nGenerationID="+strconv.Itoa(a_OfferData.nGenerationID)+" : sPrice="+strconv.FormatFloat(a_OfferData.sPrice, 'f', -1, 64))
			//printOfferData(a_OfferData)
		}
	}
}

func processEventEdit(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	var (
		NewBookOffer BookOfferType
		BookOffer    BookOfferType
		BookOfferAux BookOfferType
		OldBookOffer BookOfferType
		NewBookPrice BookPriceType
		BookPrice    BookPriceType
		BookPriceAux BookPriceType
		lstData      *list.List
		TempAux      *list.Element
		Temp         *list.Element
		bInserted    bool
		bRemoved     bool
	)
	NewBookOffer.sPrice = a_OfferData.sPrice
	NewBookOffer.nQuantity = a_OfferData.nCurrentQuantity
	NewBookOffer.nPrimaryID = a_OfferData.nPrimaryID
	NewBookOffer.nSecondaryID = a_OfferData.nSecondaryID
	NewBookOffer.nGenerationID = a_OfferData.nGenerationID
	NewBookOffer.strAccount = a_OfferData.strAccount

	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyOffers
	} else {
		lstData = &a_DataInfo.lstSellOffers
	}

	bRemoved = false
	bInserted = false
	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookOffer = Temp.Value.(BookOfferType)
			// Verifica se encontrou a oferta antiga
			if !bRemoved && BookOffer.nPrimaryID == NewBookOffer.nPrimaryID {
				// Salva oferta antiga para atualizar nivel de preco
				OldBookOffer = BookOffer
				if BookOffer.sPrice != NewBookOffer.sPrice {
					// 1 - Remove antiga se preco eh diferente
					lstData.Remove(Temp)
					bRemoved = true
					if bInserted {
						break
					}
				} else {
					// 2 - Atualiza oferta se nao mudou preco
					lstData.InsertAfter(NewBookOffer, Temp)
					lstData.Remove(Temp)
					break
				}
			}
			// Verifica se deve inserir a oferta atualizada pois mudou de preco
			if !bInserted {
				if BookOffer.sPrice == NewBookOffer.sPrice {
					if BookOffer.nGenerationID < NewBookOffer.nGenerationID {
						TempAux = Temp.Next()
						if TempAux != nil {
							BookOfferAux = Temp.Value.(BookOfferType)
							if BookOfferAux.sPrice == NewBookOffer.sPrice {
								// 3 - Inserida entre ofertas do mesmo preco, seguindo a ordem crescente do numero de geracao
								if BookOfferAux.nGenerationID > NewBookOffer.nGenerationID {
									lstData.InsertAfter(NewBookOffer, Temp)
									bInserted = true
									if bRemoved {
										break
									}
								}
							} else {
								// 4 - Inserida depois da ultima oferta do preco em questao
								lstData.InsertAfter(NewBookOffer, Temp)
								bInserted = true
								if bRemoved {
									break
								}
							}
						} else {
							// 5 - Inserida no final da lista, o numero de geracao eh maior e nao existe proxima oferta
							lstData.PushBack(NewBookOffer)
							bInserted = true
							if bRemoved {
								break
							}
						}
					} else {
						// 6 - Inserida antes da oferta analisada, pois possui o mesmo preco e numero de geracao menor
						lstData.InsertBefore(NewBookOffer, Temp)
						bInserted = true
						if bRemoved {
							break
						}
					}
				} else if (BookOffer.sPrice > NewBookOffer.sPrice && a_bBuyEvent) || (BookOffer.sPrice < NewBookOffer.sPrice && !a_bBuyEvent) {
					// 7 - Inserida antes da oferta analisada, pois possui preco maior
					lstData.InsertBefore(NewBookOffer, Temp)
					bInserted = true
					if bRemoved {
						break
					}
				}
			}

			Temp = Temp.Next()

			if Temp == nil {
				// 8 - Inserida no final da lista
				lstData.PushBack(NewBookOffer)
			}
		}
	}

	NewBookPrice.sPrice = a_OfferData.sPrice
	NewBookPrice.nQuantity = a_OfferData.nCurrentQuantity
	NewBookPrice.nCount = 1

	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyBookPrice
	} else {
		lstData = &a_DataInfo.lstSellBookPrice
	}

	bRemoved = false
	bInserted = false
	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookPrice = Temp.Value.(BookPriceType)
			if !bRemoved {
				if BookPrice.sPrice == a_OfferData.sPrice {
					BookPriceAux = BookPrice
					if OldBookOffer.sPrice == a_OfferData.sPrice {
						// 1 - Atualiza grupo de preco depois da edicao da oferta (mesmo preco)
						BookPriceAux.nQuantity -= OldBookOffer.nQuantity
						BookPriceAux.nQuantity += a_OfferData.nCurrentQuantity
						lstData.InsertAfter(BookPriceAux, Temp)
						lstData.Remove(Temp)
						break
					} else {
						// 2 - Atualiza novo grupo de preco depois da edicao da oferta (preco diferente)
						BookPriceAux.nQuantity += a_OfferData.nCurrentQuantity
						BookPriceAux.nCount++
						lstData.InsertAfter(BookPriceAux, Temp)
						lstData.Remove(Temp)
						bRemoved = true
						if bInserted {
							break
						}
					}
				}

				if BookPrice.sPrice == OldBookOffer.sPrice {
					// 3 - Atualiza grupo de preco da oferta antiga
					BookPriceAux = BookPrice
					BookPriceAux.nQuantity -= OldBookOffer.nQuantity
					BookPriceAux.nCount--
					if BookPriceAux.nCount > 0 {
						lstData.InsertAfter(BookPriceAux, Temp)
					}
					lstData.Remove(Temp)
					bRemoved = true
					if bInserted {
						break
					}
				}
			}

			if !bInserted {
				if (BookPrice.sPrice > NewBookPrice.sPrice && a_bBuyEvent) || (BookPrice.sPrice < NewBookPrice.sPrice && !a_bBuyEvent) {
					// 4 - Inserida antes do grupo analisado, pois possui preco maior
					lstData.InsertBefore(NewBookPrice, Temp)
					bInserted = true
					if bRemoved {
						break
					}
				} else {
					TempAux = Temp.Next()
					if TempAux != nil {
						BookPriceAux = Temp.Value.(BookPriceType)
						if (BookPriceAux.sPrice > NewBookPrice.sPrice && a_bBuyEvent) || (BookPriceAux.sPrice < NewBookPrice.sPrice && !a_bBuyEvent) {
							// 5 - Inserido entre grupos, seguindo a ordem crescente do preco
							lstData.InsertAfter(NewBookPrice, Temp)
							bInserted = true
							if bRemoved {
								break
							}
						}
					} else {
						// 6 - Inserido no final da lista, o preco eh maior e nao existe proximo grupo
						lstData.PushBack(NewBookPrice)
						bInserted = true
						if bRemoved {
							break
						}
					}
				}
			}

			Temp = Temp.Next()
		}
	} else {
		// 7 - Lista esta vazia, inserido no final
		lstData.PushBack(NewBookPrice)
	}
}

func processEventExpired(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventReafirmed(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func processEventTrade(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {

}

func exportResults(a_TickerData TickerDataType) {

}
