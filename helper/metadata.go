package helper

import (
	"context"
	"jito_client/connection"

	bin "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"

	// "github.com/gagliardetto/solana-go/binary"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetMetadata(mint solana.PublicKey) (token_metadata.Data, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			token_metadata.ProgramID.Bytes(),
			mint.Bytes(),
		},
		token_metadata.ProgramID,
	)
	meta := GetMetadataData(connection.RpcClient(), addr)
	return meta, err
}
func GetMetadataData(rpcClient *rpc.Client, metadataPda solana.PublicKey) token_metadata.Data {
	accInfo, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), metadataPda, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if accInfo == nil {
		return token_metadata.Data{}
	}
	var data token_metadata.Metadata
	decoder := bin.NewBorshDecoder(accInfo.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}
	return data.Data

}
