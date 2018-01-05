package solidfire

import (
	"crypto/tls"
	"net/http"

	"github.com/solidfire/terraform-provider-solidfire/solidfire/element"
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

func (c *Config) Client() (*element.Client, error) {
	client := &element.Client{
		Host:     "https://" + c.SolidFireServer,
		Username: c.User,
		Password: c.Password,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
		},
	}

	client.SetAPIVersion(c.APIVersion)

	return client, nil
}
