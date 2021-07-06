package images

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// Image describes an object returned by
// https://susepubliccloudinfo.suse.com/VERSION/FRAMEWORK/REGION/images.json
//
// {
//   "name": "suse-sles-15-sp1-v20190624-hvm-ssd-x86_64",
//   "state": "active",
//   "replacementname": "",
//   "replacementid": "",
//   "publishedon": "20190624",
//   "deprecatedon": "",
//   "region": "eu-central-1",
//   "id": "ami-0352b14942c00b04b",
//   "deletedon": ""
// },
type Image struct {
	Name            string `json:"name"`
	State           string `jsong:"state"`
	ReplacementName string `json:"replacementname,omitempty"`
	ReplacementID   string `json:"replacementid,omitempty"`
	PublishedOn     string `json:"publishedon"`
	DeprecatedOn    string `json:"deprecatedon,omitempty"`
	Region          string `json:"region"`
	ID              string `json:"id"`
	DeletedOn       string `json:"deletedon,omitempty"`
}

// Internally used to parse the response from
// SUSE public cloud info service API
type imagesReply struct {
	Images []Image `json:"images"`
}

// SearchParams is used to describe the search criteria to find one or more
// images
type SearchParams struct {
	APIEndpoint   string
	APIVersion    string
	Cloud         string
	NameRegex     string
	Region        string
	SortAscending bool
	State         string
}

// APIEndpoint is the endoint of the public instance of
// https://github.com/SUSE-Enceladus/public-cloud-info-service
const APIEndpoint = "https://susepubliccloudinfo.suse.com"

// APIVersion is the version of the enceladus API to be queried
const APIVersion = "v1"

// ValidImageStates holds the valid states of public cloud images as documented here:
// https://github.com/SUSE-Enceladus/public-cloud-info-service#server-design
var ValidImageStates = []string{
	"active",
	"inactive",
	"deprecated",
}

// GetImages returns a list of images that match the search criteria provided by
// the user.
func GetImages(params SearchParams) ([]Image, error) {
	images := make([]Image, 0)

	if err := ValidateState(params.State); err != nil {
		return images, err
	}

	if params.APIEndpoint == "" {
		params.APIEndpoint = APIEndpoint
	}

	if params.APIVersion == "" {
		params.APIVersion = APIVersion
	}

	baseURL, err := url.Parse(params.APIEndpoint)
	if err != nil {
		return images, err
	}

	urlPath := filepath.Join(
		params.APIVersion,
		params.Cloud,
		params.Region,
		"images",
		fmt.Sprintf("%s.json", params.State))
	relURL, err := baseURL.Parse(urlPath)
	if err != nil {
		return images, err
	}

	resp, err := http.Get(relURL.String())
	if err != nil {
		return images,
			fmt.Errorf("Error while accessing %v: %v", relURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return images,
			fmt.Errorf("Unexpected HTTP status %d while accessing %v",
				resp.StatusCode, relURL)
	}

	var reply imagesReply
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return images,
			fmt.Errorf("Error while decoding remote response from %s: %v",
				relURL, err)
	}

	if params.NameRegex != "" {
		r := regexp.MustCompile(params.NameRegex)
		for _, image := range reply.Images {
			if r.MatchString(image.Name) {
				images = append(images, image)
			}
		}
	} else {
		images = reply.Images[:]
	}

	sort.Slice(images, func(i, j int) bool {
		itime, _ := time.Parse("20060102", images[i].PublishedOn)
		jtime, _ := time.Parse("20060102", images[j].PublishedOn)
		if params.SortAscending {
			return itime.Unix() < jtime.Unix()
		}
		return itime.Unix() > jtime.Unix()
	})

	return images, nil
}

// ValidateState raises an error if the specified image state is not a valid one
func ValidateState(state string) error {
	for _, vs := range ValidImageStates {
		if state == vs {
			return nil
		}
	}

	return fmt.Errorf("Invalid image state: %s", state)
}
