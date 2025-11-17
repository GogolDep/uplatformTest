package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"uplatform/internal/config"
	"uplatform/internal/services"
)

type GIFHandler struct {
	cfg      *config.Config
	exchange services.ExchangeClient
	giphy    services.GiphyClient
}

func NewGIFHandler(cfg *config.Config, ex services.ExchangeClient, gh services.GiphyClient) *GIFHandler {
	return &GIFHandler{cfg: cfg, exchange: ex, giphy: gh}
}

func (h *GIFHandler) GetGIF(c *gin.Context) {
	code := strings.ToUpper(c.DefaultQuery("code", "USD"))
	base := strings.ToUpper(h.cfg.BaseCurrency)
	if code == base {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code must differ from base currency"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 7*time.Second)
	defer cancel()

	todayRates, err := h.exchange.Latest(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	yesterday := time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02")
	yRates, err := h.exchange.Historical(ctx, yesterday)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	todayBase, ok1 := todayRates[base]
	todayCode, ok2 := todayRates[code]
	yBase, ok3 := yRates[base]
	yCode, ok4 := yRates[code]
	if !(ok1 && ok2 && ok3 && ok4) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown currency code"})
		return
	}

	curToBaseToday := todayBase / todayCode
	curToBaseYesterday := yBase / yCode

	query := "broke"
	if curToBaseToday > curToBaseYesterday {
		query = "rich"
	}

	gifURL, err := h.giphy.RandomFromSearch(ctx, query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, gifURL)
}
