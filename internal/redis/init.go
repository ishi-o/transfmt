// Package redis
package redis

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/ishi-o/transfmt/config"
)

var SessionStore sessions.Store

func Init() error {
	cfg := &config.C.Redis
	appCfg := &config.C.App

	if len(appCfg.SecretKey) < 32 {
		return fmt.Errorf("secret key too short, need at least 32 characters, got %d", len(appCfg.SecretKey))
	}

	store, err := redis.NewStore(
		cfg.PoolSize,
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		"",
		cfg.Password,
		[]byte(appCfg.SecretKey),
	)
	if err != nil {
		return fmt.Errorf("failed to create redis store: %w", err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionTTL,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	SessionStore = store
	return nil
}
