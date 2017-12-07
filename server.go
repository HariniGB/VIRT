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
	http.HandleFunc("/details", controllers.FlavorsApplnPageHandler)
	http.HandleFunc("/instances", controllers.CreateInstancePageHandler)
	http.HandleFunc("/api/v1/openstack/dashboard", controllers.OpenStackPageHandler)
	http.HandleFunc("/api/v1/openstack/flavors", controllers.FlavorsHandler)
	http.HandleFunc("/api/v1/openstack/images", controllers.ImagesHandler)
	http.HandleFunc("/api/v1/openstack/instances", controllers.InstancesHandler)
	http.HandleFunc("/api/v1/openstack/instance", controllers.CreateInstance	)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/instance", controllers.InstanceHandler)
	// To add all static contents like CSS, JS, IMAGES for the HTML page at any given route (/*filepath)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":3000", nil)
}
