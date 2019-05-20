package database

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	// MyDB specifies name of database
	MyDB = "go_influx"
)

type Point struct {
	Key   string
	Value float32
}

// Insert saves points to database
func Insert(points []Point, t time.Time) error {
	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "username",
		Password: "password",
	})
	if err != nil {
		return err
	}
	defer httpClient.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "temperatures",
		Precision: "us",
	})
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})
	for _, t := range points {
		fields[t.Key] = t.Value
	}

	point, err := client.NewPoint(
		"temperatures",
		nil,
		fields,
		t,
	)

	if err != nil {
		return err
	}
	bp.AddPoint(point)

	return httpClient.Write(bp)
}
