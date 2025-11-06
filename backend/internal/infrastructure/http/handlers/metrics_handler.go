package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"
)

type MetricsHandler struct {
	service *services.MetricsService
}

func NewMetricsHandler(service *services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		service: service,
	}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	hp := promhttp.Handler()
	hp.ServeHTTP(c.Writer, c.Request)
}
