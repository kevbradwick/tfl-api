package app

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
	"strings"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Total    int       `json:"total"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Station `json:"results"`
}

func PanicHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(404, gin.H{"status": 500, "error": "Internal server error"})
			}
		}()
		c.Next()
	}
}

func GetStationHandler(c *gin.Context) {
	if station, err := FindOne(bson.M{"id": c.Param("id")}); err != nil {
		panic(err)
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

	count, _ := Count(q) // total number of documents returned
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

	// TODO move this higher up, outside of the handler function
	if err != nil {
		if err == mgo.ErrNotFound {
			panic(&HttpError{})
		}
	}
	c.JSON(200, response)
}
