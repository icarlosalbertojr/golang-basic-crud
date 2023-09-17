package main

import (
	"my-db-module/server"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/users", server.CreateNewUser).Methods(http.MethodPost)
	router.HandleFunc("/v1/users", server.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/v1/users/{id}", server.GetUsersById).Methods(http.MethodGet)
	router.HandleFunc("/v1/users/{id}", server.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/v1/users/{id}", server.DeleteUserById).Methods(http.MethodDelete)

	http.ListenAndServe(":5000", router)
}
