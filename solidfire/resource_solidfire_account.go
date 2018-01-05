package solidfire

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/solidfire/solidfire-sdk-golang/sfapi"
	"github.com/solidfire/solidfire-sdk-golang/sftypes"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element"
	"github.com/solidfire/terraform-provider-solidfire/solidfire/element/jsonrpc"
)

type CreateAccountRequest struct {
	Username        string      `structs:"username"`
	InitiatorSecret string      `structs:"initiatorSecret,omitempty"`
	TargetSecret    string      `structs:"targetSecret,omitempty"`
	Attributes      interface{} `structs:"attributes,omitempty"`
}

type CreateAccountResult struct {
	Account element.Account `json:"account"`
}

type ModifyAccountRequest struct {
	AccountID       int         `structs:"accountID"`
	InitiatorSecret string      `structs:"initiatorSecret,omitempty"`
	TargetSecret    string      `structs:"targetSecret,omitempty"`
	Attributes      interface{} `structs:"attributes,omitempty"`
	Username        string      `structs:"username,omitempty"`
}

type RemoveAccountRequest struct {
	AccountID int `structs:"accountID"`
}

func resourceSolidFireAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolidFireAccountCreate,
		Read:   resourceSolidFireAccountRead,
		Update: resourceSolidFireAccountUpdate,
		Delete: resourceSolidFireAccountDelete,
		Exists: resourceSolidFireAccountExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initiator_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceSolidFireAccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating account: %#v", d)
	client := meta.(*sfapi.Client)

	acct := sftypes.AddAccountRequest{}

	if v, ok := d.GetOk("username"); ok {
		acct.Username = v.(string)
	} else {
		return fmt.Errorf("username argument is required")
	}

	if v, ok := d.GetOk("initiator_secret"); ok {
		acct.InitiatorSecret = sftypes.CHAPSecret{
			Secret: v.(string),
		}
	}

	if v, ok := d.GetOk("target_secret"); ok {
		acct.TargetSecret = sftypes.CHAPSecret{
			Secret: v.(string),
		}
	}

	resp, err := createAccount(client, acct)
	if err != nil {
		log.Print("Error creating account")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.Account.AccountID))

	log.Printf("Created account: %v %v", acct.Username, resp.Account.AccountID)

	return resourceSolidFireAccountRead(d, meta)
}

func createAccount(client *sfapi.Client, request sftypes.AddAccountRequest) (*sftypes.AddAccountResult, error) {
	log.Printf("Request: %v", request)

	response, err := client.AddAccount(request)
	if err != nil {
		log.Print("CreateAccount request failed")
		return &sftypes.AddAccountResult{}, err
	}

	return response, nil
}

func resourceSolidFireAccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading account: %#v", d)
	client := meta.(*sfapi.Client)

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	acct := sftypes.GetAccountByIDRequest{
		AccountID: convID,
	}

	res, err := client.GetAccountByID(acct)
	if err != nil {
		log.Print("GetAccountByID failed")
		return err
	}

	if _, ok := d.GetOk("username"); ok {
		d.Set("username", res.Account.Username)
	}

	if _, ok := d.GetOk("initiator_secret"); ok {
		d.Set("initiator_secret", res.Account.InitiatorSecret)
	}

	if _, ok := d.GetOk("target_secret"); ok {
		d.Set("target_secret", res.Account.TargetSecret)
	}

	return nil
}

func resourceSolidFireAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating account %#v", d)
	client := meta.(*element.Client)

	acct := ModifyAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	acct.AccountID = convID

	if v, ok := d.GetOk("username"); ok {
		acct.Username = v.(string)
	}

	if v, ok := d.GetOk("initiator_secret"); ok {
		acct.InitiatorSecret = v.(string)
	}

	if v, ok := d.GetOk("target_secret"); ok {
		acct.TargetSecret = v.(string)
	}

	err := modifyAccount(client, acct)
	if err != nil {
		return err
	}

	return nil
}

func modifyAccount(client *element.Client, request ModifyAccountRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyAccount", params)
	if err != nil {
		log.Print("ModifyAccount request failed")
		return err
	}

	return nil
}

func resourceSolidFireAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting account: %#v", d)
	client := meta.(*element.Client)

	acct := RemoveAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	acct.AccountID = convID

	err := removeAccount(client, acct)
	if err != nil {
		return err
	}

	return nil
}

func removeAccount(client *element.Client, request RemoveAccountRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("RemoveAccount", params)
	if err != nil {
		log.Print("DeleteAccount request failed")
		return err
	}

	return nil
}

func resourceSolidFireAccountExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of account: %#v", d)
	client := meta.(*sfapi.Client)

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	acct := sftypes.GetAccountByIDRequest{
		AccountID: convID,
	}

	_, err := client.GetAccountByID(acct)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknownAccount" {
				d.SetId("")
				return false, nil
			}
		}
		log.Print("AccountExists failed")
		return false, err
	}

	return true, nil
}
