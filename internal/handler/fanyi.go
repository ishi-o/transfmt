// Package handler
package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ishi-o/transfmt/internal/redis"
	"github.com/ishi-o/transfmt/pkg/response"
)

func RegisterFanyiHandlers(rg *gin.RouterGroup) {
	gp := rg.Group("fanyi/")
	{
		gp.GET("setting", getSetting)
	}
}

func getSetting(c *gin.Context) {
	response.Success(c, redis.GetUserSetting(sessions.Default(c)))
}
