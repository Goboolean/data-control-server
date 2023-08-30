package broker

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/api/model"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/kafka"
	"github.com/Goboolean/fetch-server.v1/internal/util/prometheus"
)

type Adapter struct {
	conf *kafka.Configurator
	pub  *kafka.Publisher
}

func NewAdapter(conf *kafka.Configurator, pub *kafka.Publisher) out.TransmissionPort {
	return &Adapter{
		conf: conf,
		pub:  pub,
	}
}

func (a *Adapter) TransmitStockBatch(ctx context.Context, stock string, batch []*vo.StockAggregate) error {

	prometheus.BrokerCounter.Add(float64(len(batch)))

	converted := make([]*model.StockAggregate, len(batch))

	for idx := range converted {
		converted[idx] = &model.StockAggregate{
			EventType: batch[idx].EventType,
			Volume:    batch[idx].Volume,
			Min:       batch[idx].Min,
			Max:       batch[idx].Max,
			Open:      batch[idx].Open,
			Closed:    batch[idx].Closed,
			StartTime: batch[idx].Time,
		}
	}

	return a.pub.SendDataBatch(stock, converted)
}

func (a *Adapter) CreateStockBroker(ctx context.Context, stock string) error {
	return a.conf.CreateTopic(ctx, stock)
}
