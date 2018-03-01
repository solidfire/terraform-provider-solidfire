package jsonrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// Request represents a request to a JSON RPC API
type Request struct {
	BaseURL string      `json:"-"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// BuildHTTPReq builds an HTTP request to carry out the JSON RPC request
func (r *Request) BuildHTTPReq(host, user, pass string) (*http.Request, error) {
	bodyJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	url := host + r.BaseURL
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	basicAuth := fmt.Sprintf("%s:%s", user, pass)
	basicAuthEncoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))
	req.Header.Set("Authorization", "Basic "+basicAuthEncoded)

	return req, nil
}
