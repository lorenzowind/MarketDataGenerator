package src

import (
	"container/list"
	"time"
)

type TradeOperationType byte
type OfferOperationType byte

const (
	tropBuy     TradeOperationType = 'C' // operacao de compra
	tropSell                       = 'V' // operacao de venda
	tropUnknown                    = iota
)

const (
	ofopCreation  OfferOperationType = '0' // evento de criacao da oferta
	ofopCancel                       = '4' // evento de cancelamento da oferta
	ofopEdit                         = '5' // evento de cancelamento da oferta
	ofopTrade                        = 'F' // evento de fechamento de negocio
	ofopExpired                      = 'C' // evento de expiracao da oferta
	ofopReafirmed                    = 'D' // evento de reafirmacao da oferta
	ofopUnknown                      = iota
)

type TradeRunInfoType struct {
	strTickerName string
	dtTickerDate  time.Time
}

type FilesInfoType struct {
	TradeRunInfo TradeRunInfoType
	strTradePath string
	strBuyPath   string
	strSellPath  string
}

type TickerDataType struct {
	FilesInfo FilesInfoType
	lstTrade  list.List // doubly linked list de dados de trade
	lstBuy    list.List // doubly linked list de dados de ofertas de compra
	lstSell   list.List // doubly linked list de dados de ofertas de venda
}

type TradeDataType struct {
	chOperation        TradeOperationType
	dtTime             time.Time
	strAccount         string
	nID                int
	nQuantity          int
	sPrice             float64
	nOfferGenerationID int
	nOfferPrimaryID    int
	nOfferSecondaryID  int
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
