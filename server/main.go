package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	coap "github.com/moroen/go-tradfricoap"
)

type returnMessageSimple struct {
	Action string
	Status string
	Result string
}

const port = "8085"

func main() {

	err := coap.LoadConfig()
	if err != nil {
		panic("\nNo config found!")
	}

	router := mux.NewRouter()

	router.HandleFunc("/lights", GetLights).Methods("GET")
	router.HandleFunc("/lights/{id}", GetLight).Methods("GET")
	router.HandleFunc("/lights/{id}/{command}", SetState).Methods("GET")
	router.HandleFunc("/lights/{id}/level/{value}", SetDimmer).Methods("GET")

	// router.HandleFunc("lights/{id}/on", SetState).Methods("GET")
	log.Printf("Starting server - listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
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
	lights, err := coap.GetDevices()
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(lights)
}

func GetLight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println("Getting light: ", params["id"])
}

func SetState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	state := 0

	log.Println("SetState: ", params["id"], params["command"])
	if params["command"] == "on" {
		state = 1
	}

	coap.SetState(params["id"], state)
}

func SetDimmer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	log.Println("SetDimmer: ", params["id"], params["value"])

	if value, err := strconv.Atoi(params["value"]); err == nil {
		coap.SetLevel(params["id"], value)
	} else {
		log.Println("Failed to set level")
		errMsg := returnMessageSimple{Action: "setLevel", Status: "error", Result: err.Error()}
		json.NewEncoder(w).Encode(errMsg)
	}
}
