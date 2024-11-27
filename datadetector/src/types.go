package src

import (
	"container/list"
	"time"
)

type OfferOperationType byte

type DetectionOperationType byte

const (
	ofopCreation  OfferOperationType = '0' // evento de criacao da oferta
	ofopCancel                       = '4' // evento de cancelamento da oferta
	ofopEdit                         = '5' // evento de cancelamento da oferta
	ofopTrade                        = 'F' // evento de fechamento de negocio
	ofopExpired                      = 'C' // evento de expiracao da oferta
	ofopReafirmed                    = 'D' // evento de reafirmacao da oferta
	ofopUnknown                      = iota
)

const (
	dtopSpoofing               DetectionOperationType = 0 // deteccao de spoofing - cenario tradicional
	dtopLayering                                      = 1 // deteccao de layering - cenario tradicional
	dtopLayeringModifiedOffers                        = 2 // deteccao de layering - cenario de modificacao de preco das ofertas
	dtopUnknown                                       = iota
)

type TradeRunInfoType struct {
	ProgressInfo  ProgressInfoType
	nCurrentRun   int
	strTickerName string
	dtTickerDate  time.Time
}

type ProgressInfoType struct {
	nMaxProgress     int // usando o valor do campo nGenerationID (esta ordenado mo arquivo externo)
	nCurrentProgress int // 0 a 100
}

type InfoForAllTickersType struct {
	nProcessors int
}

type ReportRunType struct {
	nTickers int
	dtStart  time.Time
}

type FilesInfoType struct {
	TradeRunInfo     TradeRunInfoType
	strBuyPath       string
	strSellPath      string
	strBenchmarkPath string
}

type TickerDataType struct {
	FilesInfo    *FilesInfoType
	lstBuy       list.List // doubly linked list de dados de ofertas de compra
	lstSell      list.List // doubly linked list de dados de ofertas de venda
	AuxiliarData AuxiliarDataType
	TempData     TempDataType
}

type AuxiliarDataType struct {
	hshOffersByPrimary   map[int][]*OfferDataType
	hshOffersBySecondary map[int][]*OfferDataType
	hshTradesBySecondary map[int][]*OfferDataType
	hshFullTrade         map[int]*FullTradeType
	hshTradesByAccount   map[string][]*FullTradeType
	BenchmarkData        BenchmarkDataType
}

type TempDataType struct {
	hshTradePrice map[int]TradePriceType
}

type TradePriceType struct {
	sTopBuyPriceLevel  float64
	sTopSellPriceLevel float64
	dtTradeTime        time.Time
}

type BenchmarkDataType struct {
	bHasBenchmarkData    bool
	dtAvgTradeInterval   time.Duration
	sAvgOfferSize        float64
	sSDOfferSize         float64
	sExpressiveOfferSize float64
}

type EventInfoType struct {
	bBuyEvent      bool
	bProcessEvent  bool
	bBuyEventsEnd  bool
	bSellEventsEnd bool
}

type DataInfoType struct {
	lstBuyBookPrice  list.List // doubly linked list dos grupos de preco de compra
	lstSellBookPrice list.List // doubly linked list dos grupos de preco de venda
	lstBuyOffers     list.List // doubly linked list das ofertas de compra no livro
	lstSellOffers    list.List // doubly linked list das ofertas de venda no livro
	arrDetectionData []DetectionDataType
}

type DetectionDataType struct {
	nOperation DetectionOperationType
}

type BookPriceType struct {
	sPrice    float64
	nQuantity int
	nCount    int
}

type BookOfferType struct {
	sPrice        float64
	nQuantity     int
	nPrimaryID    int
	nSecondaryID  int
	nGenerationID int
	strAccount    string
}

type OfferDataType struct {
	nOperation       OfferOperationType
	dtTime           time.Time
	strAccount       string
	nGenerationID    int
	nPrimaryID       int
	nSecondaryID     int
	nTradeID         int
	nCurrentQuantity int
	nTradeQuantity   int
	nTotalQuantity   int
	sPrice           float64
}

type FullTradeType struct {
	BuyOfferTrade  *OfferDataType
	SellOfferTrade *OfferDataType
}
