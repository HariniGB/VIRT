package controllers

import (
	"github.com/HariniGB/openstack-api/models"

	// To access the golang SDK for keystone authentication
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

//openstack compute lists of flavors,images and instances
func computeList() ([]models.FlavorsData, []models.ImagesData, []models.InstancesData){
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
	instanceDataList := []models.InstancesData{}

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

	// We have the option of filtering the server list. If we want the full
	// collection, leave it as an empty struct
	opts3 := servers.ListOpts{}

	// Retrieve a pager (i.e. a paginated collection)
	pager2 := servers.List(client, opts3)

	// Define an anonymous function to be executed on each page's iteration
	pager2.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		for _, s := range serverList {
			instanceDataList = append(instanceDataList,
				models.InstancesData{
					s.Name,
					s.ID,
					s.UserID,
					s.Image,
					s.Flavor,
					s.HostID,
					s.Status,
					s.SecurityGroups,
					s.Created,
					s.Updated,
				})
		}
		return true, err
	})

	return flavorDataList, imageDataList, instanceDataList
}
