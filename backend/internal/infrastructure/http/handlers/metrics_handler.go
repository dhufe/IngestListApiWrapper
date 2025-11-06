package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
	c.JSON(http.StatusOK, nil)
}
