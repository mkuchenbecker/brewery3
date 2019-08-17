package logger

import (
	"context"

	"github.com/mkuchenbecker/brewery3/brewery/utils"
)

const (
	// MyDB specifies name of database
	MyDB = "go_influx"
)

type Logger interface {
	InsertTemperature(ctx context.Context, tps *utils.TemperaturePointSink) error
}

// type defaultLogger struct {
// }

func NewDefault() Logger {
	return &fakeLogger{}
}

func NewFake() Logger {
	return &fakeLogger{}
}

type fakeLogger struct {
}

func (logger fakeLogger) InsertTemperature(ctx context.Context, tps *utils.TemperaturePointSink) (err error) {
	utils.Printf("%v", tps.Temps)
	return nil
}

// // Insert saves points to database
// func (logger defaultLogger) InsertTemperature(ctx context.Context, tps *utils.TemperaturePointSink) (err error) {
// 	utils.Printf("%v", tps.Temps)
// 	c, err := client.NewHTTPClient(client.HTTPConfig{
// 		Addr: "http://localhost:8086",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		closeErr := c.Close()
// 		if err == nil {
// 			err = closeErr
// 		}
// 	}()

// 	// Create a new point batch
// 	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
// 		Database:  MyDB,
// 		Precision: "s",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	// Create a point and add to batch

// 	tags := map[string]string{}
// 	fields := tps.ToInterface()

// 	pt, err := client.NewPoint("temperature", tags, fields, time.Now())
// 	if err != nil {
// 		return err
// 	}
// 	bp.AddPoint(pt)

// 	// Write the batch
// 	if err := c.Write(bp); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // queryDB convenience function to query the database
// func (logger defaultLogger) queryDB(cmd string) (res []client.Result, err error) {
// 	q := client.Query{
// 		Command:  cmd,
// 		Database: MyDB,
// 	}
// 	c, err := client.NewHTTPClient(client.HTTPConfig{
// 		Addr: "http://localhost:8086",
// 	})
// 	if err != nil {
// 		return []client.Result{}, nil
// 	}
// 	if response, err := c.Query(q); err == nil {
// 		if response.Error() != nil {
// 			return res, response.Error()
// 		}
// 		res = response.Results
// 	} else {
// 		return res, err
// 	}
// 	return res, nil
// }
