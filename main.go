package main

import (
	"github.com/cprokopiak/terraform-provider-solidfire/solidfire"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: solidfire.Provider,
	})
}
