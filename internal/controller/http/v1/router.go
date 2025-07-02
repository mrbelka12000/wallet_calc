package v1

import "github.com/gorilla/mux"

func setUpRoutes(mx *mux.Router, c *Controller) {

	userRoutes := mx.PathPrefix("/user").Subrouter()
	userRoutes.HandleFunc("/register", c.Register).Methods("POST", "OPTIONS")
	userRoutes.HandleFunc("/login", c.Login).Methods("POST", "OPTIONS")
	userRoutes.HandleFunc("/profile", c.Profile).Methods("GET", "OPTIONS")
}
