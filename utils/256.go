package utils

import (
	"math/big"

	ethamath "github.com/ethereum/go-ethereum/common/math"
)

func StringToBig256(str string) *big.Int {
	return ethamath.MustParseBig256(str)
}
