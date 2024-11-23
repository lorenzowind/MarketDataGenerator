package src

import (
	"container/list"
	"time"
)

type OfferOperationType byte

const (
	ofopCreation  OfferOperationType = '0' // evento de criacao da oferta
	ofopCancel                       = '4' // evento de cancelamento da oferta
	ofopEdit                         = '5' // evento de cancelamento da oferta
	ofopTrade                        = 'F' // evento de fechamento de negocio
	ofopExpired                      = 'C' // evento de expiracao da oferta
	ofopReafirmed                    = 'D' // evento de reafirmacao da oferta
	ofopUnknown                      = iota
)

type GenerationInfoType struct {
	strTickerName          string
	dtTickerDate           time.Time
	strReferenceTickerName string
	dtReferenceTickerDate  time.Time
}

type FilesInfoType struct {
	GenerationInfo            GenerationInfoType
	strReferenceBuyPath       string
	strReferenceSellPath      string
	strReferenceBenchmarkPath string
	strBuyPath                string
	strSellPath               string
	strBenchmarkPath          string
}

type MaskDataInfoType struct {
	hshMaskAccount     map[string]int
	hshMaskPrimaryID   map[int]int
	hshMaskSecondaryID map[int]int
	nCurrentAccount    int // tipo int sera convertido para string
	nCurrentPrimaryID  int
	nCurrentSecondaryD int
}

type TickerDataType struct {
	FilesInfo     *FilesInfoType
	MaskDataInfo  MaskDataInfoType
	lstBuy        list.List // doubly linked list de dados de ofertas de compra
	lstSell       list.List // doubly linked list de dados de ofertas de venda
	BenchmarkData BenchmarkDataType
}

type BenchmarkDataType struct {
	bHasBenchmarkData   bool
	dtAvgTradeInterval  time.Time
	sAvgOfferSize       float64
	sBiggerSDOfferSize  float64
	sSmallerSDOfferSize float64
}

type OfferDataType struct {
	chOperation      OfferOperationType
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
