package solidfire

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element/jsonrpc"
)

type CreateVolumeAccessGroupRequest struct {
	Name       string      `structs:"name"`
	Initiators []string    `structs:"initiators"`
	Volumes    []int       `structs:"volumes"`
	Attributes interface{} `structs:"attributes"`
	ID         int         `structs:"id"`
}

type CreateVolumeAccessGroupResult struct {
	VolumeAccessGroupID int `json:"volumeAccessGroupID"`
	element.VolumeAccessGroup
}

type DeleteVolumeAccessGroupRequest struct {
	VolumeAccessGroupID    int  `structs:"volumeAccessGroupID"`
	DeleteOrphanInitiators bool `structs:"deleteOrphanInitiators"`
	Force                  bool `structs:"force"`
}

type ModifyVolumeAccessGroupRequest struct {
	VolumeAccessGroupID    int         `structs:"volumeAccessGroupID"`
	Name                   string      `structs:"name"`
	Attributes             interface{} `structs:"attributes"`
	Initiators             []int       `structs:"initiators"`
	DeleteOrphanInitiators bool        `structs:"deleteOrphanInitiators"`
	Volumes                []int       `structs:"volumes"`
}

func resourceSolidFireVolumeAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolidFireVolumeAccessGroupCreate,
		Read:   resourceSolidFireVolumeAccessGroupRead,
		Update: resourceSolidFireVolumeAccessGroupUpdate,
		Delete: resourceSolidFireVolumeAccessGroupDelete,
		Exists: resourceSolidFireVolumeAccessGroupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"initiators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSolidFireVolumeAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume access group: %#v", d)
	client := meta.(*element.Client)

	vag := CreateVolumeAccessGroupRequest{}

	if v, ok := d.GetOk("name"); ok {
		vag.Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			vag.Volumes = append(vag.Volumes, v.(int))
		}
	}

	resp, err := createVolumeAccessGroup(client, vag)
	if err != nil {
		log.Print("Error creating volume access group")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.VolumeAccessGroupID))
	log.Printf("Created volume access group: %v %v", vag.Name, resp.VolumeAccessGroupID)

	return resourceSolidFireVolumeAccessGroupRead(d, meta)
}

func createVolumeAccessGroup(client *element.Client, request CreateVolumeAccessGroupRequest) (CreateVolumeAccessGroupResult, error) {
	params := structs.Map(request)

	log.Printf("Parameters: %v", params)

	response, err := client.CallAPIMethod("CreateVolumeAccessGroup", params)
	if err != nil {
		log.Print("CreateVolumeAccessGroup request failed")
		return CreateVolumeAccessGroupResult{}, err
	}

	var result CreateVolumeAccessGroupResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from CreateVolumeAccessGroup")
		return CreateVolumeAccessGroupResult{}, err
	}
	return result, nil
}

func resourceSolidFireVolumeAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume access group: %#v", d)
	client := meta.(*element.Client)

	vags := element.ListVolumeAccessGroupsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	vags.VolumeAccessGroups = s

	res, err := listVolumeAccessGroups(client, vags)
	if err != nil {
		return err
	}

	if len(res.VolumeAccessGroupsNotFound) > 0 {
		return fmt.Errorf("Unable to find Volume Access Groups with the ID of %v", res.VolumeAccessGroupsNotFound)
	}

	if len(res.VolumeAccessGroups) != 1 {
		return fmt.Errorf("Expected one Volume Access Group to be found. Response contained %v results", len(res.VolumeAccessGroups))
	}

	d.Set("name", res.VolumeAccessGroups[0].Name)
	d.Set("initiators", res.VolumeAccessGroups[0].Initiators)
	d.Set("volumes", res.VolumeAccessGroups[0].Volumes)

	return nil
}

func listVolumeAccessGroups(client *element.Client, request element.ListVolumeAccessGroupsRequest) (element.ListVolumeAccessGroupsResult, error) {
	params := structs.Map(request)

	response, err := client.CallAPIMethod("ListVolumeAccessGroups", params)
	if err != nil {
		log.Print("ListVolumeAccessGroups request failed")
		return element.ListVolumeAccessGroupsResult{}, err
	}

	var result element.ListVolumeAccessGroupsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumeAccessGroups")
		return element.ListVolumeAccessGroupsResult{}, err
	}

	return result, nil
}

func resourceSolidFireVolumeAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume access group %#v", d)
	client := meta.(*element.Client)

	vag := ModifyVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	vag.VolumeAccessGroupID = convID

	if v, ok := d.GetOk("name"); ok {
		vag.Name = v.(string)

	} else {
		return fmt.Errorf("name argument is required during update")
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			vag.Volumes = append(vag.Volumes, v.(int))
		}
	} else {
		return fmt.Errorf("expecting an array of volume ids to change")
	}

	err := modifyVolumeAccessGroup(client, vag)
	if err != nil {
		return err
	}

	return nil
}

func modifyVolumeAccessGroup(client *element.Client, request ModifyVolumeAccessGroupRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyVolumeAccessGroup", params)
	if err != nil {
		log.Print("ModifyVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume access group: %#v", d)
	client := meta.(*element.Client)

	vag := DeleteVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	vag.VolumeAccessGroupID = convID

	err := deleteVolumeAccessGroup(client, vag)
	if err != nil {
		return err
	}

	return nil
}

func deleteVolumeAccessGroup(client *element.Client, request DeleteVolumeAccessGroupRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("DeleteVolumeAccessGroup", params)
	if err != nil {
		log.Print("DeleteVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeAccessGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume access group: %#v", d)
	client := meta.(*element.Client)

	vags := element.ListVolumeAccessGroupsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	vags.VolumeAccessGroups = s

	res, err := listVolumeAccessGroups(client, vags)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.VolumeAccessGroupsNotFound) > 0 {
		d.SetId("")
		return false, nil
	}

	if len(res.VolumeAccessGroups) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
