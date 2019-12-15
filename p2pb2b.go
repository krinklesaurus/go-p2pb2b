// Package p2pb2b is a library for accessing the P2pb2b API
package p2pb2b

import (
	"math"
	"net/http"
	"time"
)

// baseAPI is the p2pb2b API endpoint
const baseAPI = "https://api.p2pb2b.io"

// for testing purposes only
func newClientWithURL(url string, apiKey string, apiSecret string) (Client, error) {
	return &client{
		http: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		auth: &auth{
			APIKey:    apiKey,
			APISecret: apiSecret,
		},
		url: url,
	}, nil
}

// NewClient creates a new p2pb2b client with apiKey and apiSecret
func NewClient(apiKey string, apiSecret string) (Client, error) {
	return newClientWithURL(baseAPI, apiKey, apiSecret)
}

// Client is the basic p2pb2b client interface
type Client interface {
	PostCurrencyBalance(request *AccountCurrencyBalanceRequest) (*AccountCurrencyBalanceResp, error)
	PostBalances(request *AccountBalancesRequest) (*AccountBalancesResp, error)
	CreateOrder(request *CreateOrderRequest) (*CreateOrderResp, error)
	CancelOrder(request *CancelOrderRequest) (*CancelOrderResp, error)
	QueryUnexecuted(request *QueryUnexecutedRequest) (*QueryUnexecutedResp, error)
	QueryExecuted(request *QueryExecutedRequest) (*QueryExecutedResp, error)
	QueryDeals(request *QueryDealsRequest) (*QueryDealsResp, error)
	GetMarkets() (*MarketsResp, error)
	GetTickers() (*TickersResp, error)
	GetTicker(market string) (*TickerResp, error)
	GetOrderBook(market string, side string, offset int64, limit int64) (*OrderBookResp, error)
	GetHistory(market string, lastID int64, limit int64) (*HistoryResp, error)
	GetDepthResult(market string, limit int64) (*DepthResultResp, error)
	GetProducts() (*ProductsResp, error)
	GetSymbols() (*SymbolsResp, error)
}

// Response is the basic http response struct
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Request is the basic http request struct
type request struct {
	Request string `json:"request"`
	Nonce   int64  `json:"nonce"`
}

// TimestampToTime is a convenience function to convert a float64 timestamp to time.Time
func TimestampToTime(timestamp float64) time.Time {
	sec, dec := math.Modf(timestamp)
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
