package susepubliccloud

import (
	"fmt"
	"log"

	"github.com/SUSE/terraform-provider-susepubliccloud/pkg/info-service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceSUSEPublicCloudImageIds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSUSEPublicCloudImageIdsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"cloud": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "active",
				ValidateFunc: validateState,
			},
			"sort_ascending": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// validaetState is a SchemaValidateFunc which tests if the provided value is
// an accepted Image state
func validateState(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	err := images.ValidateState(v)
	if err != nil {
		es = append(es, err)
		return
	}

	return
}

func dataSourceSUSEPublicCloudImageIdsRead(d *schema.ResourceData, meta interface{}) error {
	params := images.SearchParams{
		Cloud:  d.Get("cloud").(string),
		Region: d.Get("region").(string),
	}

	if v, ok := d.GetOk("state"); ok {
		params.State = v.(string)
	}

	if d.Get("sort_ascending").(bool) {
		params.SortAscending = true
	} else {
		params.SortAscending = false
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		params.NameRegex = nameRegex.(string)
	}

	log.Printf("[DEBUG] Reading image IDs: %+v", params)
	images, err := images.GetImages(params)
	if err != nil {
		return err
	}

	imageIds := make([]string, 0)
	for _, image := range images {
		imageIds = append(imageIds, image.Id)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(fmt.Sprintf("%+v", params))))
	d.Set("ids", imageIds)

	return nil
}
