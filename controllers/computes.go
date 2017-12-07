package controllers

import (
	"github.com/HariniGB/openstack-api/models"

	// To access the golang SDK for keystone authentication
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
	"os/exec"
	"log"
	"k8s.io/apimachinery/pkg/util/json"
)

//openstack compute lists of flavors,images and instances
func computeList() ([]models.FlavorsData, []models.ImagesData, []models.StackResponse){
	// Empty the struct array to avoid repetition
	flavorDataList := []models.FlavorsData{}
	imageDataList := []models.ImagesData{
		{
			Name: "cirros",
			Type: "Operating System",

		},
		{
			Name: "mysql",
			Type: "database",
		},
	}

	//Pass in the values yourself
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://controller:5000/v3/",
		Username:         "admin",
		Password:         "admin_user_secret",
		TenantName: "admin",
		DomainName: "Default",
	}

	//Once you have the opts variable, you can pass it in and get back a ProviderClient struct:
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		panic(err)
	}

	// We have the option of filtering the flavor list. If we want the full
	// collection, leave it as an empty struct
	opts1 := flavors.ListOpts{ChangesSince: "2014-01-01T01:02:03Z"}

	// Retrieve a pager (i.e. a paginated collection)
	pager := flavors.ListDetail(client, opts1)
	// Define an anonymous function to be executed on each page's iteration
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		for _, f := range flavorList {
			// "f" will be a flavors.Flavor
			flavorDataList = append(flavorDataList,
				models.FlavorsData{
					f.Name,
					f.ID,
					f.RAM,
					f.VCPUs,
					f.Disk,
					f.RxTxFactor,
				})
		}
		return true, err
	})

	instanceDataList := GetStacks()
	return flavorDataList, imageDataList, instanceDataList
}

func GetQuotas() *models.QuotaResponse {
	arr := []string {"limits", "show", "--absolute",  "-f", "json"}

	cmd := exec.Command("openstack",  arr...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(b)
		return nil
	}

	out := []models.QuotaEntry{}
	err = json.Unmarshal(b, &out)
	if err != nil {
		log.Print(err)
		return nil
	}

	resp := &models.QuotaResponse{}
	resp.FromQuotas(out)
	return resp
}