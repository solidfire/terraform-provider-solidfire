package jsonrpc

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestSuccessfulCall(t *testing.T) {
	defer gock.Off()

	fakeHost := "http://fakehost"
	fakeUsername := "user"
	fakePassword := "pass"

	gock.New(fakeHost).
		Post("/path").
		// $ printf user:pass | base64
		// dXNlcjpwYXNz
		MatchHeader("Authorization", "Basic dXNlcjpwYXNz").
		JSON(map[string]interface{}{
			"method": "ArbitraryMethodName",
			"params": map[string]interface{}{
				"fakeParam1": "fakeParam1Value",
				"fakeParam2": "fakeParam2Value",
			},
		}).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"result": map[string]interface{}{
				"fakeResult1": 12,
				"fakeResultObj": map[string]interface{}{
					"nestedFakeVal": "test",
				},
			},
		})

	client := &Client{
		Host:     fakeHost,
		Username: fakeUsername,
		Password: fakePassword,
	}

	result, err := client.Do(&Request{
		BaseURL: "/path",
		Method:  "ArbitraryMethodName",
		Params: map[string]interface{}{
			"fakeParam1": "fakeParam1Value",
			"fakeParam2": "fakeParam2Value",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, result, &json.RawMessage{})
	assert.NotNil(t, result)
}

func TestMalformedResponse(t *testing.T) {
	defer gock.Off()

	fakeHost := "http://fakehost"
	fakeUsername := "user"
	fakePassword := "pass"

	gock.New(fakeHost).
		Post("/path").
		JSON(map[string]interface{}{
			"method": "ArbitraryMethodName",
			"params": map[string]interface{}{
				"fakeParam1": "fakeParam1Value",
				"fakeParam2": "fakeParam2Value",
			},
		}).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"response-with-no-result": map[string]interface{}{},
		})

	client := &Client{
		Host:     fakeHost,
		Username: fakeUsername,
		Password: fakePassword,
	}

	_, err := client.Do(&Request{
		BaseURL: "/path",
		Method:  "ArbitraryMethodName",
		Params: map[string]interface{}{
			"fakeParam1": "fakeParam1Value",
			"fakeParam2": "fakeParam2Value",
		},
	})
	if err == nil {
		t.Fatal("Expected error to be returned")
	}
}
