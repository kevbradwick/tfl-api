package app

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const db string = "tfldata"
const collectionName string = "tube_stations"

// Coordinate
// A struct that represents a GeoJSON object so that we can perform 2d
// spherical searches.
//
// see https://docs.mongodb.org/manual/applications/geospatial-indexes/ for
// more info.
type Coordinate struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Facility struct {
	Name  string `json:"name"`
	Count string `json:"count"`
}

type PlatformToTrain struct {
	Line     string `json:"line"`
	Platform string `json:"platform"`
}

// Entrance
// A station can have one or more entrances.
type Entrance struct {
	Name             string            `json:"name"`
	PlatformToTrains []PlatformToTrain `json:"platform_to_trains"`
}

// Station
// A London tube station struct.
type Station struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Telephone  string     `json:"telephone"`
	Lines      []string   `json:"lines"`
	Zones      []string   `json:"zones"`
	Facilities []Facility `json:"facilities"`
	Entrances  []Entrance `json:"entrances"`
	Location   Coordinate `json:"location"`
}

// Create a new MongoDB session.
//
// It's the calling function's responsibility to close the session
func session() (s *mgo.Session) {
	s, err := mgo.Dial(MongoUrl())
	if err != nil {
		log.Fatal("Unable to connect to Mongo instance.")
	}
	s.SetMode(mgo.Monotonic, true)
	return s
}

func execQuery(q interface{}) (s *mgo.Session, query *mgo.Query) {
	s = session()
	c := s.DB(db).C(collectionName)
	return s, c.Find(q)
}

// Find a station.
//
// Return a single tube station
func FindOne(query interface{}) (station Station, err error) {
	s, q := execQuery(query)
	defer s.Close()
	station = Station{}
	err = q.One(&station)
	return station, err
}

// Count the number of documents in a query
func Count(query interface{}) (count int, err error) {
	s, q := execQuery(query)
	defer s.Close()
	return q.Count()
}

// FindMany
//
// Make a query for multiple documents. The limit by default is 10 but can be
// adjusted to a maximum of 30. The default offset is 0.
func FindMany(query interface{}, limit, offset int) (stations []Station, err error) {
	if limit < 0 || limit > 30 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	s, q := execQuery(query)
	defer s.Close()
	stations = []Station{}
	err = q.Limit(limit).Skip(offset).All(&stations)
	return
}

// DistinctQuery will query the entire collection for distinct values
func DistinctQuery(field string) (values []string, err error) {
	s := session()
	c := s.DB(db).C(collectionName)
	err = c.Find(bson.M{}).Distinct(field, &values)
	return
}
