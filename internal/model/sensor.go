package model

type Sensor struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Tags     []string `json:"tags"`
}
