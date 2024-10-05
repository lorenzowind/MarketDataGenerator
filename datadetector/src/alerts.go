package src

func processDetection(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	checkSpoofing(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
	checkLayering(a_TickerData, a_DataInfo, a_OfferData, a_bBuyEvent)
}

func checkSpoofing(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
	// Ciclo de spoofing - precisa ocorrer um cancelamento da oferta
	if a_OfferData.chOperation == ofopCancel {
		// Verifica se a oferta cancelada possui uma quantidade expressiva
		if IsExpressiveOffer(a_TickerData, a_DataInfo, a_OfferData) {
		}
	}
}

func checkLayering(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType, a_bBuyEvent bool) {
}

func IsExpressiveOffer(a_TickerData *TickerDataType, a_DataInfo *DataInfoType, a_OfferData OfferDataType) bool {
	return true
}

func exportResults(a_TickerData *TickerDataType) {
}
