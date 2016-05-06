package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
	"strings"
)

// ListResponse is a struct for rendering responses with many results. It
// contains next and previous links as well as the total count of records for
// the given query.
type ListResponse struct {
	Total    int       `json:"total"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Station `json:"results"`
}

// PanicHandler will recover from panics and return the correct http status for
// the given error. If the error is not an instance of HttpError then a 500
// status will be returned.
func PanicHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e := err.(HttpError); e != nil {
					c.JSON(e.StatusCode(), gin.H{
						"error":   true,
						"message": e.Message(),
					})
				} else {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "An error occured. Please try again later",
					})
				}
			}
		}()
		c.Next()
	}
}

// GetStationHandler will get a single tube station using the given id in the
// path parameter.
func GetStationHandler(c *gin.Context) {
	if station, err := FindOne(bson.M{"id": c.Param("id")}); err != nil {
		if err == mgo.ErrNotFound {
			panic(&GenericHttpError{404, "Not found."})
		} else {
			panic(err)
		}
	} else {
		c.JSON(200, station)
	}
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
func SearchHandler(c *gin.Context) {
	val := ""
	page := 1
	limit := 10
	offset := 0

	// build up the final query in this var
	q := bson.M{}

	if val = c.Query("name"); val != "" {
		q["name"] = val
	}

	if val = c.Query("zones"); val != "" {
		q["zones"] = bson.M{"$in": strings.Split(val, ",")}
	}

	if val = c.Query("lines"); val != "" {
		q["lines"] = bson.M{"$in": strings.Split(val, ",")}
	}

	// geospatial search
	if val = c.Query("near"); val != "" {
		coords := strings.Split(val, ",")
		lon, lonErr := strconv.ParseFloat(coords[0], 64)
		lat, latErr := strconv.ParseFloat(coords[1], 64)
		if latErr == nil && lonErr == nil {
			geoPoint := bson.M{"type": "Point", "coordinates": []float64{lon, lat}}
			q["location"] = bson.M{"$near": bson.M{"$geometry": geoPoint}}
		}
	}

	if val = c.Query("page"); val != "" {
		page, _ = strconv.Atoi(val)
	}

	count, err := Count(q) // total number of documents returned
	if err != nil {
		panic(err)
	}

	var previousPage string
	var nextPage string
	totalPages := math.Ceil(float64(count) / float64(limit))

	if page > 1 {
		qs := c.Request.URL.Query()
		qs.Set("page", strconv.Itoa(page-1))
		previousPage = fmt.Sprintf("/station/search?%s", qs.Encode())
	}

	if (float64(page + 1)) <= totalPages {
		qs := c.Request.URL.Query()
		qs.Set("page", strconv.Itoa(page+1))
		nextPage = fmt.Sprintf("/stations/search?%s", qs.Encode())
	}

	offset = (page * limit) - limit
	stations, err := FindMany(q, limit, offset)
	response := &ListResponse{count, nextPage, previousPage, stations}

	if err != nil {
		panic(err)
	}

	c.JSON(200, response)
}

// Lines handler will list all the lines tube stations are on.
func LinesHandler(c *gin.Context) {
	lines, err := DistinctQuery("lines")
	if err != nil {
		panic(err)
	}
	c.JSON(200, lines)
}

// ZonesHandler will list all the zones the tube stations runs through.
func ZonesHandler(c *gin.Context) {
	zones, err := DistinctQuery("zones")
	if err != nil {
		panic(err)
	}
	c.JSON(200, zones)
}
