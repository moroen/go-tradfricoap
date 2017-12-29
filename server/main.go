package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	coap "github.com/moroen/go-tradfricoap"
)

type Todo struct {
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

const port = "8085"

func main() {

	err := coap.LoadConfig()
	if err != nil {
		panic("\nNo config found!")
	}

	router := NewRouter()
	log.Printf("Starting server - listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
