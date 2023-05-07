package main

import (
	"log"
	"net/http"
	"sensor-api/internal/api"
	"sensor-api/internal/store"

	"github.com/sirupsen/logrus"
)

func main() {
	sensorStore := store.NewInMemorySensorStore()
	sensorAPI := api.NewSensorAPI(sensorStore)

	http.HandleFunc("/sensors", sensorAPI.SensorsHandler)
	http.HandleFunc("/sensors/", sensorAPI.SensorHandler)
	http.HandleFunc("/nearest-sensor", sensorAPI.NearestSensorHandler)
	/*
		// Functions to add!
		http.HandleFunc("/sensors", sensorAPI.getAllSensors)
		http.HandleFunc("/sensors/tag/{tag}", sensorAPI.getSensorsByTag)
		http.HandleFunc("/sensors/search", sensorAPI.searchSensorsByName)
		http.HandleFunc("/sensors/{name}", sensorAPI.deleteSensorByName)
		http.HandleFunc("/sensors/within", sensorAPI.getSensorsWithinBoundingBox)
		http.HandleFunc("/sensors/near", sensorAPI.getSensorsWithinRadius)
		http.HandleFunc("/sensors/count", sensorAPI.getTotalSensors)
		http.HandleFunc("/sensors/tags", sensorAPI.getUniqueTags)
	*/
	logrus.Info("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
