package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/query"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Total   int           `json:"total"`
	Results []interface{} `json:"results"`
}

func GetStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	station, err := query.FindOne(bson.M{"id": vars["id"]})

	// err can only be a 404
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{404, "Station not found"})
		return
	}

	json.NewEncoder(w).Encode(station)
}

// Search
//
// The search handler will take all allowed query parameters and build
// a boolean (AND) query and return a paginated response.
func Search(w http.ResponseWriter, r *http.Request) {
	qName := r.URL.Query().Get("name")
	q := bson.M{}

	if qName != "" {
		q["name"] = qName
	}

	stations, err := query.FindMany(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{404, "Station not found"})
		return
	}

	json.NewEncoder(w).Encode(stations)
}
