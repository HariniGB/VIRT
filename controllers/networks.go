package controllers

import (
	"fmt"

	// To access the golang SDK for keystone authentication
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"log"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)


func CreateNetwork(username, subnet string, provider *gophercloud.ProviderClient) error {

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Type: "network",
		Name:   "neutron",
		Region: "RegionOne",
	})

	if err != nil {
		log.Println(err)
		return fmt.Errorf("unable to create network client")
	}

	tr := true

	opts := networks.CreateOpts{Name: username, AdminStateUp: &tr}
	network, err := networks.Create(client, opts).Extract()

	if err != nil {
		log.Println(err)
		return fmt.Errorf("unable to create network")
	}

	subOpts := subnets.CreateOpts{
		NetworkID:  network.ID,
		CIDR:       subnet,
		IPVersion:  subnets.IPv4,
		Name:       fmt.Sprintf("%s-subnet", username),
	}

	// Execute the operation and get back a subnets.Subnet struct
	sub, err := subnets.Create(client, subOpts).Extract()
	if err != nil {
		return fmt.Errorf("unable to create subnet")
	}

	up := true
	routerOpts := routers.CreateOpts{
		Name: username,
		AdminStateUp: &up,
		TenantID: username,
		GatewayInfo: &routers.GatewayInfo{NetworkID: "2cbde268-077c-4428-ab9e-d605d611212d"},
	}

	router, err := routers.Create(client, routerOpts).Extract()
	if err != nil {
		return fmt.Errorf("unable to create router")
	}

	intOpts := routers.InterfaceOpts{
		SubnetID: sub.ID,
	}

	routers.AddInterface(client, router.ID, intOpts)
	return nil
}