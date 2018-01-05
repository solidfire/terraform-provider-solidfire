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

func resourceSolidFireVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolidFireVolumeCreate,
		Read:   resourceSolidFireVolumeRead,
		Update: resourceSolidFireVolumeUpdate,
		Delete: resourceSolidFireVolumeDelete,
		Exists: resourceSolidFireVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"total_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable512e": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"min_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"burst_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSolidFireVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)
	client := meta.(*sfapi.Client)

	volume := sftypes.CreateVolumeRequest{}

	if v, ok := d.GetOk("name"); ok {
		volume.Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("account_id"); ok {
		volume.AccountID = int64(v.(int))
	} else {
		return fmt.Errorf("account_id argument is required")
	}

	if v, ok := d.GetOk("total_size"); ok {
		volume.TotalSize = int64(v.(int))
	} else {
		return fmt.Errorf("total_size argument is required")
	}

	if v, ok := d.GetOk("enable512e"); ok {
		volume.Enable512e = v.(bool)
	} else {
		return fmt.Errorf("enable512e argument is required")
	}

	if v, ok := d.GetOk("min_iops"); ok {
		volume.Qos.MinIOPS = int64(v.(int))
	}

	if v, ok := d.GetOk("max_iops"); ok {
		volume.Qos.MaxIOPS = int64(v.(int))
	}

	if v, ok := d.GetOk("burst_iops"); ok {
		volume.Qos.BurstIOPS = int64(v.(int))
	}

	resp, err := createVolume(client, volume)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.VolumeID))
	log.Printf("Created volume: %v %v", volume.Name, resp.VolumeID)

	return resourceSolidFireVolumeRead(d, meta)
}

func createVolume(client *sfapi.Client, request sftypes.CreateVolumeRequest) (*sftypes.CreateVolumeResult, error) {
	response, err := client.CreateVolume(request)
	if err != nil {
		log.Print("CreateVolume request failed")
		return &sftypes.CreateVolumeResult{}, err
	}

	return response, nil
}

func resourceSolidFireVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*sfapi.Client)

	volumes := sftypes.ListVolumesRequest{}

	id := d.Id()
	s := make([]int64, 1)
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	volumes.VolumeIDs = s

	res, err := listVolumes(client, volumes)
	if err != nil {
		return err
	}

	if len(res.Volumes) != 1 {
		return fmt.Errorf("Expected one Volume to be found. Response contained %v results", len(res.Volumes))
	}

	d.Set("name", res.Volumes[0].Name)

	return nil
}

func listVolumes(client *sfapi.Client, request sftypes.ListVolumesRequest) (*sftypes.ListVolumesResult, error) {
	response, err := client.ListVolumes(request)
	if err != nil {
		log.Print("ListVolumes request failed")
		return &sftypes.ListVolumesResult{}, err
	}

	return response, nil
}

func resourceSolidFireVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume access group %#v", d)
	client := meta.(*sfapi.Client)

	volume := sftypes.ModifyVolumeRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	volume.VolumeID = convID

	err := updateVolume(client, volume)
	if err != nil {
		return err
	}

	return nil
}

func updateVolume(client *sfapi.Client, request sftypes.ModifyVolumeRequest) error {
	_, err := client.ModifyVolume(request)
	if err != nil {
		log.Print("ModifyVolume request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume access group: %#v", d)
	client := meta.(*sfapi.Client)

	deleteVolumeReq := sftypes.DeleteVolumeRequest{}
	purgeVolumeReq := sftypes.PurgeDeletedVolumeRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	deleteVolumeReq.VolumeID = convID
	purgeVolumeReq.VolumeID = convID

	deleteErr := deleteVolume(client, deleteVolumeReq)
	if deleteErr != nil {
		return deleteErr
	}

	purgeErr := purgeDeletedVolume(client, purgeVolumeReq)
	if purgeErr != nil {
		return purgeErr
	}

	return nil
}

func deleteVolume(client *sfapi.Client, request sftypes.DeleteVolumeRequest) error {
	err := client.DeleteVolume(request)
	if err != nil {
		log.Print("DeleteVolume request failed")
		return err
	}

	return nil
}

func purgeDeletedVolume(client *sfapi.Client, request sftypes.PurgeDeletedVolumeRequest) error {
	err := client.PurgeDeletedVolume(request)
	if err != nil {
		log.Print("PurgeDeletedVolume request failed")
		return err
	}

	return nil
}

func resourceSolidFireVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %#v", d)
	client := meta.(*sfapi.Client)

	volumes := sftypes.ListVolumesRequest{}

	id := d.Id()
	s := make([]int64, 1)
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	volumes.VolumeIDs = s

	res, err := listVolumes(client, volumes)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.Volumes) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
