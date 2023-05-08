package store

import (
	"fmt"
	"math"
	"net/http"
	"sensor-api/internal/model"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/rtree"
)

// InMemorySensorStore is an in-memory implementation of SensorStore.
type InMemorySensorStore struct {
	mu      sync.Mutex
	sensors map[string]model.Sensor
	// UC Berkeley's RTree implementation
	rt *rtree.RTreeGN[float64, string]
}

// NewInMemorySensorStore creates a new InMemorySensorStore.
func NewInMemorySensorStore() *InMemorySensorStore {
	return &InMemorySensorStore{
		sensors: make(map[string]model.Sensor),
		rt:      &rtree.RTreeGN[float64, string]{},
	}
}

// AddSensor adds a sensor to the store.
func (store *InMemorySensorStore) AddSensor(sensor model.Sensor) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, exists := store.sensors[sensor.Name]
	if exists {
		log.Error("Sensor already exists: ", sensor.Name)
		return http.StatusBadRequest, fmt.Errorf("sensor already exists")
	}

	store.sensors[sensor.Name] = sensor

	point := [2]float64{sensor.Location.Latitude, sensor.Location.Longitude}
	store.rt.Insert(point, point, sensor.Name)
	return http.StatusCreated, nil
}

func (store *InMemorySensorStore) GetSensor(name string) (model.Sensor, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	sensor, ok := store.sensors[name]
	if !ok {
		log.Error("Sensor not found: ", name)
		return model.Sensor{}, http.StatusNotFound, fmt.Errorf("sensor not found")
	}

	return sensor, http.StatusOK, nil
}

// UpdateSensor updates a sensor in the store.
func (store *InMemorySensorStore) UpdateSensor(name string, updatedSensor *model.Sensor) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if updatedSensor == nil {
		log.Error("Sensor is nil")
		return http.StatusBadRequest, fmt.Errorf("sensor is nil")
	}

	if updatedSensor.Name == "" {
		log.Error("Sensor name is empty")
		return http.StatusBadRequest, fmt.Errorf("sensor name is empty")
	}

	sensor, ok := store.sensors[name]
	if !ok {
		log.Error("Sensor not found: ", name)
		return http.StatusNotFound, fmt.Errorf("sensor not found")
	}

	// if the name changed, update the sensor name in the rtree
	oldPoint := [2]float64{sensor.Location.Latitude, sensor.Location.Longitude}
	newPoint := [2]float64{updatedSensor.Location.Latitude, updatedSensor.Location.Longitude}
	if oldPoint != newPoint || sensor.Name != updatedSensor.Name {
		// only update the rtree if name or location data has been updated
		store.rt.Replace(oldPoint, oldPoint, sensor.Name, newPoint, newPoint, updatedSensor.Name)
	}
	delete(store.sensors, sensor.Name)
	store.sensors[updatedSensor.Name] = *updatedSensor
	return http.StatusNoContent, nil
}

// RemoveSensor removes a sensor from the store.
func (store *InMemorySensorStore) RemoveSensor(name string) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	sensor, ok := store.sensors[name]
	if !ok {
		log.Error("Sensor not found: ", name)
		return http.StatusNotFound, fmt.Errorf("sensor not found")
	}

	point := [2]float64{sensor.Location.Latitude, sensor.Location.Longitude}
	store.rt.Delete(point, point, sensor.Name)
	delete(store.sensors, name)

	return http.StatusNoContent, nil
}

// GetNearestSensor returns the nearest sensor to the given location.
func (store *InMemorySensorStore) GetNearestSensor(location model.Location) (*model.Sensor, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if location == (model.Location{}) {
		log.Error("Location is nil")
		return nil, http.StatusBadRequest, fmt.Errorf("location is nil")
	}

	if location.Latitude < -90 || location.Latitude > 90 {
		log.Error("Latitude is out of range: ", location.Latitude)
		return nil, http.StatusBadRequest, fmt.Errorf("latitude is out of range")
	}

	if location.Longitude < -180 || location.Longitude > 180 {
		log.Error("Longitude is out of range: ", location.Longitude)
		return nil, http.StatusBadRequest, fmt.Errorf("longitude is out of range")
	}

	if len(store.sensors) == 0 || store.rt.Len() == 0 {
		log.Error("No sensors in store")
		return nil, http.StatusNotFound, fmt.Errorf("no sensors in store")
	}

	if len(store.sensors) == 1 {
		for _, sensor := range store.sensors {
			return &sensor, http.StatusOK, nil
		}
	}

	var (
		closestSensor   *model.Sensor
		closestDistance float64 = math.MaxFloat64
	)

	point := [2]float64{location.Latitude, location.Longitude}
	store.rt.Nearby(
		/* func(min, max [2]float64, data string, item bool) float64 {
			return haversineDistance(min[1], min[0], point[1], point[0])
		}, */
		rtree.BoxDist[float64, string](point, point, nil),
		func(min, max [2]float64, data string, dist float64) bool {
			sensor := store.sensors[data]
			fmt.Println(location, sensor, dist)
			if closestSensor == nil || dist < closestDistance {
				closestSensor = &sensor
				closestDistance = dist
			}

			return true
		},
	)

	return closestSensor, http.StatusOK, nil
}

/*
Commenting out for test coverage
// haversineDistance calculates the distance between two points on Earth using the Haversine formula.
// https://en.wikipedia.org/wiki/Haversine_formula
// Good for short distances, less accurate for larger distances.
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Radius of Earth in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	latDiff := lat2Rad - lat1Rad
	lonDiff := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(latDiff/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(lonDiff/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}
*/
