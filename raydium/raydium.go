package raydium

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"jito_client/utils"
	"jito_client/wallet"
	"log"

	associatedTokenAccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/spf13/viper"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	TokenAccountSize = 165
	NativeSOL        = "11111111111111111111111111111111"
	WrappedSOL       = "So11111111111111111111111111111111111111112"
)

// var (
// 	TESTED_DATA = RaydiumTokenData{
// 		Id:               "ABF4G5MXV76WkosBjwDXy438S8UiDnMUjvvvdb2Uiu2i",
// 		BaseVault:        "DeTDv9YTFmoZyrEZnhEqaP2qw7qht1vswx3xrxYAYJk2",
// 		QuoteVault:       "EZ85M5RuJe6LNc2gC1ikTEUG5ATqXgJxj7KuCNMMfZsN",
// 		Authority:        "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1",
// 		OpenOrders:       "tgqJZ9hWkwXDkU6M6Pyb2pnCH6PMbqmNhHBxpA3kP1A",
// 		TargetOrders:     "HpPUAJXVUH1sjmm3zRymAsvUd55KgpYc6EyU8h7Uv4Lc",
// 		MarketProgramId:  "srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX",
// 		MarketId:         "Goh6jNH7uUXSV1u2gzdbbYFexSGjLYrMAC8CwsaYCF7T",
// 		MarketBids:       "NQHBrbbEe2hrgauE1mm1KgrdZ51J7JuKMcF8pttgVXk",
// 		MarketAsks:       "6W5aMiCX2h7qXSQqdxRh3AUWCMhqdT2frd6xzk5215Et",
// 		MarketEventQueue: "AG3oxthFzN4aFdR3eQXBzfZA7zPoNqqaceiECexqZQyA",
// 		MarketBaseVault:  "35nGe2g2JD7shNoLrGXQASnKcoRK1Rpq1LAAkgFrdizE",
// 		MarketQuoteVault: "AQq78yMTgDUpoRmm8hEBLLaJdbYmUZJ3LvxBx5Pepxz7",
// 		MarketAuthority:  "9meqxQ2W3qdj84JogUCUXwCDAdp92y2UbGAnz8tTKpnH",
// 	}
// )

type RaydiumSwap struct {
	clientRPC *rpc.Client
	account   solana.PrivateKey
}

type RaydiumTokenData struct {
	Id                 solana.PublicKey `json:"id"`
	BaseMint           solana.PublicKey `json:"baseMint"`
	QuoteMint          solana.PublicKey `json:"quoteMint"`
	LpMint             solana.PublicKey `json:"lpMint"`
	BaseDecimals       int64            `json:"baseDecimals"`
	QuoteDecimals      int64            `json:"quoteDecimals"`
	LpDecimals         int64            `json:"lpDecimals"`
	Version            int64            `json:"version"`
	ProgramId          solana.PublicKey `json:"programId"`
	Authority          solana.PublicKey `json:"authority"`
	OpenOrders         solana.PublicKey `json:"openOrders"`
	TargetOrders       solana.PublicKey `json:"targetOrders"`
	BaseVault          solana.PublicKey `json:"baseVault"`
	QuoteVault         solana.PublicKey `json:"quoteVault"`
	WithdrawQueue      solana.PublicKey `json:"withdrawQueue"`
	LpVault            solana.PublicKey `json:"lpVault"`
	MarketVersion      int64            `json:"marketVersion"`
	MarketProgramId    solana.PublicKey `json:"marketProgramId"`
	MarketId           solana.PublicKey `json:"marketId"`
	MarketAuthority    solana.PublicKey `json:"marketAuthority"`
	MarketBaseVault    solana.PublicKey `json:"marketBaseVault"`
	MarketQuoteVault   solana.PublicKey `json:"marketQuoteVault"`
	MarketBids         solana.PublicKey `json:"marketBids"`
	MarketAsks         solana.PublicKey `json:"marketAsks"`
	MarketEventQueue   solana.PublicKey `json:"marketEventQueue"`
	LookupTableAccount solana.PublicKey `json:"lookupTableAccount"`
}

func FindProgram(marketid solana.PublicKey, seed string) solana.PublicKey {
	RAYV4 := solana.MustPublicKeyFromBase58(utils.RAY_V4)
	res, _, err := solana.FindProgramAddress([][]byte{RAYV4.Bytes(), marketid.Bytes(), []byte(seed)}, RAYV4)
	if err != nil {
		log.Fatal(err, "FindProgramAddress")
	}
	return res
}
func (s *RaydiumSwap) WsollSell(
	ctx context.Context,
	pool RaydiumTokenData,
	amount uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	needApprove bool,
	WSOLAssociated solana.PublicKey,

) (*solana.Transaction, error) {
	minimumOutAmount := uint64(0)
	var instruction []solana.Instruction

	signers := []solana.PrivateKey{s.account}
	setFee, _ := computebudget.NewSetComputeUnitPriceInstruction(9205000).ValidateAndBuild()
	setUnitLimit, _ := computebudget.NewSetComputeUnitLimitInstruction(600000).ValidateAndBuild()
	instruction = append(instruction, setFee)
	instruction = append(instruction, setUnitLimit)

	instAprove, err := associatedTokenAccount.NewCreateInstruction(
		s.account.PublicKey(),
		s.account.PublicKey(),
		solana.MustPublicKeyFromBase58(WrappedSOL),
	).ValidateAndBuild()
	if err != nil {
		log.Fatalln(err.Error())
	}
	instruction = append(instruction, instAprove)
	instruction = append(instruction, NewRaydiumSwapInstruction(
		amount,
		minimumOutAmount,
		solana.TokenProgramID,
		pool.Id,
		pool.Authority,
		pool.OpenOrders,
		pool.TargetOrders,
		pool.BaseVault,
		pool.QuoteVault,
		pool.MarketProgramId,
		pool.MarketId,
		pool.MarketBids,
		pool.MarketAsks,
		pool.MarketEventQueue,
		pool.MarketBaseVault,
		pool.MarketQuoteVault,
		pool.MarketAuthority,
		fromAccount,
		WSOLAssociated,
		s.account.PublicKey(),
	))

	closeInst, err := token.NewCloseAccountInstruction(
		WSOLAssociated,
		s.account.PublicKey(),
		s.account.PublicKey(),
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	instruction = append(instruction, closeInst)

	tx, err := ExecuteInstructions(ctx, s.clientRPC, signers, instruction...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
func PredictAddress(marketid solana.PublicKey, tokenData RaydiumTokenData) RaydiumTokenData {
	id := FindProgram(marketid, "amm_associated_seed")
	basevault := FindProgram(marketid, "coin_vault_associated_seed")
	quotevault := FindProgram(marketid, "pc_vault_associated_seed")
	lpmint := FindProgram(marketid, "lp_mint_associated_seed")
	lpvault := FindProgram(marketid, "temp_lp_token_associated_seed")
	targetOrder := FindProgram(marketid, "target_associated_seed")
	openOrder := FindProgram(marketid, "open_order_associated_seed")
	tokenData.Id = id
	tokenData.BaseVault = basevault
	tokenData.MarketId = marketid
	tokenData.QuoteVault = quotevault
	tokenData.LpMint = lpmint
	tokenData.LpVault = lpvault
	tokenData.TargetOrders = targetOrder
	tokenData.OpenOrders = openOrder

	return tokenData
}

// func GetTokenMetadata(tokenAddress string) RaydiumTokenData {
// 	resp, err := http.Get("https://api.raydium.io/v2/sdk/liquidity/mainnet.json")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	jsonString := []string{"unOfficial"}
// 	var response RaydiumTokenData

// 	_, _ = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 		arraySearch, err := jsonparser.GetString(value, "baseMint")
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		arraySearch2, err := jsonparser.GetString(value, "quoteMint")
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		if arraySearch == tokenAddress || arraySearch2 == tokenAddress {
// 			response.Id, _ = jsonparser.GetString(value, "id")
// 			response.BaseMint, _ = jsonparser.GetString(value, "baseMint")
// 			response.QuoteMint, _ = jsonparser.GetString(value, "quoteMint")
// 			response.LpMint, _ = jsonparser.GetString(value, "lpMint")
// 			response.BaseDecimals, _ = jsonparser.GetInt(value, "baseDecimals")
// 			response.QuoteDecimals, _ = jsonparser.GetInt(value, "quoteDecimals")
// 			response.LpDecimals, _ = jsonparser.GetInt(value, "lpDecimals")
// 			response.Version, _ = jsonparser.GetInt(value, "version")
// 			response.ProgramId, _ = jsonparser.GetString(value, "programId")
// 			response.Authority, _ = jsonparser.GetString(value, "authority")
// 			response.OpenOrders, _ = jsonparser.GetString(value, "openOrders")
// 			response.TargetOrders, _ = jsonparser.GetString(value, "targetOrders")
// 			response.BaseVault, _ = jsonparser.GetString(value, "baseVault")
// 			response.QuoteVault, _ = jsonparser.GetString(value, "quoteVault")
// 			response.WithdrawQueue, _ = jsonparser.GetString(value, "withdrawQueue")
// 			response.LpVault, _ = jsonparser.GetString(value, "lpVault")
// 			response.MarketVersion, _ = jsonparser.GetInt(value, "marketVersion")
// 			response.MarketProgramId, _ = jsonparser.GetString(value, "marketProgramId")
// 			response.MarketId, _ = jsonparser.GetString(value, "marketId")
// 			response.MarketAuthority, _ = jsonparser.GetString(value, "marketAuthority")
// 			response.MarketBaseVault, _ = jsonparser.GetString(value, "marketBaseVault")
// 			response.MarketQuoteVault, _ = jsonparser.GetString(value, "marketQuoteVault")
// 			response.MarketBids, _ = jsonparser.GetString(value, "marketBids")
// 			response.MarketAsks, _ = jsonparser.GetString(value, "marketAsks")
// 			response.MarketEventQueue, _ = jsonparser.GetString(value, "marketEventQueue")
// 			response.LookupTableAccount, _ = jsonparser.GetString(value, "lookupTableAccount")
// 		}
// 	}, jsonString...)

// 	return response
// }

func NewRaydiumSwap(clientRPC *rpc.Client, account solana.PrivateKey) *RaydiumSwap {
	return &RaydiumSwap{clientRPC: clientRPC, account: account}
}

func (s *RaydiumSwap) Swap(
	ctx context.Context,
	pool RaydiumTokenData,
	amount uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	unitprice uint64,
	unitlimit uint32,
	needApprove bool,
) (*solana.Transaction, error) {
	minimumOutAmount := uint64(50)

	var instruction []solana.Instruction
	signers := []solana.PrivateKey{s.account}
	tempAccount := solana.NewWallet()
	needWrapSOL := fromToken == NativeSOL || toToken == NativeSOL
	setFee, err := computebudget.NewSetComputeUnitPriceInstruction(unitprice).ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	setUnitLimit, err := computebudget.NewSetComputeUnitLimitInstruction(unitlimit).ValidateAndBuild()
	instruction = append(instruction, setUnitLimit)
	instruction = append(instruction, setFee)
	if needApprove {
		instAprove, err := associatedTokenAccount.NewCreateInstruction(
			wallet.MainKeyPair.PublicKey(),
			wallet.MainKeyPair.PublicKey(),

			solana.MustPublicKeyFromBase58(toToken),
		).ValidateAndBuild()
		if err != nil {
			log.Fatalln(err.Error())
		}

		instruction = append(instruction, instAprove)
	}
	if needWrapSOL {
		rentCost, err := s.clientRPC.GetMinimumBalanceForRentExemption(
			ctx,
			TokenAccountSize,
			"",
		)
		if err != nil {
			return nil, err
		}
		accountLamports := rentCost
		if fromToken == NativeSOL {
			accountLamports += amount
		}

		createInst, err := system.NewCreateAccountInstruction(
			accountLamports,
			TokenAccountSize,
			solana.TokenProgramID,
			s.account.PublicKey(),
			tempAccount.PublicKey(),
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instruction = append(instruction, createInst)
		initInst, err := token.NewInitializeAccountInstruction(
			tempAccount.PublicKey(),
			solana.MustPublicKeyFromBase58(WrappedSOL),
			s.account.PublicKey(),
			solana.SysVarRentPubkey,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instruction = append(instruction, initInst)
		signers = append(signers, tempAccount.PrivateKey)
		// if fromToken == NativeSOL {
		// 	fromAccount = tempAccount.PublicKey()
		// }
		// if toToken == NativeSOL {
		// 	toAccount = tempAccount.PublicKey()
		// }
	}

	instruction = append(instruction, NewRaydiumSwapInstruction(
		amount,
		minimumOutAmount,
		solana.TokenProgramID,
		pool.Id,
		pool.Authority,
		pool.OpenOrders,
		pool.TargetOrders,
		pool.BaseVault,
		pool.QuoteVault,
		pool.MarketProgramId,
		pool.MarketId,
		pool.MarketBids,
		pool.MarketAsks,
		pool.MarketEventQueue,
		pool.MarketBaseVault,
		pool.MarketQuoteVault,
		pool.MarketAuthority,
		tempAccount.PublicKey(),
		fromAccount,
		s.account.PublicKey(),
	))

	if needWrapSOL {
		closeInst, err := token.NewCloseAccountInstruction(
			tempAccount.PublicKey(),
			s.account.PublicKey(),
			s.account.PublicKey(),
			[]solana.PublicKey{},
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instruction = append(instruction, closeInst)
	}

	tx, err := ExecuteInstructions(ctx, s.clientRPC, signers, instruction...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *RaydiumSwap) WSOLSwap(
	ctx context.Context,
	pool RaydiumTokenData,
	amount uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	unitprice uint64,
	unitlimit uint32,
	needApprove bool,
	WSOLAssociated solana.PublicKey,

) (*solana.Transaction, error) {
	minimumOutAmount := uint64(0)

	var instruction []solana.Instruction
	signers := []solana.PrivateKey{s.account}
	// needWrapSOL := fromToken == NativeSOL || toToken == NativeSOL
	setFee, _ := computebudget.NewSetComputeUnitPriceInstruction(unitprice).ValidateAndBuild()
	setUnitLimit, _ := computebudget.NewSetComputeUnitLimitInstruction(unitlimit).ValidateAndBuild()
	instruction = append(instruction, setUnitLimit)
	instruction = append(instruction, setFee)
	if needApprove {
		instAprove, err := associatedTokenAccount.NewCreateInstruction(
			wallet.MainWallet().PublicKey(),
			wallet.MainWallet().PublicKey(),
			solana.MustPublicKeyFromBase58(toToken),
		).ValidateAndBuild()
		if err != nil {
			log.Fatalln(err.Error())
		}
		instruction = append(instruction, instAprove)
	}

	instruction = append(instruction, NewRaydiumSwapInstruction(
		amount,
		minimumOutAmount,
		solana.TokenProgramID,
		pool.Id,
		pool.Authority,
		pool.OpenOrders,
		pool.TargetOrders,
		pool.BaseVault,
		pool.QuoteVault,
		pool.MarketProgramId,
		pool.MarketId,
		pool.MarketBids,
		pool.MarketAsks,
		pool.MarketEventQueue,
		pool.MarketBaseVault,
		pool.MarketQuoteVault,
		pool.MarketAuthority,
		WSOLAssociated,
		fromAccount,
		s.account.PublicKey(),
	))

	tx, err := ExecuteInstructions(ctx, s.clientRPC, signers, instruction...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func ExecuteInstructions(
	ctx context.Context,
	clientRPC *rpc.Client,
	signers []solana.PrivateKey,
	instrs ...solana.Instruction,
) (*solana.Transaction, error) {

	tx, err := BuildTransacion(ctx, clientRPC, signers, instrs...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func BuildTransacion(ctx context.Context, clientRPC *rpc.Client, signers []solana.PrivateKey, instrs ...solana.Instruction) (*solana.Transaction, error) {
	recent, err := clientRPC.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		instrs,
		recent.Value.Blockhash,
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return nil, err
	}
	if viper.GetInt("TransactionVersion") == 1 {
		tx.Message.SetVersion(1)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, payer := range signers {
				if payer.PublicKey().Equals(key) {
					return &payer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

type RaySwapInstruction struct {
	bin.BaseVariant
	InAmount                uint64
	MinimumOutAmount        uint64
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (inst *RaySwapInstruction) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(utils.RAY_V4)
}

func (inst *RaySwapInstruction) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *RaySwapInstruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *RaySwapInstruction) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	err = encoder.WriteUint8(9)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.InAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.MinimumOutAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	return nil
}

func NewRaydiumSwapInstruction(
	// Parameters:
	inAmount uint64,
	minimumOutAmount uint64,
	// Accounts:
	tokenProgram solana.PublicKey,
	ammId solana.PublicKey,
	ammAuthority solana.PublicKey,
	ammOpenOrders solana.PublicKey,
	ammTargetOrders solana.PublicKey,
	poolCoinTokenAccount solana.PublicKey,
	poolPcTokenAccount solana.PublicKey,
	serumProgramId solana.PublicKey,
	serumMarket solana.PublicKey,
	serumBids solana.PublicKey,
	serumAsks solana.PublicKey,
	serumEventQueue solana.PublicKey,
	serumCoinVaultAccount solana.PublicKey,
	serumPcVaultAccount solana.PublicKey,
	serumVaultSigner solana.PublicKey,
	userSourceTokenAccount solana.PublicKey,
	userDestTokenAccount solana.PublicKey,
	userOwner solana.PublicKey,
) *RaySwapInstruction {

	inst := RaySwapInstruction{
		InAmount:         inAmount,
		MinimumOutAmount: minimumOutAmount,
		AccountMetaSlice: make(solana.AccountMetaSlice, 18),
	}
	inst.BaseVariant = bin.BaseVariant{
		Impl: inst,
	}

	inst.AccountMetaSlice[0] = solana.Meta(tokenProgram)
	inst.AccountMetaSlice[1] = solana.Meta(ammId).WRITE()
	inst.AccountMetaSlice[2] = solana.Meta(ammAuthority)
	inst.AccountMetaSlice[3] = solana.Meta(ammOpenOrders).WRITE()
	inst.AccountMetaSlice[4] = solana.Meta(ammTargetOrders).WRITE()
	inst.AccountMetaSlice[5] = solana.Meta(poolCoinTokenAccount).WRITE()
	inst.AccountMetaSlice[6] = solana.Meta(poolPcTokenAccount).WRITE()
	inst.AccountMetaSlice[7] = solana.Meta(serumProgramId)
	inst.AccountMetaSlice[8] = solana.Meta(serumMarket).WRITE()
	inst.AccountMetaSlice[9] = solana.Meta(serumBids).WRITE()
	inst.AccountMetaSlice[10] = solana.Meta(serumAsks).WRITE()
	inst.AccountMetaSlice[11] = solana.Meta(serumEventQueue).WRITE()
	inst.AccountMetaSlice[12] = solana.Meta(serumCoinVaultAccount).WRITE()
	inst.AccountMetaSlice[13] = solana.Meta(serumPcVaultAccount).WRITE()
	inst.AccountMetaSlice[14] = solana.Meta(serumVaultSigner)
	inst.AccountMetaSlice[15] = solana.Meta(userSourceTokenAccount).WRITE()
	inst.AccountMetaSlice[16] = solana.Meta(userDestTokenAccount).WRITE()
	inst.AccountMetaSlice[17] = solana.Meta(userOwner).SIGNER()

	return &inst
}

func (s *RaydiumSwap) WSOLSwap2(
	ctx context.Context,
	pool RaydiumTokenData,
	amount uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	unitprice uint64,
	unitlimit uint32,
	needApprove bool,
	WSOLAssociated solana.PublicKey,

) (*solana.Transaction, error) {
	minimumOutAmount := uint64(0)

	var instruction []solana.Instruction
	signers := []solana.PrivateKey{s.account}
	// needWrapSOL := fromToken == NativeSOL || toToken == NativeSOL
	setFee, _ := computebudget.NewSetComputeUnitPriceInstruction(unitprice).ValidateAndBuild()
	setUnitLimit, _ := computebudget.NewSetComputeUnitLimitInstruction(unitlimit).ValidateAndBuild()
	instruction = append(instruction, setUnitLimit)
	instruction = append(instruction, setFee)
	if needApprove {
		instAprove, err := associatedTokenAccount.NewCreateInstruction(
			wallet.SecondWallet().PublicKey(),
			wallet.SecondWallet().PublicKey(),
			solana.MustPublicKeyFromBase58(toToken),
		).ValidateAndBuild()
		if err != nil {
			log.Fatalln(err.Error())
		}
		instruction = append(instruction, instAprove)
	}

	instruction = append(instruction, NewRaydiumSwapInstruction(
		amount,
		minimumOutAmount,
		solana.TokenProgramID,
		pool.Id,
		pool.Authority,
		pool.OpenOrders,
		pool.TargetOrders,
		pool.BaseVault,
		pool.QuoteVault,
		pool.MarketProgramId,
		pool.MarketId,
		pool.MarketBids,
		pool.MarketAsks,
		pool.MarketEventQueue,
		pool.MarketBaseVault,
		pool.MarketQuoteVault,
		pool.MarketAuthority,
		WSOLAssociated,
		fromAccount,
		s.account.PublicKey(),
	))

	tx, err := ExecuteInstructions(ctx, s.clientRPC, signers, instruction...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
