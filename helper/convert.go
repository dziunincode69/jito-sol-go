package helper

import (
	"jito_client/lib/packet"

	"github.com/gagliardetto/solana-go"
)

func ConvertTransactionToProtobufPacket(transaction *solana.Transaction) (*packet.Packet, error) {
	rawTx, err := transaction.MarshalBinary()
	if err != nil {
		return &packet.Packet{}, err
	}
	return &packet.Packet{
		Data: rawTx,
		Meta: &packet.Meta{
			Size:        uint64(len(rawTx)),
			Addr:        "0.0.0.0",
			Port:        0,
			Flags:       nil,
			SenderStake: 0,
		},
	}, nil
}
func ConvertSolToLamport(amount float64) uint64 {
	return uint64(amount * 1000000000)
}
