package controllers

import (
  "net/http"
  "html/template"
  "fmt"
  "github.com/julienschmidt/httprouter"

  // To access the golang SDK for keystone authentication
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/openstack"

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
)


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
  if err != nil {
    panic(err)
  }
   tmpl.Execute(w, "Dashboard page")
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

