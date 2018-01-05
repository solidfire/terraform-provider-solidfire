package solidfire

import (
	"fmt"

	"github.com/solidfire/solidfire-sdk-golang/sfapi"
)

type Config struct {
	User            string
	Password        string
	SolidFireServer string
	APIVersion      string
}

type Client struct {
	Endpoint string
}

type APIError struct {
	Id    int `json:"id"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error"`
}

func (c *Config) Client() (*sfapi.Client, error) {
	client, err := sfapi.Create(c.SolidFireServer, c.User, c.Password, c.APIVersion, 443, 30)
	if err != nil {
		fmt.Printf("Received an error trying to create API client: %+v", err)
	}

	return client, nil
}
