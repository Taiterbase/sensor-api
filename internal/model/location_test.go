package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationEquals(t *testing.T) {
	loc1 := Location{Latitude: 12.34, Longitude: -56.78}
	loc2 := Location{Latitude: 12.34, Longitude: -56.78}
	loc3 := Location{Latitude: 43.21, Longitude: -87.65}

	assert.True(t, loc1.Equals(loc2), "Expected loc1 and loc2 to be equal")
	assert.False(t, loc1.Equals(loc3), "Expected loc1 and loc3 to not be equal")
}

func TestLocationIsValid(t *testing.T) {
	validLoc := Location{Latitude: 12.34, Longitude: -56.78}
	invalidLoc1 := Location{Latitude: 91, Longitude: -56.78}
	invalidLoc2 := Location{Latitude: 12.34, Longitude: 181}
	invalidLoc3 := Location{Latitude: -91, Longitude: -56.78}
	invalidLoc4 := Location{Latitude: 12.34, Longitude: -181}
	emptyLoc := Location{}

	assert.True(t, validLoc.IsValid(), "Expected validLoc to be valid")
	assert.False(t, invalidLoc1.IsValid(), "Expected invalidLoc1 to be invalid")
	assert.False(t, invalidLoc2.IsValid(), "Expected invalidLoc2 to be invalid")
	assert.False(t, invalidLoc3.IsValid(), "Expected invalidLoc3 to be invalid")
	assert.False(t, invalidLoc4.IsValid(), "Expected invalidLoc4 to be invalid")
	assert.False(t, emptyLoc.IsValid(), "Expected emptyLoc to be invalid")
}
