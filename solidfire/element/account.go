package element

import (
	"encoding/json"
	"github.com/fatih/structs"
)

type GetAccountByIDRequest struct {
	AccountID int `structs:"accountID"`
}

type GetAccountByIDResult struct {
	Account Account `json:"account"`
}

type Account struct {
	AccountID       int         `json:"accountID"`
	Attributes      interface{} `json:"attributes"`
	InitiatorSecret string      `json:"initiatorSecret"`
	Status          string      `json:"status"`
	TargetSecret    string      `json:"targetSecret"`
	Username        string      `json:"username"`
}

func (c *Client) GetAccountByID(id int) (Account, error) {
	params := structs.Map(GetAccountByIDRequest{AccountID: id})

	response, err := c.CallAPIMethod("GetAccountByID", params)
	if err != nil {
		log.Print("GetAccountByID request failed")
		return Account{}, err
	}

	var result GetAccountByIDResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal response from GetAccountByID")
		return Account{}, err
	}

	return result.Account, nil
}
