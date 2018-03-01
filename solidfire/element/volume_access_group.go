package element

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"strconv"
)

type ListVolumeAccessGroupsRequest struct {
	VolumeAccessGroups []int `structs:"volumeAccessGroups"`
}

type ListVolumeAccessGroupsResult struct {
	VolumeAccessGroups         []VolumeAccessGroup `json:"volumeAccessGroups"`
	VolumeAccessGroupsNotFound []int               `json:"volumeAccessGroupsNotFound"`
}

type VolumeAccessGroup struct {
	VolumeAccessGroupID int      `json:"volumeAccessGroupID"`
	Name                string   `json:"name"`
	Initiators          []string `json:"initiators"`
	Volumes             []int    `json:"volumes"`
	ID                  int      `json:"id"`
}

func (c *Client) GetVolumeAccessGroupByID(id string) (VolumeAccessGroup, error) {
	convID, err := strconv.Atoi(id)
	if err != nil {
		return VolumeAccessGroup{}, err
	}

	vagIDs := make([]int, 1)
	vagIDs[0] = convID

	params := structs.Map(ListVolumeAccessGroupsRequest{VolumeAccessGroups: vagIDs})

	response, err := c.CallAPIMethod("ListVolumeAccessGroups", params)
	if err != nil {
		log.Print("ListVolumeAccessGroups request failed")
		return VolumeAccessGroup{}, err
	}

	var result ListVolumeAccessGroupsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal respone from ListVolumeAccessGroups")
		return VolumeAccessGroup{}, err
	}

	if len(result.VolumeAccessGroupsNotFound) > 0 {
		return VolumeAccessGroup{}, errors.New(fmt.Sprintf("Unable to find Volume Access Groups with the ID of %v", result.VolumeAccessGroupsNotFound))
	}

	if len(result.VolumeAccessGroups) != 1 {
		return VolumeAccessGroup{}, errors.New(fmt.Sprintf("Expected one Volume Access Group to be found. Response contained %v results", len(result.VolumeAccessGroups)))
	}

	return result.VolumeAccessGroups[0], nil
}
