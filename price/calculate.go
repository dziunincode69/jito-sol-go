package price

import (
	"jito_client/utils"
	"math/big"
)

func CalculateAmountOut(amount *big.Int, poolInfo PoolData) string {
	reserves := []*big.Int{poolInfo.PoolCoinAmount, poolInfo.PoolPcAmount}
	feeRaw := new(big.Int).Mul(amount, big.NewInt(int64(utils.LIQUIDITY_FEES_NUMERATOR)))
	feeRaw.Div(feeRaw, big.NewInt(int64(utils.LIQUIDITY_FEES_DENOMINATOR)))
	amountInWithFee := new(big.Int).Sub(amount, feeRaw)
	denominator := new(big.Int).Add(reserves[0], amountInWithFee)
	amountOut := new(big.Int).Mul(reserves[1], amountInWithFee)
	amountOut.Div(amountOut, denominator)
	solDecimals := big.NewFloat(1e9)
	amountOutFloat := new(big.Float).SetInt(amountOut)
	solValue := new(big.Float).Quo(amountOutFloat, solDecimals)
	solValueText := solValue.Text('f', 6)

	return solValueText
}
