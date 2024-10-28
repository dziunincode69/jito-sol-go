package price

import (
	"jito_client/helper"
	"jito_client/raydium"
	"jito_client/utils"

	"github.com/gagliardetto/solana-go"
)

func MakeSimulatePoolInfoInstruction(poolkeys raydium.RaydiumTokenData) []solana.Instruction {
	layoutStruct := helper.POOL_INFO_LAYOUT{
		Instruction:  12, // Corresponds to `instruction: 12` in the JS code
		SimulateType: 0,  // Corresponds to `simulateType: 0` in the JS code
	}
	layout, err := layoutStruct.Data()
	if err != nil {
		panic(err)
	}
	keys = []*solana.AccountMeta{
		solana.Meta(poolkeys.Id),
		solana.Meta(poolkeys.Authority),
		solana.Meta(poolkeys.OpenOrders),
		solana.Meta(poolkeys.BaseVault),
		solana.Meta(poolkeys.QuoteVault),
		solana.Meta(poolkeys.LpMint),
		solana.Meta(poolkeys.MarketId),
		solana.Meta(poolkeys.MarketEventQueue),
	}

	insarr := []solana.Instruction{solana.NewInstruction(solana.MustPublicKeyFromBase58(utils.RAY_V4), keys, layout)}
	return insarr
}
