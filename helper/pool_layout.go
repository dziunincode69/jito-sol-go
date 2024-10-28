package helper

import (
	"bytes"
	"fmt"
	"jito_client/utils"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type POOL_INFO_LAYOUT struct {
	bin.BaseVariant
	Instruction             uint8
	SimulateType            uint8
	solana.AccountMetaSlice `bin:"-" borsh_skip:"false"`
}

func (inst *POOL_INFO_LAYOUT) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(utils.RAY_V4)
}

func (inst *POOL_INFO_LAYOUT) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *POOL_INFO_LAYOUT) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	encoder.WriteUint8(inst.Instruction)
	encoder.WriteUint8(inst.SimulateType)
	return nil
}
func (inst *POOL_INFO_LAYOUT) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}
