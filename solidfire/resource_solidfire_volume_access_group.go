package solidfire

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/solidfire/solidfire-sdk-golang/sfapi"
	"github.com/solidfire/solidfire-sdk-golang/sftypes"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element/jsonrpc"
)

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
	client := meta.(*sfapi.Client)

	vag := sftypes.CreateVolumeAccessGroupRequest{}

	if v, ok := d.GetOk("name"); ok {
		vag.Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			vag.Volumes = append(vag.Volumes, int64(v.(int)))
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

func createVolumeAccessGroup(client *sfapi.Client, request sftypes.CreateVolumeAccessGroupRequest) (*sftypes.CreateVolumeAccessGroupResult, error) {
	response, err := client.CreateVolumeAccessGroup(request)
	if err != nil {
		log.Print("CreateVolumeAccessGroup request failed")
		return &sftypes.CreateVolumeAccessGroupResult{}, err
	}
	return response, nil
}

func resourceSolidFireVolumeAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume access group: %#v", d)
	client := meta.(*sfapi.Client)

	vags := sftypes.ListVolumeAccessGroupsRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	vags.StartVolumeAccessGroupID = convID

	res, err := listVolumeAccessGroups(client, vags)
	if err != nil {
		return err
	}

	if len(res.VolumeAccessGroups) != 1 {
		return fmt.Errorf("Expected one Volume Access Group to be found. Response contained %v results", len(res.VolumeAccessGroups))
	}

	d.Set("name", res.VolumeAccessGroups[0].Name)
	d.Set("initiators", res.VolumeAccessGroups[0].Initiators)
	d.Set("volumes", res.VolumeAccessGroups[0].Volumes)

	return nil
}

func listVolumeAccessGroups(client *sfapi.Client, request sftypes.ListVolumeAccessGroupsRequest) (*sftypes.ListVolumeAccessGroupsResult, error) {
	response, err := client.ListVolumeAccessGroups(request)
	if err != nil {
		log.Print("ListVolumeAccessGroups request failed")
		return &sftypes.ListVolumeAccessGroupsResult{}, err
	}

	return response, nil
}

func resourceSolidFireVolumeAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume access group %#v", d)
	client := meta.(*sfapi.Client)

	vag := sftypes.ModifyVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

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
			vag.Volumes = append(vag.Volumes, int64(v.(int)))
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

func modifyVolumeAccessGroup(client *sfapi.Client, request sftypes.ModifyVolumeAccessGroupRequest) error {
	err := client.ModifyVolumeAccessGroup(request)
	if err != nil {
		log.Print("ModifyVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume access group: %#v", d)
	client := meta.(*sfapi.Client)

	vag := sftypes.DeleteVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

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

func deleteVolumeAccessGroup(client *sfapi.Client, request sftypes.DeleteVolumeAccessGroupRequest) error {
	err := client.DeleteVolumeAccessGroup(request)
	if err != nil {
		log.Print("DeleteVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeAccessGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume access group: %#v", d)
	client := meta.(*sfapi.Client)

	vags := sftypes.ListVolumeAccessGroupsRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	vags.StartVolumeAccessGroupID = convID

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

	if len(res.VolumeAccessGroups) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
