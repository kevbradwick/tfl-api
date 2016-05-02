package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevbradwick/tflapi/lib"
	"github.com/kevbradwick/tflapi/query"
	"gopkg.in/mgo.v2/bson"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Total    int           `json:"total"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []lib.Station `json:"results"`
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
	val := ""
	page := 1
	limit := 10
	offset := 0

	// build up the final query in this var
	q := bson.M{}

	// the original query string, parsed
	v, _ := url.ParseQuery(r.URL.RawQuery)

	if val = v.Get("name"); val != "" {
		q["name"] = val
	}

	if val = v.Get("zones"); val != "" {
		q["zones"] = bson.M{"$in": strings.Split(val, ",")}
	}

	if val = v.Get("lines"); val != "" {
		q["lines"] = bson.M{"$in": strings.Split(val, ",")}
	}

	// geospatial search
	if val = v.Get("near"); val != "" {
		coords := strings.Split(val, ",")
		lon, lonErr := strconv.ParseFloat(coords[0], 64)
		lat, latErr := strconv.ParseFloat(coords[1], 64)
		if latErr == nil && lonErr == nil {
			geoPoint := bson.M{"type": "Point", "coordinates": []float64{lon, lat}}
			q["location"] = bson.M{"$near": bson.M{"$geometry": geoPoint}}
		}
	}

	if val = v.Get("page"); val != "" {
		page, _ = strconv.Atoi(val)
	}

	count, _ := query.Count(q) // total number of documents returned
	var previousPage string
	var nextPage string
	totalPages := math.Ceil(float64(count) / float64(limit))

	if page > 1 {
		qs := r.URL.Query()
		qs.Set("page", strconv.Itoa(page-1))
		previousPage = fmt.Sprintf("/station/search?%s", qs.Encode())
	}

	if (float64(page + 1)) <= totalPages {
		qs := r.URL.Query()
		qs.Set("page", strconv.Itoa(page+1))
		nextPage = fmt.Sprintf("/station/search?%s", qs.Encode())
	}

	offset = (page * limit) - limit
	stations, err := query.FindMany(q, limit, offset)
	response := &ListResponse{count, nextPage, previousPage, stations}

	if err != nil || float64(page) > totalPages {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{404, "No results found"})
		return
	}

	json.NewEncoder(w).Encode(response)
}
