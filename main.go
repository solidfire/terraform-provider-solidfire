package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/solidfire/terraform-provider-solidfire/solidfire"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: solidfire.Provider,
	})
}
