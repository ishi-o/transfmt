// Package redis
package redis

import (
	"encoding/json"
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

func GetUserSetting(session sessions.Session) config.UserSetting {
	if data, ok := session.Get("user_setting").([]byte); ok && len(data) > 0 {
		var setting config.UserSetting
		if json.Unmarshal(data, &setting) == nil {
			return setting
		}
	}
	defaultSetting := config.C.UserSetting
	SaveUserSetting(session, defaultSetting)
	return defaultSetting
}

func SaveUserSetting(session sessions.Session, setting config.UserSetting) error {
	data, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	session.Set("user_setting", data)
	return session.Save()
}
