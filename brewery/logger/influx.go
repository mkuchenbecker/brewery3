package logger

import (
	"context"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	client "github.com/influxdata/influxdb1-client/v2"
)

const (
	// MyDB specifies name of database
	MyDB = "go_influx"
)

type Logger interface {
	InsertTemperature(ctx context.Context, points map[string]interface{}) error
}

type defaultLogger struct {
}

func NewDefault() {}

// Insert saves points to database
func (logger defaultLogger) InsertTemperature(ctx context.Context, tps *utils.TemperaturePointSink) (err error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return err
	}
	defer func() {
		closeErr := c.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	// Create a point and add to batch

	tags := map[string]string{}
	fields := tps.ToInterface()

	pt, err := client.NewPoint("temperature", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}

// queryDB convenience function to query the database
func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
