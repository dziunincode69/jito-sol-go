package helper

import (
	"context"
	"fmt"
	"jito_client/connection"
	"jito_client/utils"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type RaydiumTokenData struct {
	Id                 solana.PublicKey `json:"id"`
	BaseMint           solana.PublicKey `json:"baseMint"`
	QuoteMint          solana.PublicKey `json:"quoteMint"`
	LpMint             solana.PublicKey `json:"lpMint"`
	BaseDecimals       int64            `json:"baseDecimals"`
	QuoteDecimals      int64            `json:"quoteDecimals"`
	LpDecimals         int64            `json:"lpDecimals"`
	Version            int64            `json:"version"`
	ProgramId          solana.PublicKey `json:"programId"`
	Authority          solana.PublicKey `json:"authority"`
	OpenOrders         solana.PublicKey `json:"openOrders"`
	TargetOrders       solana.PublicKey `json:"targetOrders"`
	BaseVault          solana.PublicKey `json:"baseVault"`
	QuoteVault         solana.PublicKey `json:"quoteVault"`
	WithdrawQueue      solana.PublicKey `json:"withdrawQueue"`
	LpVault            solana.PublicKey `json:"lpVault"`
	MarketVersion      int64            `json:"marketVersion"`
	MarketProgramId    solana.PublicKey `json:"marketProgramId"`
	MarketId           solana.PublicKey `json:"marketId"`
	MarketAuthority    solana.PublicKey `json:"marketAuthority"`
	MarketBaseVault    solana.PublicKey `json:"marketBaseVault"`
	MarketQuoteVault   solana.PublicKey `json:"marketQuoteVault"`
	MarketBids         solana.PublicKey `json:"marketBids"`
	MarketAsks         solana.PublicKey `json:"marketAsks"`
	MarketEventQueue   solana.PublicKey `json:"marketEventQueue"`
	LookupTableAccount solana.PublicKey `json:"lookupTableAccount"`
}

func GetPool(token string) (solana.PublicKey, bool, error) {
	flip := false
	ctx := context.Background()
	srm := solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")

	publicKeyToken, _ := solana.PublicKeyFromBase58(token)
	publicKeySolana, _ := solana.PublicKeyFromBase58("So11111111111111111111111111111111111111112")

	offsets := []uint64{53, 85}
	var getAccount rpc.GetProgramAccountsResult
	var err error

	results := make(chan rpc.GetProgramAccountsResult, len(offsets))
	errors := make(chan error, len(offsets))

	for _, offset := range offsets {
		go func(offset uint64) {
			filters := []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{Bytes: publicKeyToken.Bytes(), Offset: offset},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{Bytes: publicKeySolana.Bytes(), Offset: 138 - offset},
				},
			}

			opts := &rpc.GetProgramAccountsOpts{
				Commitment: rpc.CommitmentConfirmed,
				Encoding:   solana.EncodingJSONParsed,
				Filters:    filters,
			}
			getAccount, err = connection.RpcClient().GetProgramAccountsWithOpts(ctx, srm, opts)
			results <- getAccount

			errors <- err
		}(offset)
	}

	for range offsets {
		select {
		case getAccount = <-results:
			if len(getAccount) > 0 {
				flip = true
			}
		case err = <-errors:
			if err != nil {
				return solana.PublicKey{}, false, err
			}
		}
	}

	if len(getAccount) == 0 {
		return solana.PublicKey{}, false, nil
	}

	return getAccount[0].Pubkey, flip, nil
}

func FindProgram(marketid solana.PublicKey, seed string) solana.PublicKey {
	RAYV4 := solana.MustPublicKeyFromBase58(utils.RAY_V4)
	res, _, err := solana.FindProgramAddress([][]byte{RAYV4.Bytes(), marketid.Bytes(), []byte(seed)}, RAYV4)
	if err != nil {
		log.Fatal(err, "FindProgramAddress")
	}
	return res
}

func PredictAddress(marketid solana.PublicKey) RaydiumTokenData {
	var tokenData RaydiumTokenData
	id := FindProgram(marketid, "amm_associated_seed")
	basevault := FindProgram(marketid, "coin_vault_associated_seed")
	quotevault := FindProgram(marketid, "pc_vault_associated_seed")
	lpmint := FindProgram(marketid, "lp_mint_associated_seed")
	lpvault := FindProgram(marketid, "temp_lp_token_associated_seed")
	targetOrder := FindProgram(marketid, "target_associated_seed")
	openOrder := FindProgram(marketid, "open_order_associated_seed")
	tokenData.Id = id
	tokenData.BaseVault = basevault
	tokenData.QuoteVault = quotevault
	tokenData.LpMint = lpmint
	tokenData.LpVault = lpvault
	tokenData.TargetOrders = targetOrder
	tokenData.OpenOrders = openOrder
	fmt.Println(id, "id")
	fmt.Println(basevault, "basevault")
	fmt.Println(quotevault, "quotevault")
	fmt.Println(lpmint, "lpmint")
	fmt.Println(lpvault, "lpvault")
	fmt.Println(targetOrder, "targetOrder")
	fmt.Println(openOrder, "openOrder")

	return tokenData
}
