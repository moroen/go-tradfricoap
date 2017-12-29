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

func Index(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome!")
}

func GetLights(w http.ResponseWriter, r *http.Request) {
	lights, _, err := coap.GetDevices()
	if err != nil {
		panic(err.Error())
	}
	answer := returnMessageDevices{Action: "getLights", Status: "Ok", Result: lights}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

func GetLight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println("Getting light: ", params["id"])
	if device, err := coap.GetLight(params["id"]); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(device)
	}
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errMsg)
	}
}
