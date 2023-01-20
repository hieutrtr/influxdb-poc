package measurementdb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/hieutrtr/influxdb-poc/services/measurements/db Store
type Store interface {
	CreateMeasurement(context.Context, Measurement) (*MeasurementID, error)
	GetMeasurement(context.Context, MeasurementID) (*Measurement, error)
	ListMeasurements(context.Context, int64, int64) ([]Measurement, error)
	ArchiveMeasurement(context.Context, string) error
}

type Measurement struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Status      string             `bson:"status,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}
type MeasurementStore struct {
	collection *mongo.Collection
}

func NewMeasurementStore(collection *mongo.Collection) *MeasurementStore {
	return &MeasurementStore{
		collection: collection,
	}
}

func (m *MeasurementStore) CreateMeasurement(ctx context.Context, measurement Measurement) (*MeasurementID, error) {
	measurement.ID = primitive.NewObjectID()
	measurement.CreatedAt = time.Now()
	measurement.UpdatedAt = time.Now()
	measurement.Status = "active"
	_, err := m.collection.InsertOne(ctx, measurement)
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return &MeasurementID{ID: measurement.ID}, nil
}

type MeasurementID struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (m *MeasurementStore) GetMeasurement(ctx context.Context, meaID MeasurementID) (*Measurement, error) {
	var measurement Measurement
	err := m.collection.FindOne(ctx, MeasurementID{ID: meaID.ID}).Decode(&measurement)
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	if measurement.Status == "deleted" {
		return nil, NewDBError("measurement not found", 404)
	}
	return &measurement, nil
}

func (m *MeasurementStore) ListMeasurements(ctx context.Context, limit int64, skip int64) ([]Measurement, error) {
	var measurements []Measurement
	cursor, err := m.collection.Find(ctx, Measurement{Status: "active"}, options.Find().SetSort(map[string]int{"_id": -1}).SetLimit(limit).SetSkip(skip))
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	if err = cursor.All(ctx, &measurements); err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return measurements, nil
}

func (m *MeasurementStore) ArchiveMeasurement(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewDBError(err.Error(), 400)
	}
	_, err = m.collection.UpdateOne(ctx, MeasurementID{ID: oid}, bson.M{"$set": bson.M{"status": "archived", "updated_at": time.Now()}})
	if err != nil {
		return NewDBError(err.Error(), 500)
	}
	return nil
}
