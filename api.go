package main

import (
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/handlers"
	"net/http"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/station/{id:[\\d]+}", handlers.GetStation)
	rtr.HandleFunc("/station/search", handlers.Search)
	http.Handle("/", rtr)
	http.ListenAndServe(":8000", nil)
}
