package influxdb

import (
	"errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"transfer-plugins/configs"
)

var influxDBClient *dbClient

type dbClient struct {
	influxDbC *influxdb2.Client
	writeApi  *api.WriteAPI
}

func NewClient() error {
	config := configs.Get().InfluxDB
	client := influxdb2.NewClient(config.Addr, config.Token)
	if client == nil {
		return errors.New("Get client nil, please check addr/token")
	}
	writeApi := client.WriteAPI(config.Org, config.Bucket)
	if writeApi == nil {
		return errors.New("Get api nil, please check org/bucket")
	}

	influxDBClient = &dbClient{
		influxDbC: &client,
		writeApi:  &writeApi,
	}
	return nil
}

func WriteApi() *api.WriteAPI {
	return influxDBClient.writeApi
}

func Close() {
	(*influxDBClient.influxDbC).Close()
	log.Printf("influxdb2 client stopped\n")
}
