package main

import (
	"github.com/flavio/terraform-provider-susepubliccloud/susepubliccloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: susepubliccloud.Provider})
}
