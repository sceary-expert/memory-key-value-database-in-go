package router

import (
	"github.com/gorilla/mux"
	"github.com/sceary-expert/memory-key-value-database-in-go/commands"
)

func UserRoute(router *mux.Router) {

	router.HandleFunc("/set", commands.Set()).Methods("POST")
	router.HandleFunc("/get", commands.Get()).Methods("POST")

}
