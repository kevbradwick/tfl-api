package main

import (
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kevbradwick/tflapi/ingest"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
)

func getString(c *cli.Context, name, defaultVal string) string {
	if c.String(name) != "" {
		return c.String(name)
	}

	return defaultVal
}

// main runner
//
// this will read the input file
func run(c *cli.Context) {
	// check the file exists
	if _, err := os.Stat(c.String("input")); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The file %q does not exist", c.String("input"))
		os.Exit(1)
	}

	dbhost := getString(c, "dbhost", "localhost")
	dbport := getString(c, "dbport", "27017")

	xmlFile, _ := os.Open(c.String("input"))
	defer xmlFile.Close()

	// `d` will contain all the unmarshalled station data
	var d ingest.Data
	xml.NewDecoder(xmlFile).Decode(&d)

	// insert them into mongo
	session, err := mgo.Dial(fmt.Sprintf("%s:%s", dbhost, dbport))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to mongodb")
		os.Exit(2)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("tfldata").C("tube_stations")
	collection.DropCollection()
	for _, s := range d.Stations {
		if s.Type == "tube" {
			lat, _ := strconv.ParseFloat(s.Latitude(), 32)
			lon, _ := strconv.ParseFloat(s.Longitude(), 32)
			latlon := ingest.LatLon{[]float64{lat, lon}, "Point"}
			doc := &ingest.StationDocument{s.Id, s.Name, s.Address, s.Telephone, s.Lines,
				s.Zones, s.Facilities, s.Entrances, latlon}
			err := collection.Insert(doc)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to insert station %q", doc)
				os.Exit(2)
			}
		}
	}
	// create geospatial index
	collection.EnsureIndex(mgo.Index{Key: []string{"$2dsphere:location"}})
}

func main() {
	app := cli.NewApp()
	app.Name = "TFL Stations Facilities Ingest Tool"

	// command line flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input",
			Value: "",
			Usage: "The XML input file",
		},
		cli.StringFlag{
			Name:  "dbhost",
			Value: "localhost",
			Usage: "The MongoDB Hostname",
		},
		cli.StringFlag{
			Name:  "dbport",
			Value: "27017",
			Usage: "The MongoDB port",
		},
	}

	app.Action = run

	app.Run(os.Args)
}
