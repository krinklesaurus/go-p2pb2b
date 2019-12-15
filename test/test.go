package main

import (
	"fmt"
	"os"
	"time"

	p2bp2b "github.com/krinklesaurus/go-p2pb2b"
)

const p2BP2BAPIKey = "P2BP2B_API_KEY"
const p2BP2BAPISecret = "P2BP2B_API_SECRET"

func main() {
	apiKey := os.Getenv(p2BP2BAPIKey)
	apiSecret := os.Getenv(p2BP2BAPISecret)

	if apiKey == "" || apiSecret == "" {
		fmt.Println(fmt.Sprintf("please provide env vars %s and %s", p2BP2BAPIKey, p2BP2BAPISecret))
		os.Exit(1)
	}

	client, _ := p2bp2b.NewClient(apiKey, apiSecret)

	fmt.Println("---- PostBalances ----")
	res, err := client.PostBalances(&p2bp2b.AccountBalancesRequest{})
	if err != nil {
		fmt.Println(fmt.Sprintf("error posting account balances, %v", err))
	}
	fmt.Println(fmt.Sprintf("res: %+v", res))
	time.Sleep(1 * time.Second)

	fmt.Println("---- PostCurrencyBalance ----")
	res2, err := client.PostCurrencyBalance(&p2bp2b.AccountCurrencyBalanceRequest{
		Currency: "PAX",
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("error posting currency balance for currency %s, %v", "ETH", err))
	}
	fmt.Println(fmt.Sprintf("res2: %+v", res2))
	time.Sleep(1 * time.Second)

	fmt.Println("---- Query Executed ----")
	res3, err := client.QueryExecuted(&p2bp2b.QueryExecutedRequest{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("error querying executed orders, %v", err))
	}
	fmt.Println(fmt.Sprintf("res3: %+v", res3))
	time.Sleep(1 * time.Second)

	fmt.Println("---- Query Unexecuted ----")
	res4, err := client.QueryUnexecuted(&p2bp2b.QueryUnexecutedRequest{
		Market: "ETH_BTC",
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("error querying deals, %v", err))
	}
	fmt.Println(fmt.Sprintf("res4: %+v", res4))
	time.Sleep(1 * time.Second)

	fmt.Println("---- Query Deals ----")
	res5, err := client.QueryDeals(&p2bp2b.QueryDealsRequest{
		OrderID: 12345,
		Limit:   100,
		Offset:  0,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("error querying deals, %v", err))
	}
	fmt.Println(fmt.Sprintf("res5: %+v", res5))
	time.Sleep(1 * time.Second)
}
