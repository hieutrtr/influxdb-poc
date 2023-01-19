package process

import (
	"time"

	"github.com/hieutrtr/influxdb-poc/data"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type MeasurementStore struct {
	Name     string
	WriteAPI api.WriteAPI
}

func NewMeasurementStore(name string, api api.WriteAPI) *MeasurementStore {
	return &MeasurementStore{
		WriteAPI: api,
		Name:     name,
	}
}

func (m *MeasurementStore) Write(measurement data.Measurement) {
	tags := make(map[string]string)
	for _, tag := range measurement.Tags {
		tags[tag.Key] = tag.Value
	}
	fields := make(map[string]interface{})
	for _, field := range measurement.Fields {
		fields[field.Key] = field.Value
	}
	point := influxdb2.NewPoint(
		m.Name,
		tags,
		fields,
		time.Unix(measurement.Timestamp, 0),
	)
	m.WriteAPI.WritePoint(point)
}

func (m *MeasurementStore) Flush() {
	m.WriteAPI.Flush()
}
