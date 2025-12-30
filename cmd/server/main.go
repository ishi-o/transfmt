package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ishi-o/transfmt/internal/config"
	"github.com/ishi-o/transfmt/internal/handler"
	"github.com/ishi-o/transfmt/internal/logger"
	"github.com/ishi-o/transfmt/internal/redis"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load config %v", err)
	}
	logger.Init(&config.C.Log)
	defer logger.Sync()
	if err := redis.Init(); err != nil {
		logger.Fatal("Failed to init redis session store", zap.Error(err))
	}

	r := gin.Default()
	r.Use(sessions.Sessions("transfmt_session", redis.SessionStore))
	setupRouter(r)
	addr := fmt.Sprintf(":%d", config.C.App.Port)
	logger.Info("Starting server on ", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("Server run failed: %v", zap.Error(err))
	}
}

func setupRouter(r *gin.Engine) {
	rg := r.Group("/api")
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	handler.RegisterFanyiHandlers(rg)
}
