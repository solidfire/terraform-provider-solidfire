package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// Response represents a Response to a JSON RPC API call
type Response struct {
	Result *json.RawMessage `json:"result"`
	Error  *ResponseError   `json:"error"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Request returned an error. %+v", *e)
}
