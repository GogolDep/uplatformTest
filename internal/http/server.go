package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"uplatform/internal/config"
	"uplatform/internal/handlers"
	"uplatform/internal/services"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	exClient := services.NewExchangeClient(cfg.OXRBaseURL, cfg.OXRAppID, cfg.OXRTimeout)
	gifClient := services.NewGiphyClient(cfg.GiphyBaseURL, cfg.GiphyAPIKey, cfg.GiphyTimeout)

	h := handlers.NewGIFHandler(cfg, exClient, gifClient)

	v1 := r.Group("/v1")
	{
		v1.GET("/gif", h.GetGIF)
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
