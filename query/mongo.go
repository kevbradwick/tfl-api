package query

import (
	"github.com/kevbradwick/tflapi/lib"
	"gopkg.in/mgo.v2"
	"log"
)

// Create a new MongoDB session.
//
// It's the calling function's responsibility to close the session
func session() (s *mgo.Session) {
	s, err := mgo.Dial(lib.MongoUrl())
	if err != nil {
		log.Fatal("Unable to connect to Mongo instance.")
	}
	s.SetMode(mgo.Monotonic, true)
	return s
}

// Find a station.
//
// Return a single tube station
func FindOne(query interface{}) (station lib.Station, err error) {
	s := session()
	defer s.Close()
	c := s.DB("tfldata").C("tube_stations")
	station = lib.Station{}
	err = c.Find(query).One(&station)

	// some other error happened??
	if err != nil && err != mgo.ErrNotFound {
		log.Fatal("MongoDB query failed. %q", err)
	}

	return station, err
}

// Count the number of documents in a query
func Count(query interface{}) (count int, err error) {
	s := session()
	defer s.Close()
	c := s.DB("tfldata").C("tube_stations")
	count, err = c.Find(query).Count()
	if err != nil {
		log.Printf("Mongo query error %q", err)
	}
	return count, err
}

func FindMany(query interface{}, limit, offset int) (stations []lib.Station, err error) {
	if limit < 0 || limit > 30 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	s := session()
	defer s.Close()
	c := s.DB("tfldata").C("tube_stations")
	stations = []lib.Station{}
	err = c.Find(query).Limit(limit).All(&stations)

	// some other error happened??
	if err != nil && err != mgo.ErrNotFound {
		log.Fatal("MongoDB query failed. %q", err)
	}

	return stations, err
}
