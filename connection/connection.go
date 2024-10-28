package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"jito_client/lib/auth"
	"jito_client/lib/block_engine"
	geyser "jito_client/lib/geyser"
	searcher "jito_client/lib/searcher"
	"jito_client/lib/shredstream"
	"jito_client/utils"
	"strings"

	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	AuthService        auth.AuthServiceClient
	SearcherService    searcher.SearcherServiceClient
	GeyserService      geyser.GeyserClient
	ctx                context.Context
	ctx2               context.Context
	rpcClient          *rpc.Client
	ShreadService      shredstream.ShredstreamClient
	BlockEngineService block_engine.BlockEngineValidatorClient
	connectionGRPCARR  []searcher.SearcherServiceClient
	cancel             context.CancelFunc
)

func ConnectRpc(url string) {
	rpcClient = rpc.New(url)
}
func CancelContext() {

	cancel()
}

func ConnectMassGRPC() {
	ctx, cancel = context.WithCancel(context.Background())
	maxMsgSize := 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 7
	strcountry := []string{utils.Tokyo.BlockEngineUrl, utils.Amsterdam.BlockEngineUrl, utils.Frankfurt.BlockEngineUrl, utils.NewYork.BlockEngineUrl}
	for _, v := range strcountry {
		conn, err := grpc.DialContext(
			ctx,
			v,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				ServerName: strings.Split(v, ":")[0],
			})),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		)
		if err != nil {
			panic(fmt.Errorf("GRPC did not connect: %v", err))
		}
		search := searcher.NewSearcherServiceClient(conn)

		connectionGRPCARR = append(connectionGRPCARR, search)
		// connectionGRPCARR = append(connectionGRPCARR, conn)

	}
}
func GetGRPCConnection() []searcher.SearcherServiceClient {
	return connectionGRPCARR
}

func ConnectGRPC() {
	ctx, cancel = context.WithCancel(context.Background())
	maxMsgSize := 1024 * 1024 * 1024
	conn, err := grpc.DialContext(
		context.Background(),
		utils.NewYork.BlockEngineUrl,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName: "ny.mainnet.block-engine.jito.wtf",
		})),
		// grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	)
	if err != nil {
		panic(fmt.Errorf("GRPC did not connect: %v", err))
	}

	AuthService = auth.NewAuthServiceClient(conn)
	SearcherService = searcher.NewSearcherServiceClient(conn)
	ShreadService = shredstream.NewShredstreamClient(conn)
	BlockEngineService = block_engine.NewBlockEngineValidatorClient(conn)

}
func JITOBES() block_engine.BlockEngineValidatorClient {
	return BlockEngineService
}

func JITOShread() shredstream.ShredstreamClient {
	return ShreadService
}
func JITOAuth() auth.AuthServiceClient {
	return AuthService
}
func JITOSearcher() searcher.SearcherServiceClient {
	return SearcherService
}
func JITOGeyser() geyser.GeyserClient {
	return GeyserService
}

func InitContext(access_token string) {
	context := context.Background()
	ctx = metadata.AppendToOutgoingContext(context, "Authorization", "Bearer "+access_token)
	ctx2 = metadata.AppendToOutgoingContext(context, "api-key", "97389dac-a75e-45ab-bdff-df5bf1cb2383")
}
func GetContext() context.Context {
	return ctx
}
func GetContext2() context.Context {
	return ctx2
}
func RpcClient() *rpc.Client {
	return rpcClient
}
