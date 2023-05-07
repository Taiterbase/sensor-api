package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sensor-api/internal/model"
	"sensor-api/internal/store"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSensorsHandler(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194}}`)
	req, err := http.NewRequest("POST", "/sensors", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestAddSensorsHandlerInvalidBody(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new bad request to add a sensor, missing closing bracket
	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194}`)
	req, err := http.NewRequest("POST", "/sensors", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAddSensorsHandlerExists(t *testing.T) {
	// Create a new in-memory store and add a sensor to it
	store := store.NewInMemorySensorStore()
	sensor := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
	code, err := store.AddSensor(sensor)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, code)

	// Create a new request to add a sensor
	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194}}`)
	req, err := http.NewRequest("POST", "/sensors", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSensorsHandlerInvalidMethod(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("PUT", "/sensors", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestSensorsHandlerOptions(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("OPTIONS", "/sensors", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the Allow header is what we expect
	assert.Equal(t, "POST, OPTIONS", recorder.Header().Get("Allow"))
}

func TestSensorsHandlerHead(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("HEAD", "/sensors", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorsHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRemoveSensorHandler(t *testing.T) {
	// Create a new in-memory store and add a sensor to it
	store := store.NewInMemorySensorStore()
	sensor := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
	code, err := store.AddSensor(sensor)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, code)

	// Create a new request to remove the sensor
	req, err := http.NewRequest("DELETE", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the RemoveSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)
	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)
	handler.ServeHTTP(recorder, req)
}

func TestRemoveSensorHandlerNotFound(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to remove the sensor
	req, err := http.NewRequest("DELETE", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the RemoveSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetSensorHandler(t *testing.T) {
	// Create a new in-memory store and add a sensor to it
	store := store.NewInMemorySensorStore()
	sensor := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
	code, err := store.AddSensor(sensor)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, code)

	// Create a new request to get the sensor
	req, err := http.NewRequest("GET", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the GetSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body is what we expect
	expectedBody := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194},"tags":null}`)
	assert.Equal(t, expectedBody, recorder.Body.String())
}

func TestGetSensorHandlerNotFound(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the sensor
	req, err := http.NewRequest("GET", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the GetSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

// update a sensor
func TestUpdateSensorHandler(t *testing.T) {
	// Create a new in-memory store and add a sensor to it
	store := store.NewInMemorySensorStore()
	sensor := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
	code, err := store.AddSensor(sensor)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, code)

	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194},"tags":["tag1","tag2"]}`)
	// Create a new request to update the sensor
	req, err := http.NewRequest("PUT", "/sensors/Sensor1", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the UpdateSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)
	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	getReq, err := http.NewRequest("GET", "/sensors/Sensor1", nil)
	assert.NoError(t, err)
	handler.ServeHTTP(recorder, getReq)
	// Check the response is what we expect, with additional tags
	assert.Equal(t, body, recorder.Body.String())
}

// update a sensor that doesn't exist
func TestUpdateSensorHandlerNotFound(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194},"tags":["tag1","tag2"]}`)
	// Create a new request to update the sensor
	req, err := http.NewRequest("PUT", "/sensors/Sensor1", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the UpdateSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)
	// Check the status code is what we expect
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

// update a sensor with invalid request body
func TestUpdateSensorHandlerInvalidBody(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	body := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194},"tags":["tag1","tag2"}`)
	// Create a new request to update the sensor
	req, err := http.NewRequest("PUT", "/sensors/Sensor1", strings.NewReader(body))
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the UpdateSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)
	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSensorHandlerInvalidMethod(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("POST", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestSensorHandlerOptions(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("OPTIONS", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the Allow header is what we expect
	assert.Equal(t, "GET, PUT, DELETE, OPTIONS", recorder.Header().Get("Allow"))
}

func TestSensorHandlerHead(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to add a sensor
	req, err := http.NewRequest("HEAD", "/sensors/Sensor1", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the AddSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).SensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestNearestSensorHandler(t *testing.T) {
	// Create a new in-memory store and add a sensor to it
	store := store.NewInMemorySensorStore()
	sensor1 := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
	code, err := store.AddSensor(sensor1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, code)

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("GET", "/nearest-sensor?latitude=37.7749&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body is what we expect
	expectedBody := fmt.Sprintln(`{"name":"Sensor1","location":{"latitude":37.7749,"longitude":-122.4194},"tags":null}`)
	assert.Equal(t, expectedBody, recorder.Body.String())
}

func TestNearestSensorHandlerInvalidLatitude(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("GET", "/nearest-sensor?latitude=invalid&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestNearestSensorHandlerInvalidLongitude(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("GET", "/nearest-sensor?latitude=37.7749&longitude=invalid", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestNearestSensorHandlerInvalidMethod(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("POST", "/nearest-sensor?latitude=37.7749&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestNearestSensorHandlerOptions(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("OPTIONS", "/nearest-sensor?latitude=37.7749&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the Allow header is what we expect
	assert.Equal(t, "GET, OPTIONS", recorder.Header().Get("Allow"))
}

func TestNearestSensorHandlerHead(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("HEAD", "/nearest-sensor?latitude=37.7749&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestNearestSensorHandlerNotFound(t *testing.T) {
	// Create a new in-memory store
	store := store.NewInMemorySensorStore()

	// Create a new request to get the nearest sensor
	req, err := http.NewRequest("GET", "/nearest-sensor?latitude=37.7749&longitude=-122.4194", nil)
	assert.NoError(t, err)

	// Create a new recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the NearestSensor handler with the request and recorder
	handler := http.HandlerFunc(NewSensorAPI(store).NearestSensorHandler)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
