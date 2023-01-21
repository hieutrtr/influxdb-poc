package batchdb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
}

type BatchStore struct {
	collection *mongo.Collection
}

func NewBatchStore(collection *mongo.Collection) *BatchStore {
	return &BatchStore{
		collection: collection,
	}
}

type Batch struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name,omitempty"`
	MeasurementID primitive.ObjectID `bson:"measurement_id"`
	StartTime     int64              `bson:"start_time"`
	EndTime       int64              `bson:"end_time"`
	Interval      int64              `bson:"interval"`
	FilePath      string             `bson:"file_path"`
	Status        string             `bson:"status,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty"`
}

type BatchID struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (m *BatchStore) CreateBatch(ctx context.Context, batch Batch) (*BatchID, error) {
	batch.ID = primitive.NewObjectID()
	batch.CreatedAt = time.Now()
	batch.UpdatedAt = time.Now()
	batch.Status = "active"
	_, err := m.collection.InsertOne(ctx, batch)
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return &BatchID{ID: batch.ID}, nil
}

func (m *BatchStore) GetBatch(ctx context.Context, batchID BatchID) (*Batch, error) {
	var batch Batch
	err := m.collection.FindOne(ctx, batchID).Decode(&batch)
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return &batch, nil
}

type Batches []Batch

func (m *BatchStore) ListBatchesByStatus(ctx context.Context, status string, limit int64, offset int64) (Batches, error) {
	var batches Batches
	cur, err := m.collection.Find(ctx, bson.M{"status": status}, options.Find().SetSort(map[string]int{"_id": -1}).SetLimit(limit).SetSkip(offset))
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	if err := cur.All(ctx, &batches); err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return batches, nil
}

type BatchesTimeRange struct {
	StartTime int64 `bson:"start_time,omitempty"`
	EndTime   int64 `bson:"end_time,omitempty"`
}

func (m *BatchStore) ChangeBatchTime(ctx context.Context, batchID BatchID, timeRange BatchesTimeRange) error {
	_, err := m.collection.UpdateOne(ctx, batchID, bson.M{"$set": timeRange})
	if err != nil {
		return NewDBError(err.Error(), 500)
	}
	return nil
}

func (m *BatchStore) ChangeBatchInterval(ctx context.Context, batchID BatchID, interval int64) error {
	_, err := m.collection.UpdateOne(ctx, batchID, bson.M{"$set": bson.M{"interval": interval}})
	if err != nil {
		return NewDBError(err.Error(), 500)
	}
	return nil
}

func (m *BatchStore) StopBatch(ctx context.Context, batchID BatchID) error {
	_, err := m.collection.UpdateOne(ctx, batchID, bson.M{"$set": bson.M{"status": "inactive"}})
	if err != nil {
		return NewDBError(err.Error(), 500)
	}
	return nil
}

func (m *BatchStore) UpdateFileBatch(ctx context.Context, batchID BatchID, filePath string) error {
	_, err := m.collection.UpdateOne(ctx, batchID, bson.M{"$set": bson.M{"file_path": filePath}})
	if err != nil {
		return NewDBError(err.Error(), 500)
	}
	return nil
}
