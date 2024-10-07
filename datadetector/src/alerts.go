package src

func processDetection(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	// So realiza a deteccao caso tenha encontrado os valores de benchmark
	if a_TickerData.AuxiliarData.BenchmarkData.bHasBenchmarkData {
		checkSpoofing(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
		checkLayering(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
	}
}

func checkSpoofing(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	var (
		lstOfferData []*OfferDataType
		nIndex       int
	)
	// Verifica cenario em que a oferta deixou de ser expressiva
	if a_OfferData.chOperation == ofopEdit {
		lstOfferData = getOffersByPrimaryID(a_TickerData, a_OfferData.nPrimaryID)
		for nIndex = 1; nIndex < len(lstOfferData); nIndex++ {
			if lstOfferData[nIndex-1].nSecondaryID == a_OfferData.nSecondaryID {
				if IsExpressiveOffer(a_TickerData, *lstOfferData[nIndex-1]) && !IsExpressiveOffer(a_TickerData, a_OfferData) {
				}
			}
		}
	}
}

func checkLayering(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
}

func IsExpressiveOffer(a_TickerData *TickerDataType, a_OfferData OfferDataType) bool {
	return float64(a_OfferData.nTotalQuantity) >= a_TickerData.AuxiliarData.BenchmarkData.sExpressiveOfferSize
}

func exportResults(a_TickerData *TickerDataType) {
}
