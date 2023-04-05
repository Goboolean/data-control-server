package aggregate

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Goboolean/data-control-server/internal/domain/valueobj"
)


const (
	defaultChannelSize = 100
)

type StockDataRelayQueue struct {
	channelList map[int] chan valueobj.StockAggregate
}

func NewStockDataRelayQueue() *StockDataRelayQueue {
	rand.Seed(time.Now().UnixNano())

	instance := &StockDataRelayQueue{}
	instance.channelList = make(map[int] chan valueobj.StockAggregate)
	return instance
}

func (s *StockDataRelayQueue) SubscribeChannel() (chan valueobj.StockAggregate) {
	newChan := make(chan valueobj.StockAggregate, defaultChannelSize)
	newChanId := rand.Intn(100000)
	s.channelList[newChanId] = newChan

	return newChan
}

func (s *StockDataRelayQueue) UnsubscribeChannel(channel chan valueobj.StockAggregate) error {
	for k, v := range s.channelList {
		if v == channel {
			delete(s.channelList, k)
			return nil
		}
	}

	return errors.New("chan not found")
}


func (s *StockDataRelayQueue) PushDataOnQueue(data *valueobj.StockAggregate) {
	for i := 0; i < len(s.channelList); i++ {
		s.channelList[i] <- *data
	}
}