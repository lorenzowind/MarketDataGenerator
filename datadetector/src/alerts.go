package src

import (
	logger "marketmanipulationdetector/logger/src"
	"strconv"
	"time"
)

const (
	c_nTopPriceLevel = 5
)

func processDetection(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	// So realiza a deteccao caso tenha encontrado os valores de benchmark
	if a_TickerData.AuxiliarData.BenchmarkData.bHasBenchmarkData {
		// Verifica se eh evento de trade
		if a_OfferData.nOperation == ofopTrade {
			// Armazena estado do livro
			processTradePrice(a_TickerData, a_DataInfo, a_OfferData)
			// Layering - detecta cenario tradicional
			checkLayering(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
			// Verifica se eh quaisquer outro evento
		} else {
			// Verifica se eh evento de edicao
			if a_OfferData.nOperation == ofopEdit {
				// Layering - detecta cenario de modificacao de preco das ofertas
				checkLayeringModifiedOffers(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
			}
			// Spoofing - detecta cenario tradicional
			checkSpoofing(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
		}
	}
}

func processTradePrice(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType) {
	var (
		TradePrice TradePriceType
		bKeyExists bool
	)
	_, bKeyExists = a_TickerData.TempData.hshTradePrice[a_OfferData.nTradeID]
	if !bKeyExists {
		TradePrice.dtTradeTime = a_OfferData.dtTime
		TradePrice.sTopBuyPriceLevel = getPriceLevel(a_DataInfo, true, c_nTopPriceLevel)
		TradePrice.sTopSellPriceLevel = getPriceLevel(a_DataInfo, false, c_nTopPriceLevel)

		a_TickerData.TempData.hshTradePrice[a_OfferData.nTradeID] = TradePrice
	}
}

func checkSpoofing(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	const (
		c_strMethodName = "alerts.checkSpoofing"
	)
	var (
		OriginalSpoofingOffer *OfferDataType
		OriginalSpoofingTrade *OfferDataType
		// SpoofingTrade         *OfferDataType
		// arrSpoofingTrades     []*OfferDataType
	)
	OriginalSpoofingOffer = getOriginalSpoofingOffer(a_TickerData, a_OfferData)
	if OriginalSpoofingOffer != nil {
		OriginalSpoofingTrade = getOriginalSpoofingTrade(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent, OriginalSpoofingOffer)
		if OriginalSpoofingTrade != nil {
			// Concatena deteccao (dtopSpoofing)
			a_DataInfo.arrDetectionData = append(a_DataInfo.arrDetectionData, DetectionDataType{
				nOperation: dtopSpoofing,
			})

			logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Spoofing detected : dtopSpoofing : nCurrent="+strconv.Itoa(getDetectionDataLength(a_DataInfo.arrDetectionData, dtopSpoofing))+" : nTotal="+strconv.Itoa(len(a_DataInfo.arrDetectionData)))

			// logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Actual offer : "+getOfferData(a_OfferData))
			// logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Original spoofing offer : "+getOfferData(*OriginalSpoofingOffer))
			// logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Original spoofing trade : "+getOfferData(*OriginalSpoofingTrade))

			// arrSpoofingTrades = getSpoofingTrades(a_TickerData, a_bBuyEvent, OriginalSpoofingTrade)
			// logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Spoofing trades count : "+strconv.Itoa(len(arrSpoofingTrades)))
			// for _, SpoofingTrade = range arrSpoofingTrades {
			// 	logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, "Spoofing trade : "+getOfferData(*SpoofingTrade))
			// }
		}
	}
}

func getOriginalSpoofingOffer(a_TickerData *TickerDataType, a_OfferData OfferDataType) *OfferDataType {
	var (
		OriginalSpoofingOffer *OfferDataType
		arrOfferData          []*OfferDataType
		nIndex                int
	)
	OriginalSpoofingOffer = nil

	arrOfferData = getOffersByPrimaryID(a_TickerData, a_OfferData.nPrimaryID)
	for nIndex = 1; nIndex < len(arrOfferData); nIndex++ {
		// Verifica se oferta antiga era expressiva
		if IsExpressiveOffer(a_TickerData, *arrOfferData[nIndex-1]) {
			if OriginalSpoofingOffer == nil {
				// Obtem oferta antiga que era expressiva
				OriginalSpoofingOffer = arrOfferData[nIndex-1]
			}
		} else {
			// Caso oferta antiga deixou de ser expressiva seta para nulo
			OriginalSpoofingOffer = nil
		}
		// Verifica se oferta eh igual a atual, pois tem o mesmo ID de geracao
		if arrOfferData[nIndex].nGenerationID == a_OfferData.nGenerationID {
			// Se oferta atual deixou de expressiva ou eh igual a cancelada
			if !IsExpressiveOffer(a_TickerData, a_OfferData) || a_OfferData.nOperation == ofopCancel {
				if OriginalSpoofingOffer != nil {
					// Verifica se validade da oferta expressiva esta entre o tempo de benchmark de intervalo entre negocios
					if IsBetweenTradeInverval(a_TickerData, a_OfferData.dtTime, OriginalSpoofingOffer.dtTime) {
						return OriginalSpoofingOffer
					}
				}
			}
			// Encerra verificacao, pois ja passou da oferta atual
			break
		}
	}

	return nil
}

func getOriginalSpoofingTrade(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool, a_OriginalSpoofingOffer *OfferDataType) *OfferDataType {
	var (
		AccountTrade        *OfferDataType
		NearestAccountTrade *OfferDataType
		arrAccountTrades    []*FullTradeType
		nIndex              int
		dtTradeDiff         time.Duration
		dtNearestTrade      time.Duration
		sTopPriceLevel      float64
	)
	dtNearestTrade = 0
	NearestAccountTrade = nil

	arrAccountTrades = getTradesByAccount(a_TickerData, a_OfferData.strAccount)
	for nIndex = 0; nIndex < len(arrAccountTrades); nIndex++ {
		// Verifica se o trade de compra e venda de uma conta eh valido (possivelmente se aconteceu no mesmo dia)
		if arrAccountTrades[nIndex].BuyOfferTrade != nil && arrAccountTrades[nIndex].SellOfferTrade != nil {
			// Obtem o evento do trade no lado oposto da oferta de spoofing
			if a_bBuyEvent {
				AccountTrade = arrAccountTrades[nIndex].SellOfferTrade
			} else {
				AccountTrade = arrAccountTrades[nIndex].BuyOfferTrade
			}
			// Verifica diferenca de tempo entre o evento de trade e a oferta de spoofing
			dtTradeDiff = a_OfferData.dtTime.Sub(AccountTrade.dtTime)
			if dtTradeDiff < dtNearestTrade || dtNearestTrade == 0 {
				dtNearestTrade = dtTradeDiff
				// Armazena o evento de trade mais proximo
				NearestAccountTrade = AccountTrade
			} else {
				break
			}
		}
	}

	if NearestAccountTrade != nil {
		// Verifica se o trade aconteceu antes da oferta de spoofing
		if a_OfferData.dtTime.Before(NearestAccountTrade.dtTime) {
			// Compara o tempo entre o trade e a oferta de spoofing para verificar se esta dentro do benchmark
			if LessOrEqualThanTradeInverval(a_TickerData, dtNearestTrade) {
				// Verifica oferta original se esta dentro dos niveis de preco de spoofing
				if a_bBuyEvent {
					sTopPriceLevel = getPriceLevel(a_DataInfo, true, c_nTopPriceLevel)
					// Lado de compra -> preco maior ou igual esta dentro dos niveis
					if a_OriginalSpoofingOffer.sPrice < sTopPriceLevel {
						return nil
					}
				} else {
					sTopPriceLevel = getPriceLevel(a_DataInfo, false, c_nTopPriceLevel)
					// Lado de venda -> preco menor ou igual esta dentro dos niveis
					if a_OriginalSpoofingOffer.sPrice > sTopPriceLevel {
						return nil
					}
				}
				return NearestAccountTrade
			}
		} else {
			// Compara o tempo entre o trade e a oferta expressiva para verificar se esta dentro do benchmark
			if IsBetweenTradeInverval(a_TickerData, a_OriginalSpoofingOffer.dtTime, NearestAccountTrade.dtTime) {
				// Verifica trade mais proximo se esta dentro dos niveis de preco de spoofing
				if a_bBuyEvent {
					sTopPriceLevel = getTradePrice(a_TickerData, true, NearestAccountTrade.nTradeID)
					// Lado de compra -> preco maior ou igual esta dentro dos niveis
					if a_OriginalSpoofingOffer.sPrice < sTopPriceLevel {
						return nil
					}
				} else {
					sTopPriceLevel = getTradePrice(a_TickerData, false, NearestAccountTrade.nTradeID)
					// Lado de venda -> preco menor ou igual esta dentro dos niveis
					if a_OriginalSpoofingOffer.sPrice > sTopPriceLevel {
						return nil
					}
				}
				return NearestAccountTrade
			}
		}
	}

	return nil
}

//lint:ignore U1000 Ignore unused function
func getSpoofingTrades(a_TickerData *TickerDataType, a_bBuyEvent bool, a_OriginalSpoofingTrade *OfferDataType) []*OfferDataType {
	var (
		arrSpoofingTrades []*OfferDataType
		TradeAux          *OfferDataType
		arrTrades         []*OfferDataType
		FullTrade         *FullTradeType
	)
	arrSpoofingTrades = make([]*OfferDataType, 0)

	arrTrades = getTradesBySecondaryID(a_TickerData, a_OriginalSpoofingTrade.nSecondaryID)
	for _, TradeAux = range arrTrades {
		FullTrade = getFullTrade(a_TickerData, TradeAux.nTradeID)
		if FullTrade != nil {
			// Obtem o trade do lado oposto que foi manipulado pela oferta expressiva
			if a_bBuyEvent {
				arrSpoofingTrades = append(arrSpoofingTrades, FullTrade.SellOfferTrade)
			} else {
				arrSpoofingTrades = append(arrSpoofingTrades, FullTrade.BuyOfferTrade)
			}
		}
	}

	return arrSpoofingTrades
}

func checkLayering(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
}

func checkLayeringModifiedOffers(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
}

func IsBetweenTradeInverval(a_TickerData *TickerDataType, a_dtLeft, a_dtRight time.Time) bool {
	return LessOrEqualThanTradeInverval(a_TickerData, a_dtLeft.Sub(a_dtRight))
}

func LessOrEqualThanTradeInverval(a_TickerData *TickerDataType, a_dtDiff time.Duration) bool {
	return a_dtDiff <= a_TickerData.AuxiliarData.BenchmarkData.dtAvgTradeInterval
}

func IsExpressiveOffer(a_TickerData *TickerDataType, a_OfferData OfferDataType) bool {
	return float64(a_OfferData.nTotalQuantity) >= a_TickerData.AuxiliarData.BenchmarkData.sExpressiveOfferSize
}

func exportResults(a_TickerData *TickerDataType, a_DataInfo *DataInfoType) {
	const (
		c_strMethodName = "alerts.exportResults"
	)
	logger.Log(m_LogInfo, "Manipulation-Spoofing", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Results exported successfully : dtopSpoofing="+strconv.Itoa(getDetectionDataLength(a_DataInfo.arrDetectionData, dtopSpoofing))+" : nTotal="+strconv.Itoa(len(a_DataInfo.arrDetectionData)))
	logger.Log(m_LogInfo, "Manipulation-Layering", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Results exported successfully : dtopLayering="+strconv.Itoa(getDetectionDataLength(a_DataInfo.arrDetectionData, dtopLayering))+" : nTotal="+strconv.Itoa(len(a_DataInfo.arrDetectionData)))
	logger.Log(m_LogInfo, "Manipulation-Layering", c_strMethodName, getHeaderRun(a_TickerData.FilesInfo.TradeRunInfo)+" : Results exported successfully : dtopLayeringModifiedOffers : nResults="+strconv.Itoa(getDetectionDataLength(a_DataInfo.arrDetectionData, dtopLayeringModifiedOffers))+" : nTotal="+strconv.Itoa(len(a_DataInfo.arrDetectionData)))
}
