package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sensor-api/internal/model"
	"sensor-api/internal/store"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SensorAPI struct {
	store store.SensorStore
}

// NewSensorAPI creates a new SensorAPI.
func NewSensorAPI(store store.SensorStore) *SensorAPI {
	return &SensorAPI{
		store: store,
	}
}

// SensorHandler handles requests to /sensors/{name}.
func (api *SensorAPI) SensorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := r.URL.Path
		name := strings.TrimPrefix(path, "/sensors/")
		sensor, code, err := api.store.GetSensor(name)
		if err != nil {
			log.Error("Failed to get sensor: ", err)
			http.Error(w, "Failed to get sensor", code)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(sensor)
	case http.MethodPut:
		path := r.URL.Path
		name := strings.TrimPrefix(path, "/sensors/")
		var updatedSensor model.Sensor
		err := json.NewDecoder(r.Body).Decode(&updatedSensor)
		if err != nil {
			log.Error("Failed to decode request body: ", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		code, err := api.store.UpdateSensor(name, &updatedSensor)
		if err != nil {
			log.Error("Failed to update sensor: ", err)
			http.Error(w, fmt.Sprint("Failed to update sensor: ", err), code)
			return
		}
		log.Info("Updated sensor: ", updatedSensor)
		w.WriteHeader(code)
	case http.MethodDelete:
		path := r.URL.Path
		name := strings.TrimPrefix(path, "/sensors/")
		code, err := api.store.RemoveSensor(name)
		if err != nil {
			log.Error("Failed to remove sensor: ", err)
			http.Error(w, fmt.Sprint("Failed to remove sensor: ", err), code)
			return
		}

		w.WriteHeader(code)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// SensorHandler handles requests to /sensors/{name}.
func (api *SensorAPI) SensorsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("request URI: ", r.RequestURI)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("count") == "true" {
			// an optional parameter to get the number of sensors
			count, code, err := api.store.GetSensorCount()
			if err != nil {
				log.Error("Failed to get sensor count: ", err)
				http.Error(w, "Failed to get sensor count", code)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(count)
			return
		}

		tags := r.URL.Query()["tags"]
		// if tags is nil, GetSensorsByTags will return all sensors
		sensor, code, err := api.store.GetSensorsByTags(tags)
		if err != nil {
			log.Error("Failed to get sensor: ", err)
			http.Error(w, "Failed to get sensor", code)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(sensor)
	case http.MethodPost:
		var sensor model.Sensor
		err := json.NewDecoder(r.Body).Decode(&sensor)
		if err != nil {
			log.Error("Failed to decode request body: ", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		code, err := api.store.AddSensor(sensor)
		if err != nil {
			log.Error("Failed to add sensor: ", err)
			http.Error(w, fmt.Sprint("Failed to add sensor: ", err), code)
			return
		}

		log.Info("Added sensor: ", sensor)
		w.WriteHeader(code)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// NearestSensorHandler handles requests to /nearest-sensor.
func (api *SensorAPI) NearestSensorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		lat, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
		if err != nil {
			log.Error("Failed to parse latitude from URL: ", err)
			http.Error(w, "Invalid latitude", http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
		if err != nil {
			log.Error("Failed to parse longitude from URL: ", err)
			http.Error(w, "Invalid longitude", http.StatusBadRequest)
			return
		}
		location := model.Location{
			Latitude:  lat,
			Longitude: lon,
		}

		log.Debug("nearest sensor endpoint")
		log.Debug("location: ", location)
		tags := r.URL.Query()["tags"]
		// if tags is nil, GetNearestSensorByTag will return the nearest sensor regardless of tags
		nearestSensor, code, err := api.store.GetNearestSensorByTag(location, tags)
		if err != nil {
			log.Error("Failed to get nearest sensor: ", err)
			http.Error(w, "Failed to get nearest sensor", code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(nearestSensor)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *SensorAPI) TagsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tags, code, err := api.store.GetUniqueTags()
		if err != nil {
			log.Error("Failed to get tags: ", err)
			http.Error(w, "Failed to get tags", code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(tags)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *SensorAPI) LocationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		locations, code, err := api.store.GetUniqueLocations()
		if err != nil {
			log.Error("Failed to get locations: ", err)
			http.Error(w, "Failed to get locations", code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(locations)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
