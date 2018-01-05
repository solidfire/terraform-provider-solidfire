package solidfire

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_USERNAME", nil),
				Description: "The user name for SolidFire API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_PASSWORD", nil),
				Description: "The user password for SolidFire API operations.",
			},
			"solidfire_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_SERVER", nil),
				Description: "The SolidFire server name for SolidFire API operations.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDFIRE_API_VERSION", nil),
				Description: "The SolidFire server API version.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"solidfire_volume_access_group": resourceSolidFireVolumeAccessGroup(),
			"solidfire_initiator":           resourceSolidFireInitiator(),
			"solidfire_volume":              resourceSolidFireVolume(),
			"solidfire_account":             resourceSolidFireAccount(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		User:            d.Get("username").(string),
		Password:        d.Get("password").(string),
		SolidFireServer: d.Get("solidfire_server").(string),
		APIVersion:      d.Get("api_version").(string),
	}

	return config.Client()
}
