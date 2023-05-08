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

// SensorsHandler handles requests to /sensors.
func (api *SensorAPI) SensorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
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

		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Invalid HTTP request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// SensorHandler handles requests to /sensors/{name}.
func (api *SensorAPI) SensorHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/sensors/")
	log.Debug("name: ", name)
	switch r.Method {
	case http.MethodGet:
		sensor, code, err := api.store.GetSensor(name)
		if err != nil {
			log.Error("Failed to get sensor: ", err)
			http.Error(w, "Failed to get sensor", code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sensor)
	case http.MethodPut:
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
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		code, err := api.store.RemoveSensor(name)
		if err != nil {
			log.Error("Failed to remove sensor: ", err)
			http.Error(w, fmt.Sprint("Failed to remove sensor: ", err), code)
			return
		}
		w.WriteHeader(http.StatusOK)
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

// NearestSensorHandler handles requests to /nearest-sensor.
func (api *SensorAPI) NearestSensorHandler(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		log.Error("Failed to parse latitude from URL: ", err)
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		log.Error("Failed to parse longitude from URL: ", err)
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	location := model.Location{
		Latitude:  lat,
		Longitude: lon,
	}
	switch r.Method {
	case http.MethodGet:
		nearestSensor, code, err := api.store.GetNearestSensor(location)
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
