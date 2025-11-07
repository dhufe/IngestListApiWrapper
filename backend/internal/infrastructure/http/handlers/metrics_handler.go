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
	// Metriken aktualisieren bevor Prometheus sie scraped
	if err := h.service.UpdateMetrics(c.Request.Context()); err != nil {
		// Fehler loggen, aber trotzdem Metriken ausgeben
		// (andere Metriken wie HTTP-Requests funktionieren ja noch)
		c.Error(err) // Gin's Error-Handler
	}

	hp := promhttp.Handler()
	hp.ServeHTTP(c.Writer, c.Request)
}
