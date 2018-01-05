package element

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"strconv"
)

type ListInitiatorRequest struct {
	Initiators []int `structs:"initiators"`
}

type ListInitiatorResult struct {
	Initiators []InitiatorResponse `json:"initiators"`
}

type Initiator struct {
	Name                string      `structs:"name,omitempty"`
	Alias               string      `structs:"alias,omitempty"`
	Attributes          interface{} `structs:"attributes,omitempty"`
	VolumeAccessGroupID int         `structs:"volumeAccessGroupID,omitempty"`
	InitiatorID         int         `structs:"initiatorID,omitempty"`
}

type InitiatorResponse struct {
	Name               string      `json:"initiatorName"`
	Alias              string      `json:"alias"`
	Attributes         interface{} `json:"attributes"`
	ID                 int         `json:"initiatorID"`
	VolumeAccessGroups []int       `json:"volumeAccessGroups"`
}

func (c *Client) GetInitiatorByID(id string) (Initiator, error) {
	convID, err := strconv.Atoi(id)
	if err != nil {
		return Initiator{}, err
	}

	initID := make([]int, 1)
	initID[0] = convID

	params := structs.Map(ListInitiatorRequest{Initiators: initID})

	response, err := c.CallAPIMethod("ListInitiators", params)
	if err != nil {
		log.Print("ListInitiators request failed")
		return Initiator{}, err
	}

	var result ListInitiatorResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal response from ListInitiators")
		return Initiator{}, err
	}

	if len(result.Initiators) != 1 {
		return Initiator{}, errors.New(fmt.Sprintf("Expected one Initiator to be found. Response contained %v results", len(result.Initiators)))
	}

	var initiator Initiator
	initiator.Name = result.Initiators[0].Name
	initiator.Alias = result.Initiators[0].Alias
	initiator.Attributes = result.Initiators[0].Attributes
	initiator.InitiatorID = result.Initiators[0].ID
	if len(result.Initiators[0].VolumeAccessGroups) == 1 {
		initiator.VolumeAccessGroupID = result.Initiators[0].VolumeAccessGroups[0]
	}

	return initiator, nil
}
