package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"transfer-plugins/internal/models"
	"transfer-plugins/pkg/logger"
)

type Service interface {
	WriteData(points []models.UserAction) error
	WritePoint(point models.UserAction) error
	Flush()
}

type service struct {
	writeApi *api.WriteAPI
}

func NewService(writeApi *api.WriteAPI) Service {
	return &service{
		writeApi: writeApi,
	}
}

func (s *service) WriteData(points []models.UserAction) error {
	// Get errors channel
	errorsCh := (*s.writeApi).Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			logger.Errorf("write error: %s\n", err.Error())
		}
	}()
	for _, value := range points {
		point := influxdb2.NewPointWithMeasurement(value.Metrics).
			AddTag("userId", value.UserId).
			AddTag("type", value.Type).
			AddField("subType", value.SubType).
			AddField("targetId", value.TargetId).
			SetTime(value.Timestamp)
		(*s.writeApi).WritePoint(point)
	}
	(*s.writeApi).Flush()
	return nil
}

func (s *service) WritePoint(point models.UserAction) error {
	p := influxdb2.NewPointWithMeasurement(point.Metrics).
		AddTag("userId", point.UserId).
		AddTag("type", point.Type).
		AddField("subType", point.SubType).
		AddField("targetId", point.TargetId).
		SetTime(point.Timestamp)
	(*s.writeApi).WritePoint(p)
	return nil
}

func (s *service) Flush() {
	//force flush buffer
	//默认时间间隔 flushInterval: 1000（ms）
	(*s.writeApi).Flush()
}
