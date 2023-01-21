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

// NewMeasurementStore returns a measurement store. If collection is nil a default collection will be used. This is primarily for testing and should not be used by third party code.
//
// @param collection - the collection to operate on. It must exist or be empty.
//
// @return MeasurementStore for the given collection. The store will be ready to accept measurements from it and store them
func NewMeasurementStore(collection *mongo.Collection) *MeasurementStore {
	return &MeasurementStore{
		collection: collection,
	}
}

// CreateMeasurement creates a new measurement in the MeasurementStore. It returns the measurement ID which can be used to refer to it later in GetMeasurement.
//
// @param m
// @param ctx - An instance of context. Context for authentication logging cancellation deadlines tracing etc.
// @param measurement - An instance of Measurement to be inserted. Must not be nil.
//
// @return MeasurementID or an error if something goes wrong with the insertion process. If the measurement already exists a measurement ID is returned
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

// GetMeasurement retrieves a measurement from the database. If the measurement is deleted an error is returned. If the measurement does not exist a 404 error is returned
//
// @param m
// @param ctx - An instance of a context. Context
// @param meaID - The ID of the measurement to retrieve.
//
// @return A Measurement or an error that satisfies the DBError interface. The error will be of type MeasurementNotFound if the measurement has been deleted
func (m *MeasurementStore) GetMeasurement(ctx context.Context, meaID MeasurementID) (*Measurement, error) {
	var measurement Measurement
	err := m.collection.FindOne(ctx, MeasurementID{ID: meaID.ID}).Decode(&measurement)
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	if measurement.Status == "deleted" {
		return nil, NewErrMeasurementNotFound("measurement not found", 404)
	}
	return &measurement, nil
}

// ListMeasurements returns a list of measurements. The measurements are sorted by _id in descending order.
//
// @param m
// @param ctx - Context to pass into request. Background. Deprecated : use MeasurementStore. List () instead.
// @param limit - Maximum number of measurements to return. Limit is a cap on the number of results returned.
// @param skip - Skip is a cap on the number of results returned. Skip is a cap on the number of results returned.
//
// @return An array of measurements or an error if something went wrong during the operation such as database errors or there was a problem getting data
func (m *MeasurementStore) ListMeasurements(ctx context.Context, limit int64, skip int64) ([]Measurement, error) {
	var measurements []Measurement
	cursor, err := m.collection.Find(ctx, bson.M{"status": "active"}, options.Find().SetSort(map[string]int{"_id": -1}).SetLimit(limit).SetSkip(skip))
	if err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	if err = cursor.All(ctx, &measurements); err != nil {
		return nil, NewDBError(err.Error(), 500)
	}
	return measurements, nil
}

// ArchiveMeasurement archives a measurement. The status will be set to archived and updated_at will be set to the current time
//
// @param m
// @param ctx - context to pass into request
// @param id - hex ID of the measurement
//
// @return DBError 400 for bad id 500 for other errors or nil if everything went fine in this case no
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
