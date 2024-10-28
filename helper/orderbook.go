package helper

import (
	"context"
	"jito_client/connection"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/serum"
)

type OrderBook struct {
	Bids                  string
	Asks                  string
	EventQueue            string
	SerumCoinVaultAccount string
	SerumPcVaultAccount   string
}

func FetchOrderBook(serumMarket solana.PublicKey) (OrderBook, error) {
	var fetch *serum.MarketMeta
	var err error
	fetch, err = serum.FetchMarket(context.Background(), connection.RpcClient(), serumMarket)
	if err != nil {
		return OrderBook{}, err
	}

	return OrderBook{
		Asks:                  fetch.MarketV2.Asks.String(),
		Bids:                  fetch.MarketV2.Bids.String(),
		EventQueue:            fetch.MarketV2.EventQueue.String(),
		SerumCoinVaultAccount: fetch.MarketV2.BaseVault.String(),
		SerumPcVaultAccount:   fetch.MarketV2.QuoteVault.String(),
	}, nil
}
