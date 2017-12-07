package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/HariniGB/openstack-api/models"

	// To access the golang SDK for keystone authentication
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"

	//To get JSON in Rest API
	"encoding/json"
	"log"
	"github.com/HariniGB/openstack-api/heat"
)

// Login to openstack horizon and get authorized token for admin user from keystone
func keystoneAdmin() (*gophercloud.ProviderClient, error) {
	return KeystoneLogin("admin", "admin_user_secret", "")
}

func KeystoneLogin(username, password, token string) (*gophercloud.ProviderClient, error){
	//Pass in the values yourself
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://controller:5000/v3/",
	}
	if password != "" {
		opts.Username = username
		opts.Password = password
		opts.DomainName = "Default"
	} else if token != "" {
		opts.TokenID = token
	} else {
		return nil, fmt.Errorf("username/password or token is mandatory")
	}

	//Once you have the opts variable, you can pass it in and get back a ProviderClient struct:
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Println("Auth failure: ", err)
		return nil, fmt.Errorf("Authentication failure")
	}

	return provider, nil
}

// login handler
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
			client, err := KeystoneLogin(name, pass, "")
			if err == nil {
				setSession(name, client.TokenID, response)
				redirectTarget = "/dashboard"
			} else {
				log.Println(err)
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

// flavors and applications page
func FlavorsAppln(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/flavorsApplications.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Flavors and Applications page")
}

func FlavorsApplnPageHandler(response http.ResponseWriter, request *http.Request) {
	FlavorsAppln(response, request)
}

func CreateInstancePageHandler(response http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("templates/form.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(response, "Create new Instance")
}


// openstack API page
func InstancesHandler(response http.ResponseWriter, request *http.Request) {
	_, _, instanceDataList := computeList()
	input, _ := json.Marshal(instanceDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func FlavorsHandler(response http.ResponseWriter, request *http.Request) {
	flavorDataList, _, _ := computeList()
	input, _ := json.Marshal(flavorDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	_, imageDataList, _ := computeList()
	input, _ := json.Marshal(imageDataList)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)
	fmt.Fprintf(response, "%s", input)
}

func InstanceHandler(response http.ResponseWriter, request *http.Request){
	u := models.ProvisionRequest{}
	if request.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(request.Body).Decode(&u)
	} else {
		u.Name = request.FormValue("name")
		u.Flavor = request.FormValue("flavor")
		u.Type = request.FormValue("type")
	}

	params := map[string]string{
		"name": "admin",
		"net_name": "admin",
		"subnet_name": fmt.Sprintf("%s-subnet", "admin"),
		"flavor_name": u.Flavor,
	}

	var err error
	switch u.Type {
	case "cirros":
		err = DeployTemplate(u.Name, heat.Cirros, params)
	case "tomcat":
		err = DeployTemplate(u.Name, heat.Tomcat, params)
	default:
		err = fmt.Errorf("unknown application type")
	}

	if err != nil {
		response.WriteHeader(500)
		fmt.Fprintf(response, "%v", err)
		return
	}

	http.Redirect(response, request, "/dashboard", 302)
}

func OpenStackPageHandler(response http.ResponseWriter, request *http.Request) {
	flavorDataList, imageDataList, instanceDataList := computeList()
	inputData := struct {
		Flavors   []models.FlavorsData
		Images    []models.ImagesData
		Instances []models.StackResponse
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

func CreateInstance(response http.ResponseWriter, request *http.Request) {
	/*user, pass := getSession(request)

	if user == "" || pass == "" {
		response.WriteHeader(401)
		fmt.Fprintf(response, "Unauthenticated")
		return
	}

	_, err := KeystoneLogin(user, "", pass)
	if err != nil {
		response.WriteHeader(403)
		fmt.Fprintf(response, "Unauthenticated")
		return
	}*/

	u := models.ProvisionRequest{}
	if request.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(request.Body).Decode(&u)
	} else {
		u.Name = request.FormValue("name")
		u.Flavor = request.FormValue("flavor")
		u.Type = request.FormValue("type")
	}

	params := map[string]string{
		"name": "admin",
		"net_name": "admin",
		"subnet_name": fmt.Sprintf("%s-subnet", "admin"),
		"flavor_name": u.Flavor,
	}

	var err error
	switch u.Type {
	case "cirros":
		err = DeployTemplate(u.Name, heat.Cirros, params)
	case "tomcat":
		err = DeployTemplate(u.Name, heat.Tomcat, params)
	default:
		err = fmt.Errorf("unknown application type")
	}

	if err != nil {
		response.WriteHeader(500)
		fmt.Fprintf(response, "%v", err)
		return
	}

	response.WriteHeader(200)
}

func HandleQuotas(response http.ResponseWriter, request *http.Request) {
	/*user, pass := getSession(request)

	if user == "" || pass == "" {
		response.WriteHeader(401)
		fmt.Fprintf(response, "Unauthenticated")
		return
	}

	_, err := KeystoneLogin(user, "", pass)
	if err != nil {
		response.WriteHeader(403)
		fmt.Fprintf(response, "Unauthenticated")
		return
	}*/

	quotas := GetQuotas()

	if quotas == nil {
		response.WriteHeader(500)
		fmt.Fprintf(response, "Unable to retrive quotas")
		return
	}

	bytes, err := json.Marshal(quotas)
	if err != nil {
		response.WriteHeader(500)
		fmt.Fprintf(response, "Unable to marshall quotas")
		return
	}

	response.WriteHeader(200)
	fmt.Fprintf(response, string(bytes))
}