package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Lights",
		"GET",
		"/lights",
		GetLights,
	},
	Route{
		"Light",
		"GET",
		"/lights/{id}",
		GetLight,
	},
	Route{
		"SetState",
		"GET",
		"/lights/{id}/{command}",
		SetState,
	},
	Route{
		"Dimmer",
		"GET",
		"/lights/{id}/level/{value}",
		SetDimmer,
	},
}

/* router := mux.NewRouter()

router.HandleFunc("/lights", GetLights).Methods("GET")
router.HandleFunc("/lights/{id}", GetLight).Methods("GET")
router.HandleFunc("/lights/{id}/{command}", SetState).Methods("GET")
router.HandleFunc("/lights/{id}/level/{value}", SetDimmer).Methods("GET")

// router.HandleFunc("lights/{id}/on", SetState).Methods("GET")
*/

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}
