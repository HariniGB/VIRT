package controllers
import (
  // "fmt"
  "github.com/gorilla/securecookie"
  "net/http"
  "html/template"
)

// cookie handling
var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32))

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
    if name == "Demo User" && pass != "password" {
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
