package main

import (
  // Standard library packages
  "net/http"
  "github.com/gorilla/mux"
  // Connecting to the controller in same folder
  "github.com/HariniGB/openstack-api/controllers"
  // Third party packages
  // "github.com/julienschmidt/httprouter"
  // "gopkg.in/mgo.v2"
)
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
