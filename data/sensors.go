package data

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type Sensor struct {
	Name         string        `json:"name"`
	Measurements []Measurement `json:"measurement"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Field struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}
type Measurement struct {
	Name      string  `json:"name"`
	Fields    []Field `json:"fields"`
	Timestamp int64   `json:"timestamp"`
	Tags      []Tag   `json:"tags"`
}

type Sensors map[string]*Sensor

func LoadAllSensors(filePath string, meaChan chan<- Measurement, limit int) error {
	// Load csv file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Parse csv file
	parser := csv.NewReader(file)
	records, err := parser.ReadAll()
	if err != nil {
		return err
	}
	// Convert records to Sensor
	var sensorsMap Sensors
	firstColumn := records[0]
	if limit == 0 || limit > len(records) {
		limit = len(records)
	}
	for _, record := range records[1:limit] {
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
			meaChan <- Measurement{
				Name: name,
				Fields: []Field{
					{
						Key:   "value",
						Value: stringToFloat(record[i]),
					},
				},
				Tags: []Tag{
					{
						Key:   "sensor",
						Value: name,
					},
				},
				Timestamp: dateToTimestamp(record[1]),
			}
		}
	}
	return nil
}

func stringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func dateToTimestamp(s string) int64 {
	date, _ := time.Parse("2006-01-02 15:04:05", s)
	return date.Unix()
}
