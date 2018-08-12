package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"sync"
	"time"
	"fmt"
)

type InfluxConsumer struct {
	Influx InfluxSettings
	client client.Client
	lock   sync.Mutex
	output chan *client.Point
	batch  client.BatchPoints
}

func (consumer *InfluxConsumer) Connect() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     consumer.Influx.Url,
		Username: consumer.Influx.Username,
		Password: consumer.Influx.Password,
	})
	if err != nil {
		panic(err)
	}

	consumer.client = c
	consumer.output = make(chan *client.Point)
	consumer.Influx.batchConfig = client.BatchPointsConfig{
		Database: consumer.Influx.Database,
	}

	// Flush every 10 seconds
	go func() {
		for {
			time.Sleep(10 * time.Second)
			consumer.flush()
		}
	}()

	// Consume queue
	go func() {
		for point := range consumer.output {
			Debug.Printf("New Point: %v", point)
			consumer.lock.Lock()
			if consumer.batch == nil {
				b, err := client.NewBatchPoints(consumer.Influx.batchConfig)
				if err != nil {
					panic(err)
				}
				consumer.batch = b;
			}
			consumer.batch.AddPoint(point)
			consumer.lock.Unlock()
		}
	}()
}

func (consumer *InfluxConsumer) flush() {
	consumer.lock.Lock()
	if consumer.batch != nil {
		Debug.Print("Flushing...")
		err := consumer.client.Write(consumer.batch)
		if err != nil {
			Warn.Printf("Failed to write to InfluxDB: %v", err)
		}
		consumer.batch = nil
	}
	consumer.lock.Unlock()
}

func (consumer *InfluxConsumer) Close() {
	if consumer.client != nil {
		consumer.client.Close()
	}
}

func (consumer *InfluxConsumer) Write(measurement Measurement) {
	Debug.Print("Writing...")
	now := time.Now()

	// Global Info
	point, err := client.NewPoint(
		"voltcraft_charger",
		map[string]string{
			"device": measurement.USBPath,
			"mode":   measurement.Mode.name(),
			"state":  measurement.State.name(),
		},
		map[string]interface{}{
			"cells":        measurement.CellInfo.Count,
			"current":      measurement.Current.Milliamps,
			"voltage":      measurement.Current.Millivolts,
			"totalCharged": measurement.MilliampHoursCharged,
		},
		now,
	)
	if err != nil {
		panic(err)
	}

	consumer.output <- point

	// Cell Info
	for index, cell := range measurement.CellInfo.Cells {

		point, err := client.NewPoint(
			"voltcraft_cell",
			map[string]string{
				"device": measurement.USBPath,
				"mode":   measurement.Mode.name(),
				"state":  measurement.State.name(),
				"cell":   fmt.Sprintf("%d", index),
			},
			map[string]interface{}{
				"voltage": cell.Millivolts,
			},
			now,
		)
		if err != nil {
			panic(err)
		}

		consumer.output <- point
	}
}

type InfluxSettings struct {
	Url         string
	Username    string
	Password    string
	Database    string
	batchConfig client.BatchPointsConfig
}
