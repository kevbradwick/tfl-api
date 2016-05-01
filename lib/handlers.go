package lib

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
)

func GetStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Request for station %q", vars["id"])
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to mongodb")
		os.Exit(2)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("tfldata").C("tube_stations")
	station := &Station{}
	query := collection.Find(bson.M{"id": vars["id"]})
	count, _ := query.Count()
	log.Println("Found %d documents", count)
	err = query.One(&station)
	if err != nil {
		log.Println(err)
	}
	log.Println(station)
	json.NewEncoder(w).Encode(station)
}
