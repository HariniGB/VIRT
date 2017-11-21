package main

import (
  // Standard library packages
  "net/http"
  // Connecting to the controller in same folder
  "github.com/HariniGB/openstack-api/controllers"
  // Third party packages
  "github.com/julienschmidt/httprouter"
  "gopkg.in/mgo.v2"
)

func main() {
  // Instantiate a new router
  r := httprouter.New()

  // Get a UserController instance
  uc := controllers.NewUserController(getSession())

  r.GET("/", uc.IndexPageHandler)
  r.GET("/dashboard", uc.DashboardPageHandler)
  r.GET("/api/v1/openstack/dashboard", uc.OpenStackPageHandler)
  r.POST("/login", uc.LoginHandler)
  r.POST("/logout", uc.LogoutHandler)
  // To add all static contents like CSS, JS, IMAGES for the HTML page at any given route (/*filepath)
  r.ServeFiles("/static/*filepath", http.Dir("static"))
  http.ListenAndServe("localhost:3000", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
  // Connect to our local mongo
  s, err := mgo.Dial("mongodb://localhost:27017")
  // Check if connection error, is mongo running?
  if err != nil {
    panic(err)
  }

  // Deliver session
  return s
}
