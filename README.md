# Sensor API

## Overview

The sensor JSON REST API exposes CRUD, geospatial, and tagged querying functionality to an inmemory sensor store.

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

#### Create

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "sensor-name", "location": {"latitude": 0.0, "longitude": 0.0}, "tags": ["tag1", "tag2"]}' http://localhost:8080/sensors
```

#### Read

```bash
curl -X GET http://localhost:8080/sensors/sensor-name
```

#### Update

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name": "sensor-name", "location": {"latitude": 0.0, "longitude": 0.0}, "tags": ["tag1", "tag2"]}' http://localhost:8080/sensors/sensor-name
```

#### Delete

```bash
curl -X DELETE http://localhost:8080/sensors/sensor-name
```

#### Nearest Sensor

```bash
curl -X GET http://localhost:8080/nearest-sensor?latitude=0.0&longitude=0.0
```

### Additional Endpoints:

Here are some additional endpoints I would implement for querying sensor data:

1. Get all sensors: Retrieve a list of all sensors with their metadata.

```

GET /sensors

```

2. Get sensors by tag: Retrieve a list of sensors that have a specific tag.

```

GET /sensors/tag/{tag}

```

3. Search sensors by name: Retrieve a list of sensors that match a given name pattern (e.g., using substring or regular expression matching).

```

GET /sensors/search?name={name_pattern}

```

4. Get sensors within a bounding box: Retrieve a list of sensors located within a specified bounding box, defined by minimum and maximum latitude and longitude values.

```

GET /sensors/within?min_lat={min_lat}&max_lat={max_lat}&min_lng={min_lng}&max_lng={max_lng}

```

5. Get sensors within a certain radius: Retrieve a list of sensors located within a specified radius (in meters or miles) from a given location.

```

GET /sensors/near?lat={latitude}&lng={longitude}&radius={radius}

```

6. Get the total number of sensors: Retrieve the count of all sensors in the store.

```

GET /sensors/count

```

7. Get the list of unique tags: Retrieve a list of all unique tags associated with the sensors.

```

GET /sensors/tags

```

Furthermore, with the right query language and indexing, we can implement more complex queries such as querying sensors by multiple tags within a bounding box or radius. We could also extend the model for a sensor, for example, to include a value and timestamp for when the sensor was last updated. This would open up the possibility to query for sensors by tags, time span, and aggregate values over space and time while querying within a bounding box or radius.
