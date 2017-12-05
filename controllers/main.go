package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/HariniGB/openstack-api/models"

	// To access the golang SDK for keystone authentication
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"

	// To maintain cookie for Login and Logout
	"github.com/gorilla/securecookie"
	//To get JSON in Rest API
	"encoding/json"
)

// cookie handling
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// variables initialization
var flavorDataList []models.FlavorsData
var imageDataList []models.ImagesData
var instanceDataList []models.InstancesData

// Login to openstack horizon and get authorized token for admin user from keystone
func keystoneAdmin() {
	//Pass in the values yourself
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://localhost:5000/v3/",
		Username:         "admin",
		Password:         "admin_user_secret",
		TenantName:       "admin",
		DomainName:       "Default",
	}
	fmt.Print("Request details:", opts)
	//Once you have the opts variable, you can pass it in and get back a ProviderClient struct:
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}
	fmt.Print("Authorization token details:", provider)
}

//Create  a session for the successful login user
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

// Clear the current session. Called in Logout handler
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

//openstack compute lists of flavors,images and instances
func computeList() {
  // Empty the struct array to avoid repetition
  flavorDataList = nil
	imageDataList = nil
	instanceDataList = nil

	//Pass in the values yourself
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://localhost:5000/v3/",
		Username:         "admin",
		Password:         "admin_user_secret",
		// TenantID: "71cbf6c1db784ce09aa75e0edc8464e9",
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
	fmt.Println("\nFlavors List:", flavorDataList)

	// We have the option of filtering the image list. If we want the full
	// collection, leave it as an empty struct
	opts2 := images.ListOpts{ChangesSince: "2014-01-01T01:02:03Z"}
	// Retrieve a pager (i.e. a paginated collection)
	pager1 := images.ListDetail(client, opts2)
	// Define an anonymous function to be executed on each page's iteration
	pager1.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		for _, i := range imageList {
			// "i" will be a images.Image
			imageDataList = append(imageDataList,
				models.ImagesData{
					i.Name,
					i.ID,
					i.MinDisk,
					i.Status,
					i.Progress,
					i.MinRAM,
					i.Metadata,
				})
		}
		return true, err
	})
	fmt.Println("Images List:", imageDataList)
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
	fmt.Println("Instances List:", instanceDataList)
}

// login handler
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		if name == "Demo User" && pass == "password" {
			keystoneAdmin()
			setSession(name, response)
			redirectTarget = "/dashboard"
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Login page")
}

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	Login(response, request)
	// fmt.Fprintf(response, indexPage)
}

// dashboard page
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Dashboard page")
}

func DashboardPageHandler(response http.ResponseWriter, request *http.Request) {
	Dashboard(response, request)
}

// openstack API page
func InstancesHandler(response http.ResponseWriter, request *http.Request) {
	computeList()
	input, _ := json.Marshal(instanceDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func FlavorsHandler(response http.ResponseWriter, request *http.Request) {
	computeList()
	input, _ := json.Marshal(flavorDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	computeList()
	input, _ := json.Marshal(imageDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func OpenStackPageHandler(response http.ResponseWriter, request *http.Request) {
	computeList()
	inputData := struct {
		Flavors   []models.FlavorsData
		Images    []models.ImagesData
		Instances []models.InstancesData
	}{
		flavorDataList,
		imageDataList,
		instanceDataList,
	}
	input, _ := json.Marshal(inputData)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
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
