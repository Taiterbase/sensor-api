package main

import (
	"net/http"
	"sensor-api/internal/api"
	"sensor-api/internal/store"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	sensorStore := store.NewInMemorySensorStore()
	sensorAPI := api.NewSensorAPI(sensorStore)
	timeout := 5 * time.Second
	http.Handle("/sensors", api.TimeoutMiddleware(timeout, http.HandlerFunc(sensorAPI.SensorsHandler)))
	http.Handle("/sensors/", api.TimeoutMiddleware(timeout, http.HandlerFunc(sensorAPI.SensorHandler)))
	http.Handle("/sensors/nearest", api.TimeoutMiddleware(timeout, http.HandlerFunc(sensorAPI.NearestSensorHandler)))
	http.Handle("/sensors/tags", api.TimeoutMiddleware(timeout, http.HandlerFunc(sensorAPI.TagsHandler)))
	http.Handle("/sensors/locations", api.TimeoutMiddleware(timeout, http.HandlerFunc(sensorAPI.LocationsHandler)))

	log.Info("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
