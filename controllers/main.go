package controllers

import (
  "net/http"
  "html/template"
  "fmt"
  "github.com/julienschmidt/httprouter"
  "strconv"

  // To access the golang SDK for keystone authentication
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/openstack"
  "github.com/rackspace/gophercloud/pagination"
  "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
  "github.com/rackspace/gophercloud/openstack/compute/v2/images"
  "github.com/rackspace/gophercloud/openstack/compute/v2/servers"

  // To maintain cookie for Login and Logout
  "github.com/gorilla/securecookie"

  // To run a mongodb session
  "gopkg.in/mgo.v2"
)

// cookie handling
var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32))

type (
  // UserController represents the controller for operating on the User resource
  UserController struct {
    session *mgo.Session
  }

  FlavorsData struct {
    Name  string
    ID    string
    RAM   string
    VCPUs string
    Disk  string
    RXTX  string
  }
)
//Public variables
var flavorDataList []FlavorsData

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
  return &UserController{s}
}

func keystoneAdmin() {
  //Pass in the values yourself
  opts := gophercloud.AuthOptions{
      IdentityEndpoint: "http://localhost:5000/v3/",
      Username: "admin",
      Password: "admin_user_secret",
      // TenantID: "71cbf6c1db784ce09aa75e0edc8464e9",
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

func listOpenstackDetails() {
    //Pass in the values yourself
    opts := gophercloud.AuthOptions{
      IdentityEndpoint: "http://localhost:5000/v3/",
      Username: "admin",
      Password: "admin_user_secret",
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
			  FlavorsData{
				  f.Name,
				  string(f.ID),
				  strconv.Itoa(f.RAM),
				  strconv.Itoa(f.VCPUs),
				  strconv.Itoa(f.Disk),
				  strconv.FormatFloat(f.RxTxFactor, 'f', 2, 64),
			  },
		  )
	  }
      return true, err
    })
	fmt.Print("FlavorsList:",flavorDataList)
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
        fmt.Print("\n Image Name:",i.Name)
        fmt.Print("\n ID:",i.ID)
        fmt.Print("\n Min.Disk:",i.MinDisk)
        fmt.Print("\n Status:",i.Status)
        fmt.Print("\n Progress:", i.Progress)
        fmt.Print("\n Min.RAM:",i.MinRAM)
        fmt.Print("\n Metadata:",i.Metadata)
        fmt.Print("\n=============================")
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
        // "s" will be a servers.Server
        fmt.Print("\n Instance Name:",s.Name)
        fmt.Print("\n ID:",s.ID)
        fmt.Print("\n Owner ID:",s.UserID)
        fmt.Print("\n Image:",s.Image)
        fmt.Print("\n Flavor:",s.Flavor)
        fmt.Print("\n Host:", s.HostID)
        fmt.Print("\n Status:",s.Status)
        fmt.Print("\n Security groups:",s.SecurityGroups)
        fmt.Print("\n Created At:",s.Created)
        fmt.Print("\n Updated At:", s.Updated)
        fmt.Print("\n=============================")
      }
      return true, err
    })
}


func getUserName(request *http.Request) (userName string) {
  if cookie, err := request.Cookie("session"); err == nil {
    cookieValue := make(map[string]string)
    if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
      userName = cookieValue["name"]
    }
  }
  return userName
}

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

func clearSession(response http.ResponseWriter) {
  cookie := &http.Cookie{
    Name:   "session",
    Value:  "",
    Path:   "/",
    MaxAge: -1,
  }
  http.SetCookie(response, cookie)
}

// login handler

func  (uc UserController) LoginHandler(response http.ResponseWriter, request *http.Request,  p httprouter.Params) {
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
func  (uc UserController) LogoutHandler(response http.ResponseWriter, request *http.Request,  p httprouter.Params) {
  clearSession(response)
  http.Redirect(response, request, "/", 302)
}

// index page
func Login(w http.ResponseWriter, r *http.Request){
  tmpl, err := template.ParseFiles("templates/login.html")
  if err != nil {
    panic(err)
  }
   tmpl.Execute(w, "Login page")
}

func  (uc UserController) IndexPageHandler(response http.ResponseWriter, request *http.Request,  p httprouter.Params) {
  Login(response, request)
  // fmt.Fprintf(response, indexPage)
}

// internal page
func Dashboard(w http.ResponseWriter, r *http.Request){
  tmpl, err := template.ParseFiles("templates/dashboard.html")
  listOpenstackDetails()
  if err != nil {
    panic(err)
  }
   tmpl.Execute(w, &flavorDataList)
}

func  (uc UserController) InternalPageHandler(response http.ResponseWriter, request *http.Request,  p httprouter.Params) {
  userName := getUserName(request)
   fmt.Print("User-name: ",userName)
  if userName != "" {
    Dashboard(response, request)
  } else {
    http.Redirect(response, request, "/", 302)
  }
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

