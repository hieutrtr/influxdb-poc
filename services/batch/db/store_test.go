package batchdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateRandomBatch(t *testing.T) Batch {
	arg := Batch{
		Name:          "test",
		MeasurementID: primitive.NewObjectID(),
		StartTime:     time.Now().Unix(),
		EndTime:       time.Now().Unix() + int64(time.Hour.Seconds())*24*7,
		Interval:      int64(time.Hour.Seconds()) / 2,
	}
	res, err := testStore.CreateBatch(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res.ID)
	arg.ID = res.ID
	return arg
}

func TestCreateBatch(t *testing.T) {
	CreateRandomBatch(t)
}

func TestGetBatch(t *testing.T) {
	arg := CreateRandomBatch(t)
	res, err := testStore.GetBatch(context.Background(), BatchID{ID: arg.ID})
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, arg.ID, res.ID)
	require.Equal(t, arg.Name, res.Name)
	require.Equal(t, arg.MeasurementID, res.MeasurementID)
	require.Equal(t, arg.StartTime, res.StartTime)
	require.Equal(t, arg.EndTime, res.EndTime)
	require.Equal(t, arg.Interval, res.Interval)
	require.Equal(t, arg.FilePath, res.FilePath)
	require.Equal(t, "active", res.Status)
}

func TestListBatchesByStatus(t *testing.T) {
	arg := CreateRandomBatch(t)
	res, err := testStore.ListBatchesByStatus(context.Background(), "active", 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, arg.ID, res[0].ID)
	require.Equal(t, arg.Name, res[0].Name)
	require.Equal(t, arg.MeasurementID, res[0].MeasurementID)
	require.Equal(t, arg.StartTime, res[0].StartTime)
	require.Equal(t, arg.EndTime, res[0].EndTime)
	require.Equal(t, arg.Interval, res[0].Interval)
	require.Equal(t, arg.FilePath, res[0].FilePath)
	require.Equal(t, "active", res[0].Status)
}

func TestChangeBatchTime(t *testing.T) {
	arg := CreateRandomBatch(t)
	err := testStore.ChangeBatchTime(context.Background(), BatchID{ID: arg.ID}, BatchesTimeRange{
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + int64(time.Hour.Seconds())*24*7,
	})
	require.NoError(t, err)
}

func TestChangeBatchInterval(t *testing.T) {
	arg := CreateRandomBatch(t)
	err := testStore.ChangeBatchInterval(context.Background(), BatchID{ID: arg.ID}, int64(time.Hour.Seconds())/2)
	require.NoError(t, err)
}

func TestStopBatch(t *testing.T) {
	arg := CreateRandomBatch(t)
	err := testStore.StopBatch(context.Background(), BatchID{ID: arg.ID})
	require.NoError(t, err)
}

func TestUpdateFileBatch(t *testing.T) {
	arg := CreateRandomBatch(t)
	err := testStore.UpdateFileBatch(context.Background(), BatchID{ID: arg.ID}, "path/file.csv")
	require.NoError(t, err)
}
