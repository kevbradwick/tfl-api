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

func FindMany(query interface{}) (stations []lib.Station, err error) {
    s := session()
    defer s.Close()
    c := s.DB("tfldata").C("tube_stations")
    stations = []lib.Station{}
    err = c.Find(query).All(&stations)
    // some other error happened??
    if err != nil && err != mgo.ErrNotFound {
		log.Fatal("MongoDB query failed. %q", err)
	}

	return stations, err
}
