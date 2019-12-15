package p2pb2b

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// HeaderXTxcAPIKey is the HTTP Header for account API key
	HeaderXTxcAPIKey = "X-TXC-APIKEY"
	// HeaderXTxcPayloard is the HTTP Header for body json encoded in base64
	HeaderXTxcPayloard = "X-TXC-PAYLOAD"
	// HeaderXTxcSignature is the HTTP Header for encrypted in HmacSHA512 payload
	HeaderXTxcSignature = "X-TXC-SIGNATURE"
)

type auth struct {
	APIKey    string
	APISecret string
}

type client struct {
	http *http.Client
	auth *auth
	url  string
}

type response struct {
	Header     http.Header
	Body       io.ReadCloser
	StatusCode int
	Status     string
}

func checkHTTPStatus(resp response, expected ...int) error {
	for _, e := range expected {
		if resp.StatusCode == e {
			return nil
		}
	}
	return fmt.Errorf("http response status != %+v, got %d", expected, resp.StatusCode)
}

func mergeHeaders(firstHeaders map[string]string, secondHeaders map[string]string) map[string]string {
	if secondHeaders == nil {
		return firstHeaders
	}
	if firstHeaders == nil {
		return secondHeaders
	}
	for k, v := range secondHeaders {
		if firstHeaders[k] == "" {
			firstHeaders[k] = v
		}
	}
	return firstHeaders
}

func (c *client) sendPost(url string, additionalHeaders map[string]string, body io.Reader) (*response, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return &response{}, fmt.Errorf("error creating POST request, %v", err)
	}

	if additionalHeaders == nil {
		additionalHeaders = make(map[string]string)
	}
	bodyBase64 := base64.StdEncoding.EncodeToString(bodyBytes)
	additionalHeaders[HeaderXTxcPayloard] = bodyBase64

	if c.auth != nil {
		h := hmac.New(sha512.New, []byte(c.auth.APISecret))
		h.Write([]byte(bodyBase64))
		signature := hex.EncodeToString(h.Sum(nil))
		additionalHeaders[HeaderXTxcSignature] = signature
	}

	return c.sendRequest(req, additionalHeaders)
}

func (c *client) sendGet(url string, additionalHeaders map[string]string) (*response, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &response{}, fmt.Errorf("error creating GET request, %v", err)
	}

	return c.sendRequest(req, additionalHeaders)
}

func (c *client) sendRequest(request *http.Request, additionalHeaders map[string]string) (*response, error) {

	for k, v := range additionalHeaders {
		request.Header.Add(k, v)
	}

	thisHeaders := map[string]string{}
	thisHeaders["Content-Type"] = "application/json"
	if c.auth != nil {
		thisHeaders[HeaderXTxcAPIKey] = c.auth.APIKey
	}
	headers := mergeHeaders(additionalHeaders, thisHeaders)
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	return &response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header,
		Body:       resp.Body,
	}, nil
}
