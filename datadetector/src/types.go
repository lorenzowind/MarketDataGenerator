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
}

type BookPriceType struct {
	sPrice    float64
	nQuantity int
	nCount    int
}

type BookOfferType struct {
	sPrice        float64
	nQuantity     int
	nSecondaryID  int
	nGenerationID int
	strAccount    string
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
