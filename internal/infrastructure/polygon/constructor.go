package polygon

import (
	"os"
	"github.com/polygon-io/client-go/websocket/models"
	polygonws "github.com/polygon-io/client-go/websocket"
)


var (
	POLYGON_API_KEY = os.Getenv("POLYGON_API_KEY")

	DEFAULT_BUFFER_SIZE = 1000
)


type PolygonWs struct {
	conn *polygonws.Client

	ch chan models.EquityAgg
}

func NewPolygonWs() *PolygonWs {

	c, err := polygonws.New(polygonws.Config{
		APIKey:    POLYGON_API_KEY,
		Feed:      polygonws.RealTime,
		Market:    polygonws.Stocks,
	})
	
	if err != nil {
		panic(err)
	}

	if err := c.Connect(); err != nil {
		panic(err)
	}
	
	return &PolygonWs{conn: c}
}






