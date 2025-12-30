package redis

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/ishi-o/transfmt/internal/config"
)

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
