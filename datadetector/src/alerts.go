package src

import "time"

func processDetection(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	// So realiza a deteccao caso tenha encontrado os valores de benchmark
	if a_TickerData.AuxiliarData.BenchmarkData.bHasBenchmarkData {
		checkSpoofing(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
		checkLayering(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
	}
}

func checkSpoofing(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	var (
		OriginalSpoofingOffer *OfferDataType
		OriginalSpoofingTrade *OfferDataType
	)
	OriginalSpoofingOffer = getOriginalSpoofingOffer(a_TickerData, a_OfferData)

	if OriginalSpoofingOffer != nil {
		OriginalSpoofingTrade = getOriginalSpoofingTrade(a_TickerData, a_OfferData, a_bBuyEvent, OriginalSpoofingOffer)
		if OriginalSpoofingTrade != nil {

		}
	}
}

func getOriginalSpoofingOffer(a_TickerData *TickerDataType, a_OfferData OfferDataType) *OfferDataType {
	var (
		OriginalSpoofingOffer *OfferDataType
		lstOfferData          []*OfferDataType
		nIndex                int
	)
	OriginalSpoofingOffer = nil

	lstOfferData = getOffersByPrimaryID(a_TickerData, a_OfferData.nPrimaryID)
	for nIndex = 1; nIndex < len(lstOfferData); nIndex++ {
		// Verifica se oferta antiga era expressiva
		if IsExpressiveOffer(a_TickerData, *lstOfferData[nIndex-1]) {
			if OriginalSpoofingOffer == nil {
				// Obtem oferta antiga que era expressiva
				OriginalSpoofingOffer = lstOfferData[nIndex-1]
			}
		} else {
			// Caso oferta antiga deixou de ser expressiva seta para nulo
			OriginalSpoofingOffer = nil
		}
		// Verifica se oferta eh igual a atual, pois tem o mesmo ID de geracao
		if lstOfferData[nIndex].nGenerationID == a_OfferData.nGenerationID {
			// Se oferta atual deixou de expressiva ou eh igual a cancelada
			if !IsExpressiveOffer(a_TickerData, a_OfferData) || a_OfferData.chOperation == ofopCancel {
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

func getOriginalSpoofingTrade(a_TickerData *TickerDataType, a_OfferData OfferDataType, a_bBuyEvent bool, a_OriginalSpoofingOffer *OfferDataType) *OfferDataType {
	var (
		AccountTrade        *OfferDataType
		NearestAccountTrade *OfferDataType
		lstAccountTrades    []*FullTradeType
		nIndex              int
		dtTradeDiff         time.Duration
		dtNearestTrade      time.Duration
	)
	dtNearestTrade = 0
	NearestAccountTrade = nil

	lstAccountTrades = getTradesByAccount(a_TickerData, a_OfferData.strAccount)
	for nIndex = 0; nIndex < len(lstAccountTrades); nIndex++ {
		// Obtem o evento do trade no lado oposto da oferta de spoofing
		if a_bBuyEvent {
			AccountTrade = lstAccountTrades[nIndex].SellOfferTrade
		} else {
			AccountTrade = lstAccountTrades[nIndex].BuyOfferTrade
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

	if NearestAccountTrade != nil {
		// Verifica se o trade aconteceu antes da oferta de spoofing
		if a_OfferData.dtTime.Before(NearestAccountTrade.dtTime) {
			// Compara o tempo entre o trade e a oferta de spoofing para verificar se esta dentro do benchmark
			if LessOrEqualThanTradeInverval(a_TickerData, dtNearestTrade) {
				return NearestAccountTrade
			}
		} else {
			// Compara o tempo entre o trade e a oferta expressiva para verificar se esta dentro do benchmark
			if IsBetweenTradeInverval(a_TickerData, a_OriginalSpoofingOffer.dtTime, NearestAccountTrade.dtTime) {
				return NearestAccountTrade
			}
		}
	}

	return nil
}

func checkLayering(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
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

func exportResults(a_TickerData *TickerDataType) {
}
