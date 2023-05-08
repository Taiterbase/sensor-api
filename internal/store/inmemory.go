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
	// mapping of tag name to sensor names
	tags map[string]map[string]struct{}
	// UC Berkeley's RTree implementation
	rt *rtree.RTreeGN[float64, string]
}

// NewInMemorySensorStore creates a new InMemorySensorStore.
func NewInMemorySensorStore() *InMemorySensorStore {
	return &InMemorySensorStore{
		sensors: make(map[string]model.Sensor),
		rt:      &rtree.RTreeGN[float64, string]{},
		tags:    make(map[string]map[string]struct{}),
	}
}

// AddSensor adds a sensor to the store.
func (store *InMemorySensorStore) AddSensor(sensor model.Sensor) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if !sensor.Location.IsValid() {
		log.Error("Invalid location: ", sensor.Location)
		return http.StatusBadRequest, fmt.Errorf("invalid location")
	}

	_, exists := store.sensors[sensor.Name]
	if exists {
		log.Error("Sensor already exists: ", sensor.Name)
		return http.StatusBadRequest, fmt.Errorf("sensor already exists")
	}

	// add sensor to store
	store.sensors[sensor.Name] = sensor

	// insert the sensor into the rtree
	point := [2]float64{sensor.Location.Latitude, sensor.Location.Longitude}
	store.rt.Insert(point, point, sensor.Name)

	// add sensor name to tags
	for _, tag := range sensor.Tags {
		if store.tags[tag] == nil {
			store.tags[tag] = make(map[string]struct{})
		}
		store.tags[tag][sensor.Name] = struct{}{}
	}

	return http.StatusCreated, nil
}

// GetSensor returns a sensor from the store.
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

// GetSensors returns all sensors in the store.
func (store *InMemorySensorStore) GetSensors() ([]model.Sensor, int, error) {
	return store.GetSensorsByTags(nil)
}

// GetSensorsByTag returns all sensors with the given tags.
func (store *InMemorySensorStore) GetSensorsByTags(tags []string) ([]model.Sensor, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	log.Debug("Getting sensors by tags: ", tags)

	if len(store.sensors) == 0 {
		log.Error("No sensors in store")
		return nil, http.StatusNotFound, fmt.Errorf("no sensors in store")
	}

	// group up unique tags
	uniqueTags := make(map[string]bool)
	for _, tag := range tags {
		uniqueTags[tag] = true
	}

	sensors := []model.Sensor{}

	if len(uniqueTags) == 0 {
		// if tags is empty, return all sensors
		for _, sensor := range store.sensors {
			sensors = append(sensors, sensor)
		}
		log.Debug("Returning all sensors", sensors)
		return sensors, http.StatusOK, nil
	}

	// get all sensors that have all the given tags, ANDing the tags
	sensorTagCount := make(map[string]int)
	uniqueSensors := map[string]struct{}{}
	for tag := range uniqueTags {
		sensorNames, exists := store.tags[tag]
		if exists {
			for sensorName := range sensorNames {
				sensorTagCount[sensorName]++
			}
		}
	}

	for tag := range uniqueTags {
		sensorNames, exists := store.tags[tag]
		if exists {
			for sensorName := range sensorNames {
				if sensorTagCount[sensorName] == len(uniqueTags) {
					uniqueSensors[sensorName] = struct{}{}
				}
			}
		}
	}

	// grab the unique sensors from the store
	for sensorName := range uniqueSensors {
		sensors = append(sensors, store.sensors[sensorName])
	}

	return sensors, http.StatusOK, nil
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

	if !updatedSensor.Location.IsValid() {
		log.Error("Invalid location: ", updatedSensor.Location)
		return http.StatusBadRequest, fmt.Errorf("invalid location")
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
	// remove old sensor from store
	delete(store.sensors, sensor.Name)
	// add updated sensor to store
	store.sensors[updatedSensor.Name] = *updatedSensor

	// update sensor name in tags
	for _, tag := range sensor.Tags {
		delete(store.tags[tag], sensor.Name)
		if len(store.tags[tag]) == 0 {
			delete(store.tags, tag)
		}
	}
	for _, tag := range updatedSensor.Tags {
		if store.tags[tag] == nil {
			store.tags[tag] = make(map[string]struct{})
		}
		store.tags[tag][updatedSensor.Name] = struct{}{}
	}

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

	// remove sensor name from tags
	for _, tag := range sensor.Tags {
		delete(store.tags[tag], name)
		// if there are no more sensors with this tag, remove the tag
		if len(store.tags[tag]) == 0 {
			delete(store.tags, tag)
		}
	}

	return http.StatusNoContent, nil
}

// GetNearestSensor returns the nearest sensor to the given location.
func (store *InMemorySensorStore) GetNearestSensor(location model.Location) (*model.Sensor, int, error) {
	return store.GetNearestSensorByTag(location, nil)
}

// GetNearestSensorByTag returns the nearest sensor to the given location with the given set of tags.
func (store *InMemorySensorStore) GetNearestSensorByTag(location model.Location, tags []string) (*model.Sensor, int, error) {
	log.Debug("Getting nearest sensor by tag: ", tags)
	if !location.IsValid() {
		log.Error("Invalid location: ", location)
		return nil, http.StatusBadRequest, fmt.Errorf("invalid location")
	}
	if len(store.sensors) == 0 {
		log.Error("No sensors in store")
		return nil, http.StatusNotFound, fmt.Errorf("no sensors in store")
	}

	sensors, code, err := store.GetSensorsByTags(tags)
	if err != nil {
		return nil, code, err
	}
	if len(sensors) == 0 {
		log.Error("No sensors with given tag(s)")
		return nil, http.StatusNotFound, fmt.Errorf("no sensors with given tag(s)")
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	// make a mapping of sensor name to sensor
	sensorMap := make(map[string]model.Sensor)
	for _, sensor := range sensors {
		sensorMap[sensor.Name] = sensor
	}

	var (
		closestSensor   *model.Sensor
		closestDistance float64 = math.MaxFloat64
	)

	point := [2]float64{location.Latitude, location.Longitude}
	log.Debug("Starting location: ", point)
	store.rt.Nearby(
		/* func(min, max [2]float64, data string, item bool) float64 {
			return haversineDistance(min[1], min[0], point[1], point[0])
		}, */
		rtree.BoxDist[float64, string](point, point, nil),
		func(min, max [2]float64, data string, dist float64) bool {
			if sensor, ok := sensorMap[data]; ok {
				log.Debug("Nearby Sensor: ", data, sensor.Location, dist)
				if closestSensor == nil || dist < closestDistance {
					closestSensor = &sensor
					closestDistance = dist
				}
			}

			return true
		},
	)

	return closestSensor, http.StatusOK, nil
}

// GetUniqueTags returns all unique tags in the store.
func (store *InMemorySensorStore) GetUniqueTags() ([]string, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	uniqueTags := make([]string, 0, len(store.tags))
	for tag := range store.tags {
		uniqueTags = append(uniqueTags, tag)
	}
	return uniqueTags, http.StatusOK, nil
}

// GetUniqueLocations returns all unique locations in the store.
func (store *InMemorySensorStore) GetUniqueLocations() ([]model.Location, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	uniqueLocations := make(map[model.Location]struct{})
	for _, sensor := range store.sensors {
		uniqueLocations[sensor.Location] = struct{}{}
	}

	locations := make([]model.Location, 0, len(uniqueLocations))
	for location := range uniqueLocations {
		locations = append(locations, location)
	}

	return locations, http.StatusOK, nil
}

// GetTotalSensors returns the total number of sensors in the store.
func (store *InMemorySensorStore) GetSensorCount() (int, int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	return len(store.sensors), http.StatusOK, nil
}

/*
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

	a := math.Pow(math.Sin(latDiff / 2), 2) + math.Cos(lat1Rad) * math.Cos(lat2Rad) * math.Pow(math.Sin(lonDiff / 2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))

	distance := earthRadius * c
	return distance
}
*/
