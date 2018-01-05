package main

import (
	"bitbucket.org/solidfire/terraform-provider-solidfire/solidfire"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: solidfire.Provider,
	})
}
