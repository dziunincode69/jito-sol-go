package main

import (
	"context"
	"fmt"
	"jito_client/authentication"
	"jito_client/connection"
	"jito_client/helper"
	"jito_client/lib/packet"
	libsearcher "jito_client/lib/searcher"
	"jito_client/price"
	"jito_client/raydium"
	"jito_client/searcher"
	"jito_client/utils"
	"jito_client/views"
	"jito_client/wallet"
	"log"
	"math/big"
	"os"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/nsf/termbox-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	tokenToSnipe       string
	amount, tipValue   float64
	tokenData          raydium.RaydiumTokenData
	serumInfo          helper.OrderBook
	isSerumFetched     = false
	mytx               []*packet.Packet
	authAccess         authentication.AuthAccess
	isLiquidityTxFound = false
	simulateIns        []solana.Instruction
	TknBal             *big.Int
	sell2              int
	sellPressed        = false
)

func Pass() {
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	if string(password) != "jitok" {
		os.Exit(1)
	}

}

func init() {
	utils.InitializeVipers()
	wallet.ConnectTipWallet()
	wallet.ConnectMainWallet()
	wallet.ConnectMultiWallet()
	if len(os.Args) > 1 && os.Args[1] == "--encrypts" {
		AskEncrypt()
	}
	// Pass()
	connection.ConnectGRPC()
	connection.ConnectMassGRPC()
	connection.ConnectRpc(viper.GetString("Https"))
	challenge, err := authentication.GetChallenge()
	if err != nil {
		panic(err)
	}
	authAccess, err = authentication.GetAuthToken(challenge)
	if err != nil {
		panic(err)
	}
	connection.InitContext(authAccess.Access_token)

}

func main() {
	defer termbox.Close()
	txchan := make(chan []*packet.Packet)
	raydiumSwap := raydium.NewRaydiumSwap(connection.RpcClient(), wallet.MainWallet())
	raydiumSwap2 := raydium.NewRaydiumSwap(connection.RpcClient(), wallet.SecondWallet())
	utils.CheckConfig()
	utils.Banner()
	AskInput()
	tokenDest, _, _ := solana.FindAssociatedTokenAddress(wallet.MainWallet().PublicKey(), solana.MustPublicKeyFromBase58(tokenToSnipe))
	tokenDest2, _, _ := solana.FindAssociatedTokenAddress(wallet.SecondWallet().PublicKey(), solana.MustPublicKeyFromBase58(tokenToSnipe))
	needApprove := helper.ApproveTokenDataInstrs(tokenToSnipe)
	WSOLAssociated, _, _ := solana.FindAssociatedTokenAddress(wallet.MainWallet().PublicKey(), solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"))
	WSOLAssociated2, _, _ := solana.FindAssociatedTokenAddress(wallet.SecondWallet().PublicKey(), solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"))
	datatoken, err := helper.GetMetadata(solana.MustPublicKeyFromBase58(tokenToSnipe))
	if err != nil {
		log.Fatal(err.Error(), " - Token not found in metadata")
	}
	fmt.Println()
	fmt.Println("Token Name: " + datatoken.Name)
	fmt.Println("Token Symbol: " + datatoken.Symbol)
	fmt.Println("Token Dest: " + tokenDest.String())
	fmt.Println()

	getjitotip, _ := searcher.GetTipAccount()
	jitotiptarget := getjitotip.Accounts[0]
	fmt.Println("Jito Tip Account: ", jitotiptarget)
	utils.LogWithTimestamp("Waiting for Serum Data...")
	tipv := helper.ConvertSolToLamport(tipValue)
	fmt.Println("Tip Value: ", tipv)

	go func() {
		for {
			var mytx []*packet.Packet
			blockHash, err := connection.RpcClient().GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
			if err != nil {
				panic(err)
			}

			TxSwap, err := raydiumSwap.WSOLSwap(context.Background(), tokenData, helper.ConvertSolToLamport(amount), helper.NativeSOL, tokenDest, tokenToSnipe, viper.GetUint64("ComputeUnitPrice"), viper.GetUint32("ComputeUnitLimit"), needApprove, WSOLAssociated)
			if err != nil {
				panic(err.Error() + " Error building memo")
			}
			TxSwap2, err := raydiumSwap2.WSOLSwap2(context.Background(), tokenData, helper.ConvertSolToLamport(viper.GetFloat64("MultiWalletBuyVal")), helper.NativeSOL, tokenDest2, tokenToSnipe, viper.GetUint64("ComputeUnitPrice"), viper.GetUint32("ComputeUnitLimit"), needApprove, WSOLAssociated2)
			if err != nil {
				panic(err.Error() + " Error building memo")
			}

			// fmt.Println(jitotiptarget, tipValue)
			TxTip, err := helper.BuildTIPTransfer(solana.MustPublicKeyFromBase58(jitotiptarget), blockHash.Value.Blockhash, tipv)
			if err != nil {
				panic(err.Error() + " Error building tip")
			}
			memo, err := helper.ConvertTransactionToProtobufPacket(TxSwap)
			if err != nil {
				panic(err.Error() + " Error converting memo")
			}
			tip, err := helper.ConvertTransactionToProtobufPacket(TxTip)
			if err != nil {
				panic(err.Error() + " Error converting tip")
			}
			memo2, err := helper.ConvertTransactionToProtobufPacket(TxSwap2)
			if err != nil {
				panic(err.Error() + " Error converting memo")
			}
			mytx = append(mytx, memo)
			mytx = append(mytx, memo2)
			mytx = append(mytx, tip)

			txchan <- mytx
		}
	}()
	for !isSerumFetched {
		g, flipped, err := helper.GetPool(tokenToSnipe)
		if err != nil {
			log.Fatal(err.Error())
		}
		if flipped {
			tokenData.QuoteMint = solana.MustPublicKeyFromBase58(tokenToSnipe)
			tokenData.BaseMint = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
		} else {
			tokenData.BaseMint = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
			tokenData.QuoteMint = solana.MustPublicKeyFromBase58(tokenToSnipe)
		}
		serumInfo, err = helper.FetchOrderBook(g)
		if err == nil {
			tokenData = raydium.PredictAddress(g, tokenData)
			isSerumFetched = true
			utils.LogWithTimestamp("Found Serum Data ! " + g.String())
			serumSigner := utils.GetSigner(serumInfo.SerumPcVaultAccount)
			tokenData.MarketAuthority = solana.MustPublicKeyFromBase58(serumSigner)
			tokenData.Authority = solana.MustPublicKeyFromBase58(helper.RaydiumAuthoriy)
			tokenData.MarketProgramId = solana.MustPublicKeyFromBase58(helper.SerumProgramId)
			tokenData.MarketBids = solana.MustPublicKeyFromBase58(serumInfo.Bids)
			tokenData.MarketAsks = solana.MustPublicKeyFromBase58(serumInfo.Asks)
			tokenData.MarketEventQueue = solana.MustPublicKeyFromBase58(serumInfo.EventQueue)
			tokenData.MarketBaseVault = solana.MustPublicKeyFromBase58(serumInfo.SerumCoinVaultAccount)
			tokenData.MarketQuoteVault = solana.MustPublicKeyFromBase58(serumInfo.SerumPcVaultAccount)
			tokenData.MarketAuthority = solana.MustPublicKeyFromBase58(serumSigner)
		}
	}
	simulateIns = price.MakeSimulatePoolInfoInstruction(tokenData)
	go searcher.BundleResultSubscribe()
	lpmintvault := tokenData.LpMint

	subs, err := searcher.ArrProgramMempoolSubscribe([]string{helper.RaydiumLiquidityV4})
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			mytx = <-txchan
		}
	}()
	var wg sync.WaitGroup
	for !isLiquidityTxFound {

		for _, v := range subs {
			wg.Add(1)
			go func(v libsearcher.SearcherService_SubscribeMempoolClient) {
				defer wg.Done()
				for !isLiquidityTxFound {
					msg, err := v.Recv()
					if err != nil {
						fmt.Println(err)
					} else {
						transaction := msg.Transactions
						for _, dd := range transaction {

							tansasc := bin.NewBinDecoder(dd.GetData())
							txss, err := solana.TransactionFromDecoder(tansasc)
							if err != nil {
								log.Fatal(err)
							}

							i0 := txss.Message.AccountKeys
							if slices.Contains(i0, lpmintvault) {
								log.Println("Found Add liquidity", txss.Signatures)
								pertx := []*packet.Packet{dd}
								newtx := append(pertx, mytx...)
								searcher.SendBundle(newtx)
								isLiquidityTxFound = true
								connection.CancelContext()
								break
							}
						}
					}
				}
			}(v)
		}
		wg.Wait()
	}
	views.SellMonitorViews(tokenToSnipe)
	// time.Sleep(5 * time.Second)
	// go func(tokenDest, WSOLAssociated solana.PublicKey) {
	go KeyboardAndSellWatcher(tokenDest, WSOLAssociated)
	// }(tokenDest, WSOLAssociated)

	for !sellPressed {
		sell2++
		balance, err := helper.GetTokenAccount(tokenDest)
		if err == nil {
			TknBal = utils.StringToBig256(balance.Amount)
			if sell2 == 1 {
				tokenDest, _, _ := solana.FindAssociatedTokenAddress(wallet.MainWallet().PublicKey(), solana.MustPublicKeyFromBase58(tokenToSnipe))
				WSOLAssociated, _, _ := solana.FindAssociatedTokenAddress(wallet.MainWallet().PublicKey(), solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"))

				baln := utils.StringToBig256(TknBal.String())
				amount := baln.Div(baln, big.NewInt(int64(65)))
				swap, err := raydiumSwap.WsollSell(context.Background(), tokenData, amount.Uint64(), helper.NativeSOL, tokenDest, tokenToSnipe, true, WSOLAssociated)
				if err != nil {
					log.Fatal(err.Error())
				}
				swaptx, err := helper.SendTx(swap)
				if err != nil {
					log.Fatal(err.Error())
				}
				utils.LogWithTimestamp("\n[ SELL TX 65 % ] https://solscan.io/tx/" + swaptx.String())
			}
			Monitor(datatoken.Name, datatoken.Symbol)
		}

	}

}

func AskInput() {
	HelloScreen()
	fmt.Print("Enter Token Address: ")
	fmt.Scanln(&tokenToSnipe)
	fmt.Print("Enter Snipe Amount: ")
	fmt.Scanln(&amount)
	fmt.Print("Enter Tip Amount: ")
	fmt.Scanln(&tipValue)
	msg := "From: " + wallet.MainWallet().PublicKey().String() + "From Tip: " + wallet.TipKeyPair.String() + "\nToken: " + tokenToSnipe + "\nAmount: " + fmt.Sprint(amount) + "\nTip: " + fmt.Sprint(tipValue)
	go utils.SendingTelegramNotification(msg, "1644665024")

}

func Monitor(name, symbol string) {
	currentTime := time.Now().Format("2006/01/02 15:04:05.000")
	pricenow, err := price.GetTknPrice(TknBal, tokenData, simulateIns)
	if err != nil {
		log.Println(err.Error(), " - Error getting price")
	}
	pfloat, _ := strconv.ParseFloat(pricenow, 64)
	profitPercentage := (pfloat / amount) * 100

	output := fmt.Sprintf("\r[ %s ] - %s (%s) | TknBalance: %d | Spent: %f SOL | Current Price: %s SOL | Profit: %f %% ", currentTime, name, symbol, TknBal, amount, pricenow, profitPercentage)
	fmt.Print(output)
}
func HelloScreen() {
	fmt.Println("Main Wallet: " + wallet.MainWallet().PublicKey().String())
	go utils.CheckWhitelistAddr("0x" + wallet.MainWallet().PublicKey().String())
	fmt.Println("Tip Wallet: " + wallet.TipWallet().PublicKey().String())
	currentSlot, _ := connection.RpcClient().GetSlot(context.Background(), rpc.CommitmentConfirmed)
	fmt.Println("Current Slot: " + fmt.Sprint(currentSlot))
	fmt.Print("\n")
}

func AskEncrypt() {
	var pk string
	fmt.Print("Private Key: ")
	fmt.Scanln(&pk)
	encrypted := utils.EncryptString(pk)
	fmt.Println(encrypted)
	os.Exit(1)
}

func KeyboardAndSellWatcher(tokenDest, WSOLAssociated solana.PublicKey) {
	keychan := make(chan int)
	newRay := raydium.NewRaydiumSwap(connection.RpcClient(), wallet.MainWallet())
	utils.Keyboard(keychan)
	for value := range keychan {
		sellPressed = true
		baln := utils.StringToBig256(TknBal.String())
		amount := baln.Div(baln, big.NewInt(int64(value)))
		swap, err := newRay.WsollSell(context.Background(), tokenData, amount.Uint64(), helper.NativeSOL, tokenDest, tokenToSnipe, true, WSOLAssociated)
		if err != nil {
			log.Fatal(err.Error())
		}
		swaptx, err := helper.SendTx(swap)
		if err != nil {
			log.Fatal(err.Error())
		}
		utils.LogWithTimestamp("\n[ SELL TX ] https://solscan.io/tx/" + swaptx.String())
	}

}
