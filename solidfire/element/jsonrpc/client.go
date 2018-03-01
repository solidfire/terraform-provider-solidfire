package jsonrpc

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
)

// Client represents a client for interaction with a JSON RPC API
type Client struct {
	Host          string
	Username      string
	Password      string
	HTTPTransport http.RoundTripper

	initOnce   sync.Once
	httpClient http.Client
}

func (c *Client) init() {
	if c.HTTPTransport != nil {
		c.httpClient.Transport = c.HTTPTransport
	}
}

// Do sends the API Request, parses the response as JSON, and returns the "result" value as raw JSON
func (c *Client) Do(req *Request) (*json.RawMessage, error) {
	c.initOnce.Do(c.init)

	httpReq, err := req.BuildHTTPReq(c.Host, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Print("HTTP req failed")
		return nil, err
	}

	if httpRes.StatusCode == 401 {
		return nil, errors.New("Unauthenticated")
	}

	defer httpRes.Body.Close()
	var res Response

	if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		log.Print("HTTP decoder failed")
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	if res.Result == nil {
		return nil, errors.New("No result returned in JSON RPC response.")
	}

	return res.Result, nil
}
