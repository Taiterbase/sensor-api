package store

import "sensor-api/internal/model"

type SensorStore interface {
	AddSensor(sensor model.Sensor) (int, error)
	GetSensor(name string) (model.Sensor, int, error)
	UpdateSensor(name string, updatedSensor *model.Sensor) (int, error)
	RemoveSensor(name string) (int, error)
	GetNearestSensor(location model.Location) (*model.Sensor, int, error)
}
