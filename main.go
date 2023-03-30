package main

import (
	"log"
	"memory-key-value-database-in-go/configs"
	"memory-key-value-database-in-go/routers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//run database
	configs.ConnectDB()
	// auth.CreateNewReferenceToken("abc@gmail")
	//routes
	// router.UserRoute(router)
	routers.UserRoute(router)
	log.Fatal(http.ListenAndServe(":8080", router))
	// router := mux.NewRouter()
	// // fmt.Println("main 13")
	// router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	// 	rw.Header().Set("Content-Type", "application/json")
	// 	fmt.Println("main 17")
	// 	json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB"})
	// }).Methods("GET")
	// // fmt.Println("main 20")
	// log.Fatal(http.ListenAndServe(":6000", router))
	// // fmt.Println("main 22")
}
