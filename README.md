# Sensor API

## Overview

The sensor JSON REST API exposes CRUD, geospatial, and tagged querying functionality to an inmemory sensor store.
The core functionality below has been tested to 100% code coverage.

```
● Storing name, location (gps position), and a list of tags for each sensor.
● Retrieving metadata for an individual sensor by name.
● Updating a sensor’s metadata.
● Querying to find the sensor nearest to a given location.
```

### Model

```json
{
  "name": "sensor-name",
  "location": {
    "latitude": 0.0,
    "longitude": 0.0
  },
  "tags": ["tag1", "tag2"]
}
```

### Endpoints

Below are cURL commands for testing every endpoint. The server runs on port 8080 by default.

1. SensorsHandler (GET, POST, OPTIONS, HEAD)

   - Get all sensors:

   ```
   curl -X GET http://localhost:8080/sensors
   ```

   - Get all sensors with specific tags:

   ```
   curl -X GET "http://localhost:8080/sensors?tags=tag1&tags=tag2"
   ```

   - Get sensor count:

   ```
   curl -X GET "http://localhost:8080/sensors?count=true"
   ```

   - Add a new sensor:

   ```
   curl -X POST -H "Content-Type: application/json" -d '{"name": "sensor1", "location": {"latitude": 12.34, "longitude": 56.78}, "tags": ["tag1", "tag2"]}' http://localhost:8080/sensors
   ```

2. SensorHandler (GET, PUT, DELETE, OPTIONS, HEAD)

   - Get a sensor by name:

   ```
   curl -X GET "http://localhost:8080/sensors/sensor1"
   ```

   - Update a sensor by name:

   ```
   curl -X PUT -H "Content-Type: application/json" -d '{"name": "sensor1", "location": {"latitude": 12.34, "longitude": 56.78}, "tags": ["tag1", "tag2"]}' "http://localhost:8080/sensors/sensor1"
   ```

   - Delete a sensor by name:

   ```
   curl -X DELETE "http://localhost:8080/sensors/sensor1"
   ```

3. NearestSensorHandler (GET, OPTIONS, HEAD)

   - Get nearest sensor by location:

   ```
   curl -X GET "http://localhost:8080/sensors/nearest?latitude=12.34&longitude=56.78"
   ```

   - Get nearest sensor by location and tags:

   ```
   curl -X GET "http://localhost:8080/sensors/nearest?latitude=12.34&longitude=56.78&tags=tag2"
   ```

4. TagsHandler (GET, OPTIONS, HEAD)

   - Get unique tags:

   ```
   curl -X GET http://localhost:8080/sensors/tags
   ```

5. LocationsHandler (GET, OPTIONS, HEAD)

   - Get unique locations:

   ```
   curl -X GET http://localhost:8080/sensors/locations
   ```

### Additional Endpoints:

Here are some additional endpoints I would implement for querying sensor data:

1. Get sensors within a bounding box: Retrieve a list of sensors located within a specified bounding box, defined by minimum and maximum latitude and longitude values.

```

GET /sensors/within?min_lat={min_lat}&max_lat={max_lat}&min_lng={min_lng}&max_lng={max_lng}

```

2. Get sensors within a certain radius: Retrieve a list of sensors located within a specified radius (in meters or miles) from a given location.

```

GET /sensors/near?lat={latitude}&lng={longitude}&radius={radius}

```

Furthermore, with the right query language and indexing, we can implement more complex queries such as querying sensors by multiple tags within a bounding box or radius. We could also extend the model for a sensor, for example, to include a value and timestamp for when the sensor was last updated. This would open up the possibility to query for sensors by tags, time span, and aggregate values over space and time while querying within a bounding box or radius.

## Future Development

With more time, I would like to implement the following:

- Tighter delineation of read/write mutex locks using a RWMutex instead of a Mutex struct. This should be a very quick chagne.
- Less redundant handler testing, and a more complete testing suite in general. The core functionality is tested to 100% code coverage, but queries delineated by tagging and location are additional functionalities and not tested.
- More complete logging.
