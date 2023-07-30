package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


type Server struct {
	srv *http.Server
}


var (
	instance *Server
	once 	   sync.Once
)

func New(c *resolver.ConfigMap) *Server {

	once.Do(func() {

		port, err := c.GetStringKey("PORT")
		if err != nil {
			panic("fetch server metric port is required")
		}

		router := gin.Default()

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		}

		router.GET("/metrics", gin.WrapH(promhttp.Handler()))

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()

		instance = &Server{
			srv: srv,
		}
	})

	return instance
}


func (s *Server) Close() {
	s.srv.Shutdown(context.Background())
}