package p2pb2b

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AccountBalancesResp struct {
	Response
	Result map[string]AccountBalance `json:"result"`
}

type AccountBalance struct {
	Available float64 `json:"available,string"`
	Freeze    float64 `json:"freeze,string"`
}

type AccountBalancesRequest struct {
	request
}

type AccountCurrencyBalanceResp struct {
	Response
	Result map[string]AccountCurrencyBalance `json:"result,omitempty"`
}

type AccountCurrencyBalance struct {
	Available float64 `json:"available,string"`
	Freeze    float64 `json:"freeze,string"`
}

type AccountCurrencyBalanceRequest struct {
	Currency string `json:"currency"`
	request
}

func (c *client) PostBalances(req *AccountBalancesRequest) (*AccountBalancesResp, error) {
	path := "/api/v1/account/balances"
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

	var result AccountBalancesResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) PostCurrencyBalance(req *AccountCurrencyBalanceRequest) (*AccountCurrencyBalanceResp, error) {
	path := "/api/v1/account/balance"
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
		return nil, fmt.Errorf("error: %v", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AccountCurrencyBalanceResp
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
