package measurementdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateMeasurement(t *testing.T) Measurement {
	t.Logf("TestCreateMeasurement")
	arg := Measurement{
		Name:        "pump_temperature",
		Description: "Pump temperature",
	}
	t.Logf("Measurement: %v", arg)
	insertedResult, err := testStore.CreateMeasurement(context.TODO(), arg)
	require.NoError(t, err)
	arg.ID = insertedResult.ID
	arg.Status = "active"
	return arg
}

func TestCreateMeasurement(t *testing.T) {
	CreateMeasurement(t)
}

func TestGetMeasurement(t *testing.T) {
	t.Logf("TestGetMeasurement")
	arg := CreateMeasurement(t)
	t.Logf("Measurement: %v", arg)
	mea, err := testStore.GetMeasurement(context.TODO(), MeasurementID{ID: arg.ID})
	t.Logf("Inserted Measurement: %v", mea)
	require.NoError(t, err)
	require.Equal(t, arg.ID, mea.ID)
	require.Equal(t, arg.Name, mea.Name)
	require.Equal(t, arg.Description, mea.Description)
}

func TestListMeasurements(t *testing.T) {
	t.Logf("TestListMeasurements")
	arg := CreateMeasurement(t)
	t.Logf("Measurement: %v", arg)
	measurements, err := testStore.ListMeasurements(context.TODO(), 10, 0)
	t.Logf("List Measurements: %v", measurements)
	require.NoError(t, err)
	require.Equal(t, arg.ID, measurements[0].ID)
	require.Equal(t, arg.Name, measurements[0].Name)
	require.Equal(t, arg.Description, measurements[0].Description)
}

func TestArchiveMeasurement(t *testing.T) {
	t.Logf("TestArchiveMeasurement")
	arg := CreateMeasurement(t)
	t.Logf("Measurement: %v", arg)
	err := testStore.ArchiveMeasurement(context.TODO(), arg.ID.Hex())
	t.Logf("Inserted Measurement: %v", arg)
	require.NoError(t, err)
	mea, err := testStore.GetMeasurement(context.TODO(), MeasurementID{ID: arg.ID})
	require.NoError(t, err)
	require.Equal(t, "archived", mea.Status)
}
