package process

import (
	"os"
	"testing"
	"time"

	"github.com/hieutrtr/influxdb-poc/data"
	"github.com/hieutrtr/influxdb-poc/utils"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func TestMain(m *testing.M) {
	// Run the tests
	os.Exit(m.Run())
}

func TestWrite(t *testing.T) {
	conf, err := utils.LoadConfig("..")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}
	t.Logf("TestWrite")
	mea := data.Measurement{
		Name: "pump_temperature",
		Tags: []data.Tag{
			{
				Key:   "sensor",
				Value: "pump1",
			},
		},
		Fields: []data.Field{
			{
				Key:   "value",
				Value: 1.0,
			},
		},
		Timestamp: time.Now().Unix(),
	}
	t.Logf("Measurement: %v", mea)
	client := influxdb2.NewClient("http://localhost:8086", conf.InfluxDBToken)
	store := NewMeasurementStore("pump_temperature", client.WriteAPI(conf.InfluxDBOrg, conf.InfluxDBBucket))
	store.Write(mea)
	store.Flush()
	client.Close()
}
