package main

import (
	"log"

	"uplatform/internal/config"
	"uplatform/internal/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	server.RegisterRoutes(r, cfg)

	if err := r.Run(cfg.HTTPAddr()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
