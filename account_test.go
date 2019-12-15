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

func TestPostBalancesNoKeyProvided(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/balances", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"success\":false,\"message\":[[\"Key not provided.\"]],\"result\":[]}"))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &AccountBalancesRequest{}
	_, err = client.PostBalances(request)
	assert.True(t, err != nil)
}

func TestPostBalances(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": {
		"ETH": {
			"available": "0.1",
			"freeze": "0.4"
		}
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/balances", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req AccountBalancesRequest
		err := json.Unmarshal(reqBody, &req)
		assert.Nil(t, err, err)
		assert.Equal(t, req.Request, "/api/v1/account/balances")
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &AccountBalancesRequest{}
	resp, err := client.PostBalances(request)

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}

func TestPostCurrencyBalanceNoKeyProvided(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/balance", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"success\":false,\"message\":[[\"Key not provided.\"]],\"result\":[]}"))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &AccountCurrencyBalanceRequest{
		Currency: "ETH",
	}
	_, err = client.PostCurrencyBalance(request)
	assert.True(t, err != nil)
}

func TestPostCurrencyBalance(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": {
		  "ETH": {
			"available": "0.63",
			"freeze": "0"
		  }
		}
	  }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/account/balance", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		reqBody, _ := ioutil.ReadAll(r.Body)
		reqBody64 := base64.StdEncoding.EncodeToString(reqBody)

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Equal(t, reqBody64, r.Header.Get("X-TXC-PAYLOAD"))
		assert.NotEmpty(t, r.Header.Get("X-TXC-SIGNATURE"))

		var req AccountCurrencyBalanceRequest
		err := json.Unmarshal(reqBody, &req)
		assert.Nil(t, err, err)
		assert.Equal(t, req.Currency, "ETH")
		assert.Equal(t, req.Request, "/api/v1/account/balance")
		assert.NotEmpty(t, req.Nonce)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}
	request := &AccountCurrencyBalanceRequest{
		Currency: "ETH",
	}
	resp, err := client.PostCurrencyBalance(request)
	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)

	respBytes, _ := json.Marshal(resp)
	equal, _ := isEqualJSON(body, string(respBytes))
	assert.True(t, equal, fmt.Sprintf("%s is not equal to %s", body, string(respBytes)))
}
