package helper

import (
	"encoding/binary"
	"math/big"
)

type LiquidityState struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       *big.Int
	SwapQuoteOutAmount     *big.Int
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      *big.Int
	SwapBaseOutAmount      *big.Int
	SwapQuote2BaseFee      uint64
	BaseVault              [32]byte // PublicKey is 32 bytes
	QuoteVault             [32]byte
	BaseMint               [32]byte
	QuoteMint              [32]byte
	LpMint                 [32]byte
	OpenOrders             [32]byte
	MarketId               [32]byte
	MarketProgramId        [32]byte
	TargetOrders           [32]byte
	WithdrawQueue          [32]byte
	LpVault                [32]byte
	Owner                  [32]byte
	LpReserve              uint64
	Padding                [24]byte // 3 uint64s
}

func DecodeLiquidityV4(data []byte) (*LiquidityState, error) {
	var state LiquidityState
	offset := 0

	readUint64 := func() uint64 {
		value := binary.LittleEndian.Uint64(data[offset : offset+8])
		offset += 8
		return value
	}

	readBigInt128 := func() *big.Int {
		bi := new(big.Int)
		bi.SetBytes(data[offset : offset+16])
		offset += 16
		return bi
	}

	readPublicKey := func() [32]byte {
		var pk [32]byte
		copy(pk[:], data[offset:offset+32])
		offset += 32
		return pk
	}

	state.Status = readUint64()
	state.Nonce = readUint64()
	state.MaxOrder = readUint64()
	state.Depth = readUint64()
	state.BaseDecimal = readUint64()
	state.QuoteDecimal = readUint64()
	state.State = readUint64()
	state.ResetFlag = readUint64()
	state.MinSize = readUint64()
	state.VolMaxCutRatio = readUint64()
	state.AmountWaveRatio = readUint64()
	state.BaseLotSize = readUint64()
	state.QuoteLotSize = readUint64()
	state.MinPriceMultiplier = readUint64()
	state.MaxPriceMultiplier = readUint64()
	state.SystemDecimalValue = readUint64()
	state.MinSeparateNumerator = readUint64()
	state.MinSeparateDenominator = readUint64()
	state.TradeFeeNumerator = readUint64()
	state.TradeFeeDenominator = readUint64()
	state.PnlNumerator = readUint64()
	state.PnlDenominator = readUint64()
	state.SwapFeeNumerator = readUint64()
	state.SwapFeeDenominator = readUint64()
	state.BaseNeedTakePnl = readUint64()
	state.QuoteNeedTakePnl = readUint64()
	state.QuoteTotalPnl = readUint64()
	state.BaseTotalPnl = readUint64()
	state.PoolOpenTime = readUint64()
	state.PunishPcAmount = readUint64()
	state.PunishCoinAmount = readUint64()
	state.OrderbookToInitTime = readUint64()
	state.SwapBaseInAmount = readBigInt128()
	state.SwapQuoteOutAmount = readBigInt128()
	state.SwapBase2QuoteFee = readUint64()
	state.SwapQuoteInAmount = readBigInt128()
	state.SwapBaseOutAmount = readBigInt128()
	state.SwapQuote2BaseFee = readUint64()
	state.BaseVault = readPublicKey()
	state.QuoteVault = readPublicKey()
	state.BaseMint = readPublicKey()
	state.QuoteMint = readPublicKey()
	state.LpMint = readPublicKey()
	state.OpenOrders = readPublicKey()
	state.MarketId = readPublicKey()
	state.MarketProgramId = readPublicKey()
	state.TargetOrders = readPublicKey()
	state.WithdrawQueue = readPublicKey()
	state.LpVault = readPublicKey()
	state.Owner = readPublicKey()
	state.LpReserve = readUint64()
	state.Padding = [24]byte{}

	return &state, nil
}
