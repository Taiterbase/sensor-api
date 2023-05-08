package main

import (
	"net/http"
	"sensor-api/internal/api"
	"sensor-api/internal/store"

	log "github.com/sirupsen/logrus"
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
		http.HandleFunc("/sensors/tags/within", sensorAPI.getSensorsByTagWithinBoundingBox)
		http.HandleFunc("/sensors/tags/near", sensorAPI.getSensorsByTagWithinRadius)
		http.HandleFunc("/sensors/count", sensorAPI.getTotalSensors)
		http.HandleFunc("/sensors/tags", sensorAPI.getUniqueTags)
		http.HandleFunc("/sensors/locations", sensorAPI.getUniqueLocations)
		http.HandleFunc("/sensors/locations/count", sensorAPI.getTotalUniqueLocations)
	*/
	log.Info("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
