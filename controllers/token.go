package controllers

import (
  "fmt"
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/openstack"
)

func keystone_admin() {

  //Pass in the values yourself
  opts := gophercloud.AuthOptions{
      IdentityEndpoint: "http://localhost:5000/v3/",
      Username: "admin",
      Password: "admin_user_secret",
      TenantID: "71cbf6c1db784ce09aa75e0edc8464e9",
      TenantName: "admin",
      DomainName: "Default",
  }
  fmt.Print("Request details:",opts)

  //Once you have the opts variable, you can pass it in and get back a ProviderClient struct:
  provider, err := openstack.AuthenticatedClient(opts)
  if err != nil {
    panic(err)
  }
  fmt.Print("Authorization token details:",provider)
}


// OS_AUTH_URL="http://localhost:5000/v3/""
// OS_PROJECT_NAME="admin"
// OS_PROJECT_ID="71cbf6c1db784ce09aa75e0edc8464e9"
// OS_PASSWORD="admin_user_secret"
// OS_USER_DOMAIN_NAME="Default"
// OS_USERNAME="admin"
// OS_REGION_NAME="RegionOne"
// OS_INTERFACE=public
// OS_IDENTITY_API_VERSION=3
