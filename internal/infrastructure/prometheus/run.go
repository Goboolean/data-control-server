package prometheus

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)






func Run(ctx context.Context) {

	port, exist := os.LookupEnv("FETCH_SERVER_METRIC_PORT")

	if !exist {
		panic(fmt.Errorf("%s required", "FETCH_SERVER_METRIC_PORT"))
	}

	router := gin.Default()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go func() {
		if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
			panic(err)
		}
	}()
}
