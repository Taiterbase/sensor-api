package store

import "sensor-api/internal/model"

type SensorStore interface {
	AddSensor(sensor model.Sensor) (int, error)
	GetSensor(name string) (model.Sensor, int, error)
	GetSensorsByTags(tags []string) ([]model.Sensor, int, error)
	UpdateSensor(name string, updatedSensor *model.Sensor) (int, error)
	RemoveSensor(name string) (int, error)
	GetNearestSensor(location model.Location) (*model.Sensor, int, error)
	GetNearestSensorByTag(location model.Location, tags []string) (*model.Sensor, int, error)
	GetSensorCount() (int, int, error)
	GetUniqueTags() ([]string, int, error)
	GetUniqueLocations() ([]model.Location, int, error)
	/*
		GetSensorCardinality(tags []string) (int, int, error) ?
		GetSensorsWithinBoundingBox(minLat, minLong, maxLat, maxLong float64) ([]model.Sensor, int, error)
		GetSensorsWithinRadius(location model.Location, radius float64) ([]model.Sensor, int, error)
		GetSensorsByTagWithinBoundingBox(tags []string, minLat, minLong, maxLat, maxLong float64) ([]model.Sensor, int, error)
		GetSensorsByTagWithinRadius(tags []string, location model.Location, radius float64) ([]model.Sensor, int, error)
	*/
}
