package store

import (
	"fmt"
	"sensor-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	store := NewInMemorySensorStore()

	sensor1 := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}

	// Test AddSensor and GetSensor
	code, err := store.AddSensor(sensor1)
	assert.NoError(t, err)
	assert.Equal(t, 201, code)

	retrievedSensor, _, err := store.GetSensor("Sensor1")
	assert.NoError(t, err)
	assert.Equal(t, sensor1, retrievedSensor)

	// Test duplicate sensor addition
	code, err = store.AddSensor(sensor1)
	assert.Error(t, err)
	assert.Equal(t, 400, code)
}

func TestGet(t *testing.T) {
	store := NewInMemorySensorStore()

	// Test getting a sensor that doesn't exist
	_, _, err := store.GetSensor("Sensor1")
	assert.Error(t, err)

	sensor1 := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}

	_, err = store.AddSensor(sensor1)
	assert.NoError(t, err)

	retrievedSensor, _, err := store.GetSensor("Sensor1")
	assert.NoError(t, err)
	assert.Equal(t, sensor1, retrievedSensor)
}

func TestUpdate(t *testing.T) {
	store := NewInMemorySensorStore()

	sensor := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}

	_, err := store.AddSensor(sensor)
	assert.NoError(t, err)

	sensor1Updated := sensor
	sensor1Updated.Location.Latitude = 37.7833
	_, err = store.UpdateSensor("Sensor1", &sensor1Updated)
	assert.NoError(t, err)

	updatedSensor, _, err := store.GetSensor("Sensor1")
	assert.NoError(t, err)
	assert.Equal(t, sensor1Updated, updatedSensor)
	assert.Equal(t, 37.7833, updatedSensor.Location.Latitude)

	assert.Equal(t, 1, len(store.sensors))
	assert.Equal(t, 1, store.rt.Len())
	min, _ := store.rt.Bounds()
	assert.Equal(t, 37.7833, min[0])
	assert.Equal(t, -122.4194, min[1])

	// Test updating a sensor that doesn't exist
	_, err = store.UpdateSensor("Sensor2", &sensor1Updated)
	assert.Error(t, err)

	// Test updating a sensor with a nil pointer
	_, err = store.UpdateSensor("Sensor1", nil)
	assert.Error(t, err)

	// Test updating a sensor with a nil name
	sensor1Updated.Location = model.Location{
		Latitude:  37.7833,
		Longitude: -122.4194,
	}
	sensor1Updated.Name = ""
	_, err = store.UpdateSensor("Sensor1", &sensor1Updated)
	assert.Error(t, err)
}

func TestRemove(t *testing.T) {
	store := NewInMemorySensorStore()

	sensor1 := model.Sensor{
		Name: "Sensor1",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}

	sensor2 := model.Sensor{
		Name: "Sensor2",
		Location: model.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}

	_, err := store.AddSensor(sensor1)
	assert.NoError(t, err)
	_, err = store.AddSensor(sensor2)
	assert.NoError(t, err)

	_, err = store.RemoveSensor("Sensor1")
	assert.NoError(t, err)

	_, _, err = store.GetSensor("Sensor1")
	assert.Error(t, err)

	// Test removing a sensor that doesn't exist
	_, err = store.RemoveSensor("Sensor1")
	assert.Error(t, err)

	assert.Equal(t, 1, len(store.sensors))
	assert.Equal(t, 1, store.rt.Len())

	_, err = store.RemoveSensor("Sensor2")
	assert.NoError(t, err)

	assert.Equal(t, 0, len(store.sensors))
	assert.Equal(t, 0, store.rt.Len())
}

func TestNearby(t *testing.T) {
	t.Run("Three close sensors", func(t *testing.T) {
		store := NewInMemorySensorStore()
		sensor1 := model.Sensor{
			Name: "Sensor1",
			Location: model.Location{
				Latitude:  39.0921,
				Longitude: -123.5222,
			},
		}
		sensor2 := model.Sensor{
			Name: "Sensor2",
			Location: model.Location{
				Latitude:  37.7833,
				Longitude: -122.4167,
			},
		}
		sensor3 := model.Sensor{
			Name: "Sensor3",
			Location: model.Location{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		_, err := store.AddSensor(sensor1)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor2)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor3)
		assert.NoError(t, err)

		// Test getting the nearest sensor
		nearestSensor, _, err := store.GetNearestSensor(model.Location{
			Latitude:  37.775,
			Longitude: -122.42,
		})
		assert.NoError(t, err)
		assert.Equal(t, sensor3, *nearestSensor)
	})

	t.Run("Three far sensors", func(t *testing.T) {
		store := NewInMemorySensorStore()
		sensor1 := model.Sensor{
			Name: "Sensor1",
			Location: model.Location{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}
		sensor2 := model.Sensor{
			Name: "Sensor2",
			Location: model.Location{
				Latitude:  37.7833,
				Longitude: -122.4167,
			},
		}
		sensor3 := model.Sensor{
			Name: "Sensor3",
			Location: model.Location{
				Latitude:  39.0921,
				Longitude: -123.5222,
			},
		}

		_, err := store.AddSensor(sensor1)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor2)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor3)
		assert.NoError(t, err)

		// Test getting the nearest sensor
		nearestSensor, _, err := store.GetNearestSensor(model.Location{
			Latitude:  39.0920,
			Longitude: -123.5221,
		})
		assert.NoError(t, err)
		assert.Equal(t, sensor3, *nearestSensor)
	})

	t.Run("Twenty evenly placed sensors", func(t *testing.T) {
		store := NewInMemorySensorStore()
		for i := 0; i < 20; i++ {
			sensor := model.Sensor{
				Name: fmt.Sprintf("Sensor%d", i),
				Location: model.Location{
					Latitude:  39,
					Longitude: -110 + float64(i)*.9*-1,
				},
			}
			_, err := store.AddSensor(sensor)
			assert.NoError(t, err)
		}

		// Test getting the nearest sensor
		nearestSensor, _, err := store.GetNearestSensor(model.Location{
			Latitude:  39.0920,
			Longitude: -123.5221,
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Sensor{Name: "Sensor15", Location: model.Location{
			Latitude:  39,
			Longitude: -110 + float64(15)*.9*-1,
		}}, *nearestSensor)
	})

	t.Run("Test bad input", func(t *testing.T) {
		store := NewInMemorySensorStore()
		sensor1 := model.Sensor{
			Name: "Sensor1",
			Location: model.Location{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}
		sensor2 := model.Sensor{
			Name: "Sensor2",
			Location: model.Location{
				Latitude:  37.7833,
				Longitude: -122.4167,
			},
		}
		sensor3 := model.Sensor{
			Name: "Sensor3",
			Location: model.Location{
				Latitude:  39.0921,
				Longitude: -123.5222,
			},
		}

		_, err := store.AddSensor(sensor1)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor2)
		assert.NoError(t, err)
		_, err = store.AddSensor(sensor3)
		assert.NoError(t, err)

		// Test getting the nearest sensor with a bad location
		_, _, err = store.GetNearestSensor(model.Location{
			Latitude:  91,
			Longitude: -123.5221,
		})
		assert.Error(t, err)

		// Test getting the nearest sensor with a bad location
		_, _, err = store.GetNearestSensor(model.Location{
			Latitude:  70,
			Longitude: -193.5221,
		})
		assert.Error(t, err)

		// Test getting the nearest sensor with a nil location
		_, _, err = store.GetNearestSensor(model.Location{})
		assert.Error(t, err)

		// Test getting the nearest sensor with no sensors
		store = NewInMemorySensorStore()
		_, _, err = store.GetNearestSensor(model.Location{
			Latitude:  39.0920,
			Longitude: -123.5221,
		})

		// Test getting the nearest sensor with one sensor
		assert.Error(t, err)
		store.AddSensor(sensor1)
		_, _, err = store.GetNearestSensor(model.Location{
			Latitude:  39.0920,
			Longitude: -123.5221,
		})
		assert.NoError(t, err)
	})
}
