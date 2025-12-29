package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ishi-o/transfmt/config"
	"github.com/ishi-o/transfmt/internal/handler"
	"github.com/ishi-o/transfmt/internal/redis"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := redis.Init(); err != nil {
		log.Fatalf("Failed to init redis session store: %v", err)
	}
	r := gin.Default()
	r.Use(sessions.Sessions("transfmt_session", redis.SessionStore))
	setupRouter(r)
	addr := fmt.Sprintf(":%d", config.C.App.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server run failed: ", err)
	}
}

func setupRouter(r *gin.Engine) {
	rg := r.Group("/api")
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	handler.RegisterFanyiHandlers(rg)
}
