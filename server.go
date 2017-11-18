package main

import (
  // Standard library packages
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/securecookie"
  "html/template"
  "fmt"
  // Connecting to the controller in same folder
  // "github.com/HariniGB/openstack-api/controllers"
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/openstack"
  // Third party packages
  // "github.com/julienschmidt/httprouter"
  // "gopkg.in/mgo.v2"
)

// cookie handling
var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32))

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

func loginHandler(response http.ResponseWriter, request *http.Request) {
  name := request.FormValue("name")
  pass := request.FormValue("password")
  redirectTarget := "/"
  if name != "" && pass != "" {
    if name == "Demo User" && pass == "password" {
    keystone_admin()
    setSession(name, response)
    redirectTarget = "/dashboard"
    }
  }
  http.Redirect(response, request, redirectTarget, 302)
}

// logout handler
func logoutHandler(response http.ResponseWriter, request *http.Request) {
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

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
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

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
  userName := getUserName(request)
  if userName != "" {
   Dashboard(response, request)
  } else {
    http.Redirect(response, request, "/", 302)
  }
}

var router = mux.NewRouter()

func main() {
  router.HandleFunc("/", indexPageHandler)
  router.HandleFunc("/dashboard", internalPageHandler)

  router.HandleFunc("/login", loginHandler).Methods("POST")
  router.HandleFunc("/logout", logoutHandler).Methods("POST")

  http.Handle("/", router)
  http.ListenAndServe(":8000", nil)
}

// getSession creates a new mongo session and panics if connection error occurs
// func getSession() *mgo.Session {
//   // Connect to our local mongo
//   s, err := mgo.Dial("mongodb://localhost:27017")

//   // Check if connection error, is mongo running?
//   if err != nil {
//     panic(err)
//   }

//   // Deliver session
//   return s
// }
