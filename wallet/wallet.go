package wallet

import (
	"fmt"
	"jito_client/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/spf13/viper"
)

var (
	TipKeyPair    solana.PrivateKey
	MainKeyPair   solana.PrivateKey
	SecondKeyPair solana.PrivateKey
)

func ConnectTipWallet() {
	PK := viper.GetString("TipPrivateKey")
	fmt.Println(utils.DecryptString(PK))
	TipKeyPair = solana.MustPrivateKeyFromBase58(utils.DecryptString(PK))
}
func ConnectMultiWallet() {
	PK := viper.GetString("PrivateKey2")
	fmt.Println(utils.DecryptString(PK))
	SecondKeyPair = solana.MustPrivateKeyFromBase58(utils.DecryptString(PK))
}

func ConnectMainWallet() {
	PK := viper.GetString("PrivateKey")
	fmt.Println(utils.DecryptString(PK))
	MainKeyPair = solana.MustPrivateKeyFromBase58(utils.DecryptString(PK))
}
func SecondWallet() solana.PrivateKey {
	return SecondKeyPair
}

func MainWallet() solana.PrivateKey {
	return MainKeyPair
}

func TipWallet() solana.PrivateKey {
	return TipKeyPair
}
