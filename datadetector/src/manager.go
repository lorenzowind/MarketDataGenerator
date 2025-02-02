package src

import (
	"container/list"
	logger "marketmanipulationdetector/logger/src"
	"strconv"
)

func processEvents(a_TickerData *TickerDataType, a_DataInfo *DataInfoType) {
	const (
		c_strMethodName = "manager.processEvents"
	)
	var (
		NextBuy   *list.Element
		NextSell  *list.Element
		LastBuy   *list.Element
		LastSell  *list.Element
		BuyData   OfferDataType
		SellData  OfferDataType
		OfferData OfferDataType
		EventInfo EventInfoType
		nProgress int
	)
	logger.Log(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Begin")

	NextBuy = a_TickerData.lstBuy.Front()
	NextSell = a_TickerData.lstSell.Front()

	LastBuy = a_TickerData.lstBuy.Back()
	LastSell = a_TickerData.lstSell.Back()

	// Obtem o progresso do ticker com base no valor maximo lido
	if LastBuy != nil {
		BuyData = LastBuy.Value.(OfferDataType)
		OfferData = BuyData
		a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nMaxProgress = OfferData.nGenerationID
	}
	if LastSell != nil {
		SellData = LastSell.Value.(OfferDataType)
		OfferData = SellData
		if OfferData.nGenerationID > a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nMaxProgress {
			a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nMaxProgress = OfferData.nGenerationID
		}
	}

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
				logger.LogError(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Generation ID is equal for buy and sell offer : nGenerationID="+strconv.Itoa(BuyData.nGenerationID))
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
				EventInfo.bSellEventsEnd = true
			}
		} else if NextSell != nil && !EventInfo.bSellEventsEnd {
			SellData = NextSell.Value.(OfferDataType)
			OfferData = SellData
			// Especifica que eh um evento de venda
			EventInfo.bBuyEvent = false
			// Obtem o proximo evento de oferta de venda
			NextSell = NextSell.Next()
			if NextSell == nil {
				EventInfo.bBuyEventsEnd = true
				EventInfo.bSellEventsEnd = true
			}
		} else {
			EventInfo.bProcessEvent = false
			logger.LogError(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : NextBuy and NextSell are nil")
		}
		// Processa evento da oferta de compra ou venda
		if EventInfo.bProcessEvent {
			processOffer(a_TickerData, a_DataInfo, OfferData, EventInfo.bBuyEvent)
			processDetection(a_TickerData, a_DataInfo, OfferData, EventInfo.bBuyEvent)

			// Calcula novo progresso do ticker se eh um evento de criacao da oferta
			if OfferData.nOperation == ofopCreation && checkIfHasSameDate(OfferData.dtTime, a_TickerData.FilesInfo.TradeRunInfo.dtTickerDate) {
				nProgress = getProgress(OfferData.nGenerationID, a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nMaxProgress)
				if nProgress > a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nCurrentProgress {
					a_TickerData.FilesInfo.TradeRunInfo.ProgressInfo.nCurrentProgress = nProgress
					logger.Log(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : nProgress="+strconv.Itoa(nProgress)+"%")
				}
			}
		}
		// Condicao de parada -> os eventos foram processados
		if EventInfo.bBuyEventsEnd && EventInfo.bSellEventsEnd {
			break
		}
	}

	logger.Log(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Ticker events processed successfully")
	logger.Log(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : End")
}

func processOffer(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	const (
		c_strMethodName = "manager.processOffer"
	)
	switch a_OfferData.nOperation {
	case ofopCreation:
		processEventCreation(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopCancel:
		processEventCancel(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopEdit:
		processEventEdit(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopExpired:
		processEventCancel(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopReafirmed:
		processEventReafirmed(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopTrade:
		processEventTrade(a_DataInfo, a_OfferData, a_bBuyEvent)
	case ofopUnknown:
		logger.LogError(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Unknown offer operation : nOperation="+string(a_OfferData.nOperation))
	default:
		logger.LogError(m_LogInfo, "Ticker-Internal-Data", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Default offer operation : nOperation="+string(a_OfferData.nOperation))
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

	if bRemoved {
		if a_bBuyEvent {
			lstData = &a_DataInfo.lstBuyBookPrice
		} else {
			lstData = &a_DataInfo.lstSellBookPrice
		}

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
					break
				}
				Temp = Temp.Next()
			}
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
						BookPriceAux = TempAux.Value.(BookPriceType)
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

func processEventReafirmed(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
}

func processEventTrade(a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	var (
		lstData      *list.List
		Temp         *list.Element
		BookOffer    BookOfferType
		NewBookOffer BookOfferType
		BookPrice    BookPriceType
		NewBookPrice BookPriceType
		bRemoved     bool
	)
	if a_bBuyEvent {
		lstData = &a_DataInfo.lstBuyOffers
	} else {
		lstData = &a_DataInfo.lstSellOffers
	}

	Temp = lstData.Front()
	if Temp != nil {
		for Temp != nil {
			BookOffer = Temp.Value.(BookOfferType)
			if BookOffer.nSecondaryID == a_OfferData.nSecondaryID {
				NewBookOffer = BookOffer
				NewBookOffer.nQuantity -= a_OfferData.nTradeQuantity
				if NewBookOffer.nQuantity > 0 {
					lstData.InsertAfter(NewBookOffer, Temp)
				}
				lstData.Remove(Temp)
				bRemoved = true
				break
			}

			Temp = Temp.Next()
		}
	}

	if bRemoved {
		if a_bBuyEvent {
			lstData = &a_DataInfo.lstBuyBookPrice
		} else {
			lstData = &a_DataInfo.lstSellBookPrice
		}

		Temp = lstData.Front()
		if Temp != nil {
			for Temp != nil {
				BookPrice = Temp.Value.(BookPriceType)
				if BookPrice.sPrice == a_OfferData.sPrice {
					NewBookPrice = BookPrice
					NewBookPrice.nQuantity -= a_OfferData.nTradeQuantity
					if a_OfferData.nCurrentQuantity == 0 {
						NewBookPrice.nCount--
						if NewBookPrice.nCount > 0 {
							lstData.InsertAfter(NewBookPrice, Temp)
						}
					} else {
						lstData.InsertAfter(NewBookPrice, Temp)
					}
					lstData.Remove(Temp)
					break
				}
				Temp = Temp.Next()
			}
		}
	}
}

func getPriceLevel(a_DataInfo *DataInfoType, bBuy bool, a_nLevel int) float64 {
	var (
		lstData       *list.List
		Temp          *list.Element
		BookPrice     BookPriceType
		nListSize     int
		nCurrentLevel int
	)
	if bBuy {
		lstData = &a_DataInfo.lstBuyBookPrice
	} else {
		lstData = &a_DataInfo.lstSellBookPrice
	}

	Temp = lstData.Front()
	if Temp != nil {
		nListSize = lstData.Len()
		// Caso lista de precos seja menor que o nivel desejado, ja retorna preco da primeira posicao
		if nListSize < a_nLevel {
			BookPrice = Temp.Value.(BookPriceType)
			return BookPrice.sPrice
		}
		nCurrentLevel = 1
		for Temp != nil {
			// Retorna o preco do nivel desejado
			if nListSize-nCurrentLevel == a_nLevel {
				BookPrice = Temp.Value.(BookPriceType)
				return BookPrice.sPrice
			}
			Temp = Temp.Next()
			nCurrentLevel++
		}
	}

	return 0
}
