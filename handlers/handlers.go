package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/query"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
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
//
// The following query parameters are permitted;
//
//	name
//		Search for an exact match
//
//	zones
//		Query multiple zones e.g. ?zones=3,4
//
//	lines
//		Query multiple lines e.g. ?lines=District,Circle
//
func Search(w http.ResponseWriter, r *http.Request) {
	var val string
	q := bson.M{}

	if val = r.URL.Query().Get("name"); val != "" {
		q["name"] = val
	}

	// could be multiple
	if val = r.URL.Query().Get("zones"); val != "" {
		q["zones"] = bson.M{"$in": strings.Split(val, ",")}
	}

	if val = r.URL.Query().Get("lines"); val != "" {
		q["lines"] = bson.M{"$in": strings.Split(val, ",")}
	}

	if val = r.URL.Query().Get("near"); val != "" {
		coords := strings.Split(val, ",")
		lon, _ := strconv.ParseFloat(coords[0], 64)
		lat, _ := strconv.ParseFloat(coords[1], 64)
		q["location"] = bson.M{"$near": bson.M{"$geometry": bson.M{"type": "Point", "coordinates": []float64{lon, lat}}}}
	}

	stations, err := query.FindMany(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{404, "Station not found"})
		return
	}

	json.NewEncoder(w).Encode(stations)
}
