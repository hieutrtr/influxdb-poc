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
	fmt.Println("loaded config")
	sensors, err := data.LoadAllSensors("./resources/sensor.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("loaded sensors")
	influxClient := influxdb2.NewClient("http://localhost:8086", conf.InfluxDBToken)
	defer influxClient.Close()
	store := process.NewMeasurementStore("pump_temp_dev2", influxClient.WriteAPI(conf.InfluxDBOrg, conf.InfluxDBBucket))
	for _, sensor := range sensors {
		for _, measurement := range sensor.Measurements {
			fmt.Println(measurement)
			measurement.Timestamp += 151000000
			store.Write(measurement)
		}
	}
	store.Flush()
}
