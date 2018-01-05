package element

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"strconv"
)

type ListVolumesRequest struct {
	Volumes               []int `structs:"volumeIDs"`
	IncludeVirtualVolumes bool  `structs:"includeVirtualVolumes"`
}

type ListVolumesResult struct {
	Volumes []Volume `json:"volumes"`
}

type Volume struct {
	Name     string `json:"name"`
	VolumeID int    `json:"volumeID"`
	Iqn      string `json:"iqn"`
}

func (c *Client) GetVolumeByID(id string) (Volume, error) {
	convID, err := strconv.Atoi(id)
	if err != nil {
		return Volume{}, err
	}

	volIDs := make([]int, 1)
	volIDs[0] = convID

	params := structs.Map(ListVolumesRequest{Volumes: volIDs})

	response, err := c.CallAPIMethod("ListVolumes", params)
	if err != nil {
		log.Print("ListVolumes request failed")
		return Volume{}, err
	}

	var result ListVolumesResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal response from ListVolumes")
		return Volume{}, err
	}

	if len(result.Volumes) != 1 {
		return Volume{}, errors.New(fmt.Sprintf("Expected one Volume to be found. Response contained %v results", len(result.Volumes)))
	}

	return result.Volumes[0], nil
}
