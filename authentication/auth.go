package authentication

import (
	"context"
	"crypto/ed25519"
	"jito_client/connection"
	"jito_client/lib/auth"
	"jito_client/wallet"
)

type AuthAccess struct {
	Access_token  string
	Refresh_token string
}

func GetChallenge() (string, error) {
	req := &auth.GenerateAuthChallengeRequest{
		Role:   auth.Role_SEARCHER,
		Pubkey: wallet.TipWallet().PublicKey().Bytes(),
	}
	tmp_chall, err := connection.JITOAuth().GenerateAuthChallenge(context.Background(), req)
	if err != nil {
		return "", err
	}
	return tmp_chall.GetChallenge(), nil
}
func GetAuthToken(challenge string) (AuthAccess, error) {
	datatosign := wallet.TipWallet().PublicKey().String() + "-" + challenge
	signedChallenge := ed25519.Sign(ed25519.PrivateKey(wallet.TipWallet()), []byte(datatosign))
	req, err := connection.AuthService.GenerateAuthTokens(context.Background(), &auth.GenerateAuthTokensRequest{
		Challenge:       datatosign,
		ClientPubkey:    wallet.TipWallet().PublicKey().Bytes(),
		SignedChallenge: signedChallenge,
	})
	if err != nil {
		return AuthAccess{}, err
	}
	return AuthAccess{
		Access_token:  req.GetAccessToken().Value,
		Refresh_token: req.GetRefreshToken().Value,
	}, nil
}

func RefreshToken(refresh_token string) (AuthAccess, error) {
	req, err := connection.AuthService.RefreshAccessToken(context.Background(), &auth.RefreshAccessTokenRequest{
		RefreshToken: refresh_token,
	})
	if err != nil {
		return AuthAccess{}, err
	}
	return AuthAccess{
		Access_token: req.GetAccessToken().Value,
	}, nil
}
