package geyser

import (
	"jito_client/connection"
	geyser "jito_client/lib/geyser"

	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
)

func SubscribeAccountUpdates(account solana.PublicKey) (geyser.Geyser_SubscribeAccountUpdatesClient, error) {

	// program := solana.MustPublicKeyFromBase58("utils.RAY_V4")
	// programparams := &geyser.SubscribeProgramsUpdatesRequest{
	// 	Programs: [][]byte{program.Bytes()},
	// }
	subs, err := connection.JITOGeyser().SubscribeAccountUpdates(connection.GetContext(), &geyser.SubscribeAccountUpdatesRequest{
		Accounts: [][]byte{account.Bytes()},
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}
	return subs, err
	// fmt.Println(subs)

}
