package broker

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/vo"
	"github.com/Goboolean/fetch-server/internal/infrastructure/broker"
	"github.com/Goboolean/fetch-server/internal/util/prometheus"
)

type Adapter struct {
	conf *broker.Configurator
	pub  *broker.Publisher
}

func NewAdapter(conf *broker.Configurator, pub *broker.Publisher) out.TransmissionPort {
	return &Adapter{
		conf: conf,
		pub:  pub,
	}
}

func (a *Adapter) TransmitStockBatch(ctx context.Context, stock string, batch []*vo.StockAggregate) error {

	prometheus.BrokerCounter.Add(float64(len(batch)))

	converted := make([]*broker.StockAggregate, len(batch))

	for idx := range converted {
		converted[idx] = &broker.StockAggregate{
			EventType: batch[idx].EventType,
			Average:   float32(batch[idx].Average),
			Min:       float32(batch[idx].Min),
			Max:       float32(batch[idx].Max),
			Start:     float32(batch[idx].Start),
			End:       float32(batch[idx].End),
			StartTime: batch[idx].StartTime,
			EndTime:   batch[idx].EndTime,
		}
	}

	return a.pub.SendDataBatch(stock, converted)
}

func (a *Adapter) CreateStockBroker(ctx context.Context, stock string) error {
	return a.conf.CreateTopic(ctx, stock)
}
