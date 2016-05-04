package app

import (
	"strings"
)

type PlatformToTrainNode struct {
	Line     string `xml:"trainName"`
	Platform string `xml:"platformToTrainSteps"`
}

type EntranceNode struct {
	Name             string                `xml:"name"`
	PlatformToTrains []PlatformToTrainNode `xml:"platformToTrain"`
}

type FacilityNode struct {
	Name  string `xml:"name,attr"`
	Count string `xml:",chardata"`
}

type StationNode struct {
	Id          string         `xml:"id,attr"`
	Type        string         `xml:"type,attr"`
	Name        string         `xml:"name"`
	Address     string         `xml:"contactDetails>address"`
	Telephone   string         `xml:"contactDetails>phone"`
	Lines       []string       `xml:"servingLines>servingLine"`
	Zones       []string       `xml:"zones>zone"`
	Facilities  []FacilityNode `xml:"facilities>facility"`
	Entrances   []EntranceNode `xml:"entrances>entrance"`
	Coordinates string         `xml:"Placemark>Point>coordinates"`
}

type DataNode struct {
	PublishedDate string        `xml:"Header>PublishDateTime"`
	Language      string        `xml:"Header>Language"`
	Stations      []StationNode `xml:"stations>station"`
}

func (s *StationNode) Longitude() string {
	return strings.Split(strings.Trim(s.Coordinates, "\n\t"), ",")[0]
}

func (s *StationNode) Latitude() string {
	return strings.Split(strings.Trim(s.Coordinates, "\n\t"), ",")[1]
}

type LatLonNode struct {
	Coordinates []float64
	Type        string
}

type StationDocument struct {
	Id         string
	Name       string
	Address    string
	Telephone  string
	Lines      []string
	Zones      []string
	Facilities []FacilityNode
	Entrances  []EntranceNode
	Location   LatLonNode
}
