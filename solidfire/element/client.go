package element

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element/jsonrpc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.WithFields(logrus.Fields{
	"prefix": "main",
})

func init() {
	logrus.SetFormatter(new(prefixed.TextFormatter))
}

// A Client to interact with the Element API
type Client struct {
	Host                  string
	Username              string
	Password              string
	MaxConcurrentRequests int
	HTTPTransport         http.RoundTripper

	apiVersion string

	initOnce      sync.Once
	jsonrpcClient *jsonrpc.Client
	requestSlots  chan int
}

// CallAPIMethod can be used to make a request to any Element API method, receiving results as raw JSON
func (c *Client) CallAPIMethod(method string, params map[string]interface{}) (*json.RawMessage, error) {
	c.initOnce.Do(c.init)

	c.waitForAvailableSlot()
	defer c.releaseSlot()

	log.WithFields(logrus.Fields{
		"method": method,
		"params": params,
	}).Debug("Calling API")

	if params == nil {
		params = map[string]interface{}{}
	}
	result, err := c.jsonrpcClient.Do(&jsonrpc.Request{
		BaseURL: "/json-rpc/" + c.GetAPIVersion(),
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"method": method,
	}).Debug("Received successful API response")
	return result, nil
}

func (c *Client) init() {
	if c.MaxConcurrentRequests == 0 {
		c.MaxConcurrentRequests = 6
	}
	c.requestSlots = make(chan int, c.MaxConcurrentRequests)
	c.jsonrpcClient = &jsonrpc.Client{
		Host:          c.Host,
		Username:      c.Username,
		Password:      c.Password,
		HTTPTransport: c.HTTPTransport,
	}
}

// SetAPIVersion for the client to use for requests to the Element API
func (c *Client) SetAPIVersion(apiVersion string) {
	c.apiVersion = apiVersion
}

// GetAPIVersion returns the API version that will be used for Element API requests
func (c *Client) GetAPIVersion() string {
	if c.apiVersion == "" {
		return "1.0"
	}
	return c.apiVersion
}

func (c *Client) waitForAvailableSlot() {
	c.requestSlots <- 1
}

func (c *Client) releaseSlot() {
	<-c.requestSlots
}
