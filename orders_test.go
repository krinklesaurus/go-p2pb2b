package p2pb2b

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": {
			"orderId": 25749,
			"market": "ETH_BTC",
			"price": "0.1",
			"side": "sell",
			"type": "limit",
			"timestamp": 1537535284.828868,
			"dealMoney": "0",
			"dealStock": "0",
			"amount": "0.1",
			"takerFee": "0.002",
			"makerFee": "0.002",
			"left": "0.1",
			"dealFee": "0"
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/order/new", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req CreateOrderRequest
		err := json.Unmarshal(reqBody, &req)

		assert.Nil(t, err, err)
		assert.Equal(t, "ETH_BTC", req.Market)
		assert.Equal(t, "buy", req.Side)
		assert.Equal(t, 0.001, req.Amount)
		assert.Equal(t, float64(1000), req.Price)
		assert.Equal(t, "/api/v1/order/new", req.Request)
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &CreateOrderRequest{
		Market: "ETH_BTC",
		Side:   "buy",
		Amount: 0.001,
		Price:  1000.0,
	}
	resp, err := client.CreateOrder(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}

func TestCancelOrder(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": {
			"orderId": 25749,
			"market": "ETH_BTC",
			"price": "0.1",
			"side": "sell",
			"type": "limit",
			"timestamp": 1537535284.828868,
			"dealMoney": "0",
			"dealStock": "0",
			"amount": "0.1",
			"takerFee": "0.002",
			"makerFee": "0.002",
			"left": "0.1",
			"dealFee": "0"
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/order/cancel", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req CancelOrderRequest
		err := json.Unmarshal(reqBody, &req)

		assert.Nil(t, err, err)
		assert.Equal(t, "ETH_BTC", req.Market)
		assert.Equal(t, 25749, int(req.OrderID))
		assert.Equal(t, req.Request, "/api/v1/order/cancel")
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &CancelOrderRequest{
		Market:  "ETH_BTC",
		OrderID: 25749,
	}
	resp, err := client.CancelOrder(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}

func TestQueryUnexecuted(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"

	body := `{
		"success": true,
		"message": "",
		"result": {
		  "limit": 100,
		  "offset": 0,
		  "total": 1,
		  "result": [
			{
			  "id": 3900714,
			  "left": "1",
			  "market": "ETH_BTC",
			  "amount": "1",
			  "type": "limit",
			  "price": "0.008",
			  "timestamp": 1546459568.376407,
			  "side": "buy",
			  "dealFee": "0",
			  "takerFee": "0.001",
			  "makerFee": "0.001",
			  "dealStock": "0",
			  "dealMoney": "0"
			}
		  ]
		}
	  }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/orders", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req QueryUnexecutedRequest
		err := json.Unmarshal(reqBody, &req)

		assert.Nil(t, err, err)
		assert.Equal(t, "ETH_BTC", req.Market)
		assert.Equal(t, int(0), req.Offset)
		assert.Equal(t, int(100), req.Limit)
		assert.Equal(t, "/api/v1/orders", req.Request)
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &QueryUnexecutedRequest{
		Market: "ETH_BTC",
		Offset: 0,
		Limit:  100,
	}
	resp, err := client.QueryUnexecuted(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}

func TestQueryExecuted(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"

	body := `{
		"success": true,
		"message": "",
		"result": {
		  "ETH_BTC": [
			{
			  "amount": "1",
			  "price": "0.01",
			  "type": "limit",
			  "id": 9740,
			  "side": "sell",
			  "ctime": 1533568890.583023,
			  "takerFee": "0.002",
			  "ftime": 1533630652.62185,
			  "market": "ETH_BTC",
			  "makerFee": "0.002",
			  "dealFee": "0.002",
			  "dealStock": "1",
			  "dealMoney": "0.01",
			  "marketName": "ETH_BTC"
			}
		  ],
		  "ATB_USD": [
			{
			  "amount": "0.3",
			  "price": "0.06296168",
			  "type": "market",
			  "id": 11669,
			  "side": "buy",
			  "ctime": 1533626329.696647,
			  "takerFee": "0.002",
			  "ftime": 1533626329.696659,
			  "market": "ATB_USD",
			  "makerFee": "0.002",
			  "dealFee": "0.000037777008",
			  "dealStock": "0.3",
			  "dealMoney": "0.018888504",
			  "marketName": "ATB_USD"
			}
		  ]
		}
	  }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/order_history", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req QueryExecutedRequest
		err := json.Unmarshal(reqBody, &req)

		assert.Nil(t, err, err)
		assert.Equal(t, 0, req.Offset)
		assert.Equal(t, 100, req.Limit)
		assert.Equal(t, "/api/v1/account/order_history", req.Request)
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &QueryExecutedRequest{
		Offset: 0,
		Limit:  100,
	}
	resp, err := client.QueryExecuted(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}

func TestQueryDeals(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": {
		  "offset": 0,
		  "limit": 50,
		  "records": [
			{
			  "time": 1533310924.935978,
			  "fee": "0",
			  "price": "80.22761599",
			  "amount": "2.12687945",
			  "id": 548,
			  "dealOrderId": 1237,
			  "role": 1,
			  "deal": "170.6344677716224"
			}
		  ]
		}
	  }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/order", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req QueryDealsRequest
		err := json.Unmarshal(reqBody, &req)

		assert.Nil(t, err, err)
		assert.Equal(t, 1234, req.OrderID)
		assert.Equal(t, 10, req.Offset)
		assert.Equal(t, 100, req.Limit)
		assert.Equal(t, "/api/v1/account/order", req.Request)
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &QueryDealsRequest{
		OrderID: 1234,
		Offset:  10,
		Limit:   100,
	}
	resp, err := client.QueryDeals(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}
