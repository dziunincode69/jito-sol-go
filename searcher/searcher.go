package searcher

import (
	"fmt"
	"jito_client/connection"
	"jito_client/lib/bundle"
	"jito_client/lib/packet"
	searcher "jito_client/lib/searcher"
	"jito_client/utils"
	"log"
	"strconv"
	"sync"
)

func ArrProgramMempoolSubscribe(AccountList []string) ([]searcher.SearcherService_SubscribeMempoolClient, error) {
	var wg sync.WaitGroup
	subscribersChan := make(chan searcher.SearcherService_SubscribeMempoolClient, len(AccountList)) // Create a buffered channel
	mempoolVar := &searcher.MempoolSubscription{
		Msg: &searcher.MempoolSubscription_ProgramV0Sub{
			ProgramV0Sub: &searcher.ProgramSubscriptionV0{
				Programs: AccountList,
			},

			// Msg: &searcher.MempoolSubscription_WlaV0Sub{
			// 	WlaV0Sub: &searcher.WriteLockedAccountSubscriptionV0{
			// 		Accounts: AccountList,
			// 	},
		},
	}

	for _, v := range connection.GetGRPCConnection() {
		wg.Add(1)
		go func(v searcher.SearcherServiceClient) {
			defer wg.Done()
			subs, err := v.SubscribeMempool(connection.GetContext(), mempoolVar)
			if err != nil {
				fmt.Println(err)
				return
			}
			subscribersChan <- subs // Send the subscriber to the channel
		}(v)
	}

	go func() {
		wg.Wait()
		close(subscribersChan) // Close the channel once all goroutines are done
	}()

	var subscribers []searcher.SearcherService_SubscribeMempoolClient
	for sub := range subscribersChan { // Collect all subscribers from the channel
		subscribers = append(subscribers, sub)
	}
	utils.LogWithTimestamp("Subscribers: " + strconv.Itoa(len(subscribers)))
	return subscribers, nil
}

func ProgramMempoolSubsribe(AccountList []string) (searcher.SearcherService_SubscribeMempoolClient, error) {
	mempoolVar := &searcher.MempoolSubscription{
		Msg: &searcher.MempoolSubscription_ProgramV0Sub{
			ProgramV0Sub: &searcher.ProgramSubscriptionV0{
				Programs: AccountList,
			},
			// Msg: &searcher.MempoolSubscription_WlaV0Sub{
			// 	WlaV0Sub: &searcher.WriteLockedAccountSubscriptionV0{
			// 		Accounts: AccountList,
			// 	},
		},
	}
	subs, err := connection.JITOSearcher().SubscribeMempool(connection.GetContext(), mempoolVar)
	if err != nil {
		return nil, err
	}
	return subs, nil
}
func GetTipAccount() (*searcher.GetTipAccountsResponse, error) {
	req, err := connection.JITOSearcher().GetTipAccounts(connection.GetContext(), &searcher.GetTipAccountsRequest{})
	if err != nil {
		return nil, err
	}
	return req, nil
}
func GetRegions() (*searcher.GetRegionsResponse, error) {
	req, err := connection.JITOSearcher().GetRegions(connection.GetContext(), &searcher.GetRegionsRequest{})
	if err != nil {
		return nil, err
	}
	return req, nil
}
func SendBundle(packets []*packet.Packet) {
	tx := &searcher.SendBundleRequest{
		Bundle: &bundle.Bundle{
			Packets: packets,
		},
	}
	req, err := connection.JITOSearcher().SendBundle(connection.GetContext(), tx)
	if err != nil {
		panic(err.Error() + " Error sending bundle")
	}
	log.Println("Success Sending Bundle:", req.Uuid)

}

func BundleResultSubscribe() {

	subs, err := connection.JITOSearcher().SubscribeBundleResults(connection.GetContext(), &searcher.SubscribeBundleResultsRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		msg, err := subs.Recv()
		if err != nil {
			fmt.Println(err)
			return
		}
		utils.LogWithTimestamp("Bundle Result: " + msg.String())
	}

}
