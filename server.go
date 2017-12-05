package main

import (
  // Standard library packages
  "net/http"
  // Connecting to the controller in same folder
  "github.com/HariniGB/openstack-api/controllers"
  // Third party packages
)

func main() {
  // Get a UserController instance
  http.HandleFunc("/", controllers.IndexPageHandler)
  http.HandleFunc("/dashboard", controllers.DashboardPageHandler)
  http.HandleFunc("/api/v1/openstack/dashboard", controllers.OpenStackPageHandler)
  http.HandleFunc("/login", controllers.LoginHandler)
  http.HandleFunc("/logout", controllers.LogoutHandler)
  // To add all static contents like CSS, JS, IMAGES for the HTML page at any given route (/*filepath)
  http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
  http.ListenAndServe("localhost:3000", nil)
}
