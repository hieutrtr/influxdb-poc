package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Run the tests
	os.Exit(m.Run())
}

func TestLoadAllSensors(t *testing.T) {
	sensors, err := LoadAllSensors("../resources/sensor.csv")
	if err != nil {
		t.Errorf("Error loading sensors: %v", err)
	}
	for _, sensor := range sensors {
		t.Logf("Sensor: %v", sensor.Measurements[0])
		require.NotEmpty(t, sensor.Measurements[0])
	}
	require.Equal(t, 1, 1)
}
