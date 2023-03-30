package routers

import (
	"memory-key-value-database-in-go/commands"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	// fmt.Println("routes 9")
	router.HandleFunc("/set", commands.Set()).Methods("POST")
	router.HandleFunc("/get", commands.Get()).Methods("POST")

}
