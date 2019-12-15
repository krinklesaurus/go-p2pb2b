package p2pb2b

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateOrderResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  Order  `json:"result"`
}

type Order struct {
	Amount    float64 `json:"amount,string"`
	DealFee   float64 `json:"dealFee,string"`
	DealMoney float64 `json:"dealMoney,string"`
	DealStock float64 `json:"dealStock,string"`
	Left      float64 `json:"left,string"`
	MakerFee  float64 `json:"makerFee,string"`
	Market    string  `json:"market"`
	OrderID   int     `json:"orderId"`
	Price     float64 `json:"price,string"`
	Side      string  `json:"side"`
	TakerFee  float64 `json:"takerFee,string"`
	Timestamp float64 `json:"timestamp"`
	Type      string  `json:"type"`
}

type CreateOrderRequest struct {
	request
	Market string  `json:"market"`
	Side   string  `json:"side"`
	Amount float64 `json:"amount,string"`
	Price  float64 `json:"price,string"`
}

type CancelOrderResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  Order  `json:"result"`
}

type CancelOrderRequest struct {
	request
	Market  string `json:"market"`
	OrderID int    `json:"orderId"`
}

type QueryUnexecutedRequest struct {
	request
	Market string `json:"market"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type QueryUnexecutedResp struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Result  QueryUnexecutedResult `json:"result,omitempty"`
}

type QueryUnexecutedResult struct {
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Total  int               `json:"total"`
	Result []UnexecutedOrder `json:"result"`
}

type UnexecutedOrder struct {
	Amount    float64 `json:"amount,string"`
	DealFee   float64 `json:"dealFee,string"`
	DealMoney float64 `json:"dealMoney,string"`
	DealStock float64 `json:"dealStock,string"`
	Left      float64 `json:"left,string"`
	MakerFee  float64 `json:"makerFee,string"`
	Market    string  `json:"market"`
	ID        int     `json:"id"`
	Price     float64 `json:"price,string"`
	Side      string  `json:"side"`
	TakerFee  float64 `json:"takerFee,string"`
	Timestamp float64 `json:"timestamp"`
	Type      string  `json:"type"`
}

type QueryExecutedRequest struct {
	request
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type QueryExecutedResp struct {
	Response
	Result map[string][]AltOrder `json:"result,omitempty"`
}

type AltOrder struct {
	Amount     float64 `json:"amount,string"`
	Price      float64 `json:"price,string"`
	Type       string  `json:"type"`
	ID         int     `json:"id"`
	Source     string  `json:"source,omitempty"`
	Side       string  `json:"side"`
	Ctime      float64 `json:"ctime"`
	TakerFee   float64 `json:"takerFee,string"`
	Ftime      float64 `json:"ftime"`
	Market     string  `json:"market"`
	MakerFee   float64 `json:"makerFee,string"`
	DealFee    float64 `json:"dealFee,string"`
	DealStock  float64 `json:"dealStock,string"`
	DealMoney  float64 `json:"dealMoney,string"`
	MarketName string  `json:"marketName"`
}

type QueryDealsRequest struct {
	request
	OrderID int `json:"orderId"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
}

type QueryDealsResp struct {
	Response
	Result QueryDealsResult `json:"result,omitempty"`
}

type QueryDealsResult struct {
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Records []Record `json:"records"`
}

type Record struct {
	Time        float64 `json:"time"`
	Fee         float64 `json:"fee,string"`
	Price       float64 `json:"price,string"`
	Amount      float64 `json:"amount,string"`
	ID          int     `json:"id"`
	DealOrderID int     `json:"dealOrderId"`
	Role        int     `json:"role"`
	Deal        float64 `json:"deal,string"`
}

func (c *client) CreateOrder(req *CreateOrderRequest) (*CreateOrderResp, error) {
	path := "/api/v1/order/new"
	url := fmt.Sprintf("%s%s", c.url, path)
	req.request = request{
		Nonce:   time.Now().UnixNano(),
		Request: path,
	}
	asJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendPost(url, nil, bytes.NewReader(asJSON))
	if err != nil {
		return nil, err
	}
	err = checkHTTPStatus(*resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CreateOrderResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) CancelOrder(req *CancelOrderRequest) (*CancelOrderResp, error) {
	path := "/api/v1/order/cancel"
	url := fmt.Sprintf("%s%s", c.url, path)
	req.request = request{
		Nonce:   time.Now().UnixNano(),
		Request: path,
	}
	asJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendPost(url, nil, bytes.NewReader(asJSON))
	if err != nil {
		return nil, err
	}
	err = checkHTTPStatus(*resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CancelOrderResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) QueryUnexecuted(req *QueryUnexecutedRequest) (*QueryUnexecutedResp, error) {
	path := "/api/v1/orders"
	url := fmt.Sprintf("%s%s", c.url, path)
	req.request = request{
		Nonce:   time.Now().UnixNano(),
		Request: path,
	}
	asJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendPost(url, nil, bytes.NewReader(asJSON))
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = checkHTTPStatus(*resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result QueryUnexecutedResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) QueryExecuted(req *QueryExecutedRequest) (*QueryExecutedResp, error) {
	path := "/api/v1/account/order_history"
	url := fmt.Sprintf("%s%s", c.url, path)
	req.request = request{
		Nonce:   time.Now().UnixNano(),
		Request: path,
	}
	asJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendPost(url, nil, bytes.NewReader(asJSON))
	if err != nil {
		return nil, err
	}

	err = checkHTTPStatus(*resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result QueryExecutedResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) QueryDeals(req *QueryDealsRequest) (*QueryDealsResp, error) {
	path := "/api/v1/account/order"
	url := fmt.Sprintf("%s%s", c.url, path)
	req.request = request{
		Nonce:   time.Now().UnixNano(),
		Request: path,
	}
	asJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendPost(url, nil, bytes.NewReader(asJSON))
	if err != nil {
		return nil, err
	}

	err = checkHTTPStatus(*resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result QueryDealsResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
