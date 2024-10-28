package helper

import (
	"bytes"
	"context"
	"jito_client/connection"
	"jito_client/wallet"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

var TraderAPIMemoProgram = solana.MustPublicKeyFromBase58("HQ2UUt18uJqKaQFJhgV9zaTdQxUZjNrsKFgoEDquBkcx")

func BuildTIPTransfer(receipent solana.PublicKey, blockHash solana.Hash, val uint64) (*solana.Transaction, error) {
	transfer := system.NewTransferInstruction(val, wallet.TipWallet().PublicKey(), receipent).Build()
	ins := []solana.Instruction{
		transfer,
	}
	txn, err := solana.NewTransaction(ins, blockHash, solana.TransactionPayer(wallet.TipKeyPair.PublicKey()))
	if err != nil {
		return nil, err
	}
	_, err = txn.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if wallet.TipKeyPair.PublicKey().Equals(key) {
				return &wallet.TipKeyPair
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return txn, nil
}
func BuildMemo(msg string, blockHash solana.Hash) (*solana.Transaction, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(msg))
	instruction := &solana.GenericInstruction{
		AccountValues: nil,
		ProgID:        TraderAPIMemoProgram,
		DataBytes:     buf.Bytes(),
	}
	txn, err := solana.NewTransaction([]solana.Instruction{instruction}, blockHash, solana.TransactionPayer(wallet.TipWallet().PublicKey()))
	if err != nil {
		return nil, err
	}
	_, err = txn.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if wallet.TipWallet().PublicKey().Equals(key) {
				return &wallet.TipKeyPair
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func SendTx(tx *solana.Transaction) (solana.Signature, error) {
	sig, err := connection.RpcClient().SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight: true,
		},
	)
	if err != nil {
		return solana.Signature{}, err
	}
	return sig, nil

}
