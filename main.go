package main

import (
	"github.com/SUSE/terraform-provider-susepubliccloud/susepubliccloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: susepubliccloud.Provider})
}
