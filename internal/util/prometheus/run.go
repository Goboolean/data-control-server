package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	srv *http.Server
}

func New(c *resolver.ConfigMap) (*Server, error) {

	var err error

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	router := gin.Default()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go func() {
		if webErr := srv.ListenAndServe(); webErr != nil && webErr != http.ErrServerClosed {
			err = webErr
		}
	}()
	time.Sleep(100 * time.Millisecond)

	return &Server{
		srv: srv,
	}, err
}

func (s *Server) Close() {
	s.srv.Shutdown(context.Background())
}