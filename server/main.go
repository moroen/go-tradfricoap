package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	coap "github.com/moroen/tradfricoap"
)

func main() {

	err := coap.LoadConfig()
	if err != nil {
		panic("\nNo config found!")
	}

	router := mux.NewRouter()

	router.HandleFunc("/lights", GetLights).Methods("GET")
	log.Fatal(http.ListenAndServe(":8085", router))
	/*
		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)
		api.Use(&rest.JsonpMiddleware{
			CallbackNameKey: "cb",
		})
		api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		}))
		log.Fatal(http.ListenAndServe(":8085", api.MakeHandler()))
	*/
}

func GetLights(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting lights")
	lights := coap.GetDevices()
	json.NewEncoder(w).Encode(lights)
}
