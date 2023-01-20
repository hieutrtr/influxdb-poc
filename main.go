package main

import (
	"fmt"

	"github.com/hieutrtr/influxdb-poc/data"
	"github.com/hieutrtr/influxdb-poc/process"
	"github.com/hieutrtr/influxdb-poc/utils"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	conf, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// influxClient := influxdb2.NewClient("http://localhost:8086", conf.InfluxDBToken)
	influxClient := influxdb2.NewClientWithOptions("http://localhost:8086", conf.InfluxDBToken, influxdb2.DefaultOptions().SetFlushInterval(3000))

	store := process.NewMeasurementStore("pump_temp_dev4", influxClient.WriteAPI(conf.InfluxDBOrg, conf.InfluxDBBucket))

	fmt.Println("loaded config")
	meaChan := make(chan data.Measurement)
	defer close(meaChan)

	go func() {
		for measurement := range meaChan {
			fmt.Println(measurement)
			measurement.Timestamp += 151000000
			store.Write(measurement)
		}
		store.Flush()
		influxClient.Close()
	}()

	if err := data.LoadAllSensors("./resources/sensor.csv", meaChan, 0); err != nil {
		panic(err)
	}
	fmt.Println("loaded sensors")
}
