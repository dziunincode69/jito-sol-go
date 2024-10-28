package price

import (
	"context"
	"fmt"
	"jito_client/connection"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func SimulateTransaction(ins []solana.Instruction) ([]string, error) {
	blockHash, err := connection.RpcClient().GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	accountNew := solana.MustPrivateKeyFromBase58("vwZ66sAP8GWPuabtoTdj49SfYQRueF7NnfaTx9veGQvYyQQDG2hNmxCfT58KGeKGKNCqN7QrehJyU3BhSUFpTtA")
	transactions, err := solana.NewTransaction(ins, blockHash.Value.Blockhash, solana.TransactionPayer(accountNew.PublicKey()))
	if err != nil {
		panic(err)
	}
	_, err = transactions.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountNew.PublicKey().Equals(key) {
				return &accountNew
			}
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	simulate, err := connection.RpcClient().SimulateTransaction(context.Background(), transactions)
	if err != nil {
		log.Fatal(err)
	}
	if len(simulate.Value.Logs) < 3 {
		return nil, fmt.Errorf("error: %s", simulate.Value.Logs)
	}
	return simulate.Value.Logs, nil
}
