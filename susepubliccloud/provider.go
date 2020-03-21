package susepubliccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"susepubliccloud_image_ids": dataSourceSUSEPublicCloudImageIds(),
		},

		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return nil, nil
}
