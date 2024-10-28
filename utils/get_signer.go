package utils

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
)

func GetSigner(serumPcVaultAccount string) string {
	client := &http.Client{}
	var data = strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"getAccountInfo","params":["` + serumPcVaultAccount + `",{"encoding":"jsonParsed"}]}`)
	req, err := http.NewRequest("POST", "https://api.mainnet-beta.solana.com", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	vaultSigner, _ := jsonparser.GetString(bodyText, "result", "value", "data", "parsed", "info", "owner")
	return vaultSigner
}
