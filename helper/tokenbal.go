package helper

import (
	"context"
	"jito_client/connection"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTokenAccount(address solana.PublicKey) (*rpc.UiTokenAmount, error) {
	res, err := connection.RpcClient().GetTokenAccountBalance(context.Background(), address, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}
	return res.Value, nil

}
