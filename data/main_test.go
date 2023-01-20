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
	meaChan := make(chan Measurement)
	defer close(meaChan)
	go func(tt *testing.T, b bool, i ...interface{}) {
		for mea := range meaChan {
			tt.Logf("Measurement: %v", mea)
			require.NotEmpty(tt, mea)
		}
	}(t, true)
	err := LoadAllSensors("../resources/sensor.csv", meaChan, 100)
	if err != nil {
		t.Errorf("Error loading sensors: %v", err)
	}
	require.Equal(t, 1, 2)
}
