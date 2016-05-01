package lib

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Facility struct {
	Name  string `json:"name"`
	Count string `json:"count"`
}

type PlatformToTrain struct {
	Line     string `json:"line"`
	Platform string `json:"platform"`
}

type Entrance struct {
	Name             string            `json:"name"`
	PlatformToTrains []PlatformToTrain `json:"platform_to_trains"`
}

type Station struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Telephone  string     `json:"telephone"`
	Lines      []string   `json:"lines"`
	Zones      []string   `json:"zones"`
	Facilities []Facility `json:"facilities"`
	Entrances  []Entrance `json:"entrances"`
	Coordinate Coordinate `json:"coordinate"`
}
