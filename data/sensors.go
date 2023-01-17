package data

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type SensorRecord struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type Sensor struct {
	Name    string         `json:"name"`
	Records []SensorRecord `json:"records"`
}

type Sensors map[string]*Sensor

func LoadAllSensors() (Sensors, error) {
	// Load csv file
	filePath := "../resources/sensor.csv"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Parse csv file
	parser := csv.NewReader(file)
	records, err := parser.ReadAll()
	if err != nil {
		return nil, err
	}
	// Convert records to Sensor
	var sensorsMap Sensors
	firstColumn := records[0]
	for _, record := range records[1:20] {
		for i := 2; i < len(record); i++ {
			name := firstColumn[i]
			if sensorsMap == nil {
				sensorsMap = make(Sensors)
			}
			if sensorsMap[name] == nil {
				sensorsMap[name] = &Sensor{
					Name: name,
				}
			}
			sensorsMap[name].Records = append(sensorsMap[name].Records, SensorRecord{
				Value:     stringToFloat(record[i]),
				Timestamp: dateToTimestamp(record[1]),
			})
		}
	}
	return sensorsMap, nil
}

func stringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func dateToTimestamp(s string) int64 {
	date, _ := time.Parse("2006-01-02 15:04:05", s)
	return date.Unix()
}
