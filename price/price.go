package price

import (
	"fmt"
	"jito_client/raydium"
	"jito_client/utils"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type PoolData struct {
	Status         int      `json:"status"`
	CoinDecimals   int      `json:"coin_decimals"`
	PcDecimals     int      `json:"pc_decimals"`
	LpDecimals     int      `json:"lp_decimals"`
	PoolPcAmount   *big.Int `json:"pool_pc_amount"`
	PoolCoinAmount *big.Int `json:"pool_coin_amount"`
	PnlPcAmount    *big.Int `json:"pnl_pc_amount"`
	PnlCoinAmount  *big.Int `json:"pnl_coin_amount"`
	PoolLpSupply   *big.Int `json:"pool_lp_supply"`
	AmmId          string   `json:"amm_id"`
}

var keys solana.AccountMetaSlice

func GetTknPrice(tokenBalanceLamports *big.Int, tokenData raydium.RaydiumTokenData, ins []solana.Instruction) (string, error) {
	isFetched := false
	var resp []string
	var err error
	for !isFetched {
		resp, err = SimulateTransaction(ins)
		if err == nil && len(resp) > 4 {
			isFetched = true
		}
	}

	jsonStr, err := utils.ParseSimulateLogToJson(resp[4], "GetPoolData")
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	var pooldata PoolData

	status, err := utils.ParseSimulateValueAsInt(jsonStr, "status")
	if err != nil {
		fmt.Println("Error parsing status:", err)
		return "", err
	}

	baseDecimals, err := utils.ParseSimulateValueAsInt(jsonStr, "coin_decimals")
	if err != nil {
		fmt.Println("Error parsing coin_decimals:", err)
		return "", err
	}

	quoteDecimals, err := utils.ParseSimulateValueAsInt(jsonStr, "pc_decimals")
	if err != nil {
		fmt.Println("Error parsing pc_decimals:", err)
		return "", err
	}

	lpDecimals, err := utils.ParseSimulateValueAsInt(jsonStr, "lp_decimals")
	if err != nil {
		fmt.Println("Error parsing lp_decimals:", err)
		return "", err
	}

	baseReserve, err := utils.ParseSimulateValue(jsonStr, "pool_coin_amount")
	if err != nil {
		fmt.Println("Error parsing pool_coin_amount:", err)
		return "", err
	}

	quoteReserve, err := utils.ParseSimulateValue(jsonStr, "pool_pc_amount")
	if err != nil {
		fmt.Println("Error parsing pool_pc_amount:", err)
		return "", err
	}

	lpSupply, err := utils.ParseSimulateValue(jsonStr, "pool_lp_supply")
	if err != nil {
		fmt.Println("Error parsing pool_lp_supply:", err)
		return "", err
	}
	pooldata.Status = status
	pooldata.CoinDecimals = baseDecimals
	pooldata.PcDecimals = quoteDecimals
	pooldata.LpDecimals = lpDecimals
	pooldata.PoolPcAmount = quoteReserve
	pooldata.PoolCoinAmount = baseReserve
	pooldata.PoolLpSupply = lpSupply

	result := CalculateAmountOut(tokenBalanceLamports, pooldata)
	return result, nil
}
