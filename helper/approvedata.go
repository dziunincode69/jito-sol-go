package helper

import (
	"context"
	"jito_client/connection"
	"jito_client/wallet"
	"log"

	"github.com/gagliardetto/solana-go/rpc"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

type TokenAccountInfo struct {
	Mint    solana.PublicKey
	Account solana.PublicKey
}

func ApproveTokenDataInstrs(tokenSnipe string) bool {
	mints := []solana.PublicKey{
		solana.MustPublicKeyFromBase58(NativeSOL),
		solana.MustPublicKeyFromBase58(tokenSnipe),
	}

	var _, missingAccounts map[string]solana.PublicKey
	_, missingAccounts, err := GetTokenAccountsFromMints(context.Background(), *connection.RpcClient(), wallet.MainWallet().PublicKey(), mints...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var mintAddress solana.PublicKey

	mintAddress = solana.PublicKey{}

	if len(missingAccounts) != 0 {
		for mint := range missingAccounts {
			if mint == NativeSOL {
				continue
			}
			mintAddress = solana.MustPublicKeyFromBase58(mint)
		}
	}
	if mintAddress.String() == "11111111111111111111111111111111" {
		return false
	}
	return true
}

func GetTokenAccountsFromMints(
	ctx context.Context,
	clientRPC rpc.Client,
	owner solana.PublicKey,
	mints ...solana.PublicKey,
) (map[string]solana.PublicKey, map[string]solana.PublicKey, error) {

	duplicates := map[string]bool{}
	var tokenAccounts []solana.PublicKey
	var tokenAccountInfos []TokenAccountInfo
	for _, m := range mints {
		if ok := duplicates[m.String()]; ok {
			continue
		}
		duplicates[m.String()] = true
		a, _, err := solana.FindAssociatedTokenAddress(owner, m)
		if err != nil {
			return nil, nil, err
		}

		if m.String() == NativeSOL {
			a = owner
		}

		tokenAccounts = append(tokenAccounts, a)
		tokenAccountInfos = append(tokenAccountInfos, TokenAccountInfo{
			Mint:    m,
			Account: a,
		})
	}

	res, err := clientRPC.GetMultipleAccounts(ctx, tokenAccounts...)
	if err != nil {
		return nil, nil, err
	}

	missingAccounts := map[string]solana.PublicKey{}
	existingAccounts := map[string]solana.PublicKey{}
	for i, a := range res.Value {
		tai := tokenAccountInfos[i]
		if a == nil {
			missingAccounts[tai.Mint.String()] = tai.Account
			continue
		}

		if tai.Mint.String() == NativeSOL {
			existingAccounts[tai.Mint.String()] = owner
			continue
		}

		var ta token.Account
		err = bin.NewBinDecoder(a.Data.GetBinary()).Decode(&ta)
		if err != nil {
			return nil, nil, err
		}
		existingAccounts[tai.Mint.String()] = tai.Account
	}

	return existingAccounts, missingAccounts, nil
}
