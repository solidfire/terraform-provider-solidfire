package solidfire

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/cprokopiak/terraform-provider-solidfire/solidfire/element"
	"github.com/cprokopiak/terraform-provider-solidfire/solidfire/element/jsonrpc"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

type StorageDevice struct {
	Device string
	IQN    string
}

type CreateInitiatorsRequest struct {
	Initiators []element.Initiator `structs:"initiators"`
}

type CreateInitiatorsResult struct {
	Initiators []element.InitiatorResponse `json:"initiators"`
}

type DeleteInitiatorsRequest struct {
	Initiators []int `structs:"initiators"`
}

type ModifyInitiatorsRequest struct {
	Initiators []element.Initiator `structs:"initiators"`
}

func resourceSolidFireInitiator() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolidFireInitiatorCreate,
		Read:   resourceSolidFireInitiatorRead,
		Update: resourceSolidFireInitiatorUpdate,
		Delete: resourceSolidFireInitiatorDelete,
		Exists: resourceSolidFireInitiatorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"volume_access_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iqns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSolidFireInitiatorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating initiator: %#v", d)
	client := meta.(*element.Client)

	initiators := CreateInitiatorsRequest{}
	newInitiator := make([]element.Initiator, 1)
	var iqns []string

	if v, ok := d.GetOk("name"); ok {
		newInitiator[0].Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("alias"); ok {
		newInitiator[0].Alias = v.(string)
	}

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		newInitiator[0].VolumeAccessGroupID = v.(int)
	}

	if v, ok := d.GetOk("iqns"); ok {

		if a, ok := v.([]interface{}); ok {
			for i := range a {
				iqns = append(iqns, a[i].(string))
			}
		}
	}

	initiators.Initiators = newInitiator

	resp, err := createInitiators(client, initiators)
	if err != nil {
		log.Print("Error creating initiator")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.Initiators[0].ID))
	log.Printf("Created initiator: %v %v", newInitiator[0].Name, resp.Initiators[0].ID)

	return resourceSolidFireInitiatorRead(d, meta)
}

func createInitiators(client *element.Client, request CreateInitiatorsRequest) (CreateInitiatorsResult, error) {
	params := structs.Map(request)

	log.Printf("Parameters: %v", params)

	response, err := client.CallAPIMethod("CreateInitiators", params)
	if err != nil {
		log.Print("CreateInitiators request failed")
		return CreateInitiatorsResult{}, err
	}

	var result CreateInitiatorsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall resposne from CreateInitiators")
		return CreateInitiatorsResult{}, err
	}
	log.Printf("Initiator Result: %v", result)
	return result, nil
}

func resourceSolidFireInitiatorRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading initiator: %#v", d)
	client := meta.(*element.Client)

	initiators := element.ListInitiatorRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	res, err := listInitiators(client, initiators)
	if err != nil {
		return err
	}

	if len(res.Initiators) != 1 {
		return fmt.Errorf("Expected one Initiator to be found. Response contained %v results", len(res.Initiators))
	}

	d.Set("name", res.Initiators[0].Name)
	d.Set("alias", res.Initiators[0].Alias)
	d.Set("attributes", res.Initiators[0].Attributes)

	if len(res.Initiators[0].VolumeAccessGroups) == 1 {
		d.Set("volume_access_group_id", res.Initiators[0].VolumeAccessGroups[0])
	}

	return nil
}

func listInitiators(client *element.Client, request element.ListInitiatorRequest) (element.ListInitiatorResult, error) {
	params := structs.Map(request)

	response, err := client.CallAPIMethod("ListInitiators", params)
	if err != nil {
		log.Print("ListInitiators request failed")
		return element.ListInitiatorResult{}, err
	}

	var result element.ListInitiatorResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListInitiators")
		return element.ListInitiatorResult{}, err
	}

	return result, nil
}

func resourceSolidFireInitiatorUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating initiator: %#v", d)
	client := meta.(*element.Client)

	initiators := ModifyInitiatorsRequest{}
	initiator := make([]element.Initiator, 1)

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	initiator[0].InitiatorID = convID

	if v, ok := d.GetOk("alias"); ok {
		initiator[0].Alias = v.(string)
	}

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		initiator[0].VolumeAccessGroupID = v.(int)
	}

	initiators.Initiators = initiator

	err := modifyInitiators(client, initiators)
	if err != nil {
		return err
	}

	return nil
}

func modifyInitiators(client *element.Client, request ModifyInitiatorsRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyInitiators", params)
	if err != nil {
		log.Print("ModifyInitiators request failed")
		return err
	}

	return nil
}

func resourceSolidFireInitiatorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting initiator: %#v", d)
	client := meta.(*element.Client)

	initiators := DeleteInitiatorsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	err := deleteInitiator(client, initiators)
	if err != nil {
		return err
	}

	return nil
}

func deleteInitiator(client *element.Client, request DeleteInitiatorsRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("DeleteInitiators", params)
	if err != nil {
		log.Print("DeleteInitiator request failed")
		return err
	}

	return nil
}

func resourceSolidFireInitiatorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of initiator: %#v", d)
	client := meta.(*element.Client)

	initiators := element.ListInitiatorRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	res, err := listInitiators(client, initiators)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.Initiators) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
