package p2pb2b

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	pseudoAPIKey := uuid.NewV4()
	pseudoAPISecret := "4a894c5c-8a7e-4337-bb6b-9fde16e3dddd"
	body := `{
		"success": true,
		"message": "",
		"result": [
			{
				"id": "ETH_BTC",
				"fromSymbol": "ETH",
				"toSymbol": "BTC"
			},
			{
				"id": "BTC_USD",
				"fromSymbol": "BTC",
				"toSymbol": "USD"
			}
		],
		"cache_time": 1574197000.65497,
		"current_time": 1574197000.655773
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/public/products", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		assert.Equal(t, pseudoAPIKey.String(), r.Header.Get("X-TXC-APIKEY"))
		assert.Empty(t, r.Header.Get("X-TXC-PAYLOAD"))
		assert.Empty(t, r.Header.Get("X-TXC-SIGNATURE"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, err := newClientWithURL(ts.URL, pseudoAPIKey.String(), pseudoAPISecret)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := client.GetProducts()
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, resp, fmt.Sprintf("error: %v", err))
	assert.Equal(t, true, resp.Success)
	assert.Equal(t, 1574197000.65497, resp.CacheTime)
	assert.Equal(t, 1574197000.655773, resp.CurrentTime)

	assert.Equal(t, 2, len(resp.Result))
	assert.Equal(t, "ETH_BTC", resp.Result[0].ID)
	assert.Equal(t, "ETH", resp.Result[0].FromSymbol)
	assert.Equal(t, "BTC", resp.Result[0].ToSymbol)

	assert.Equal(t, "BTC_USD", resp.Result[1].ID)
	assert.Equal(t, "BTC", resp.Result[1].FromSymbol)
	assert.Equal(t, "USD", resp.Result[1].ToSymbol)
}
