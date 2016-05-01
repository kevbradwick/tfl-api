package main

import (
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/lib"
	"net/http"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/station/{id:[\\d]+}", lib.GetStation)
	http.Handle("/", rtr)
	http.ListenAndServe(":8000", nil)
}
