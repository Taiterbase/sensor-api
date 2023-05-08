package model

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *Location) Equals(other Location) bool {
	return l.Latitude == other.Latitude && l.Longitude == other.Longitude
}

func (l *Location) IsValid() bool {
	return *l != Location{} && l.Latitude >= -90 && l.Latitude <= 90 && l.Longitude >= -180 && l.Longitude <= 180
}
