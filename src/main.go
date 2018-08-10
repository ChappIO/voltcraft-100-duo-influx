package main

import (
	"log"
	"os"
	"io/ioutil"
)

var Warn = log.New(os.Stdout, "[WARN ] ", log.LstdFlags)
var Info = log.New(os.Stdout, "[INFO ] ", log.LstdFlags)
var Debug = log.New(ioutil.Discard, "[DEBUG] ", log.LstdFlags)

func main() {
	Info.Print("Hello!")
	if _, ok := os.LookupEnv("DEBUG"); ok {
		Debug.SetOutput(os.Stdout)
	}
	Debug.Print("Debug logging enabled!")

	volt := NewVoltcraft()

	influx := InfluxConsumer{
		Influx: InfluxSettings{
			Url:      getEnvOr("INFLUXDB_URL", "http://localhost:8086"),
			Username: getEnvOr("INFLUXDB_WRITE_USER", ""),
			Password: getEnvOr("INFLUXDB_WRITE_USER_PASSWORD", ""),
			Database: getEnvOr("INFLUXDB_DB", "telegraf"),
		},
	}
	influx.Connect()
	defer influx.Close()

	volt.OnData = func(measurement Measurement) {
		influx.Write(measurement)
	}

	volt.StartScanning()
}

func getEnvOr(name string, def string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return def;
}
