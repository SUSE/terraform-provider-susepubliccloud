package susepubliccloud

import (
	"fmt"
	"hash/crc32"
	"log"

	images "github.com/SUSE/terraform-provider-susepubliccloud/pkg/info-service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceSUSEPublicCloudImageIDs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSUSEPublicCloudImageIDsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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

func dataSourceSUSEPublicCloudImageIDsRead(d *schema.ResourceData, meta interface{}) error {
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

	imageIDs := make([]string, 0)
	for _, image := range images {
		imageIDs = append(imageIDs, image.ID)
	}

	d.SetId(fmt.Sprintf("%d", stringTohashcode(fmt.Sprintf("%+v", params))))
	return d.Set("ids", imageIDs)
}

// String hashes a string to a unique hashcode.
//
// Copied from hashicorp/terraform-plugin-sdk/helper/hashcode/hashcode.go
// Because this is going to be dropped in future releases of the library
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func stringTohashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
