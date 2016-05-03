package main

import (
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/app"
	"net/http"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/station/{id:[\\d]+}", app.GetStationHandler)
	rtr.HandleFunc("/station/search", app.SearchHandler)
	http.Handle("/", rtr)
	http.ListenAndServe(":8000", gh.CompressHandler(rtr))
}
