// Package handler
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ishi-o/transfmt/internal/redis"
	"github.com/ishi-o/transfmt/pkg/response"
)

func RegisterFanyiHandlers(rg *gin.RouterGroup) {
	gp := rg.Group("fanyi/")
	{
		gp.GET("setting", getSetting)
		gp.POST("setting/lang", setLang)
		gp.POST("", fanyi)
	}
}

func getSetting(c *gin.Context) {
	response.Success(c, redis.GetUserSetting(sessions.Default(c)))
}

type langRequest struct {
	SourceLang string `json:"sourceLang" form:"sourceLang" binding:"required"`
	TargetLang string `json:"targetLang" form:"targetLang" binding:"required"`
}

func setLang(c *gin.Context) {
	var lang langRequest
	if err := c.ShouldBind(&lang); err != nil {
		c.Error(err)
		return
	}
	session := sessions.Default(c)
	setting := redis.GetUserSetting(session)
	setting.Fanyi.SourceLang, setting.Fanyi.TargetLang = lang.SourceLang, lang.TargetLang
	if err := redis.SaveUserSetting(session, setting); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

func fanyi(c *gin.Context) {
	var text string
	c.ShouldBind(&text)
	setting := redis.GetUserSetting(sessions.Default(c))
	backend, srclang, dstlang := setting.Fanyi.Backend, setting.Fanyi.SourceLang, setting.Fanyi.TargetLang
	switch backend {
	case "yandex":
		dsttext, err := fanyiByYandex(srclang, dstlang, text)
		if err != nil {
			response.ServerError(c, err)
			return
		}
		response.Success(c, dsttext)
	default:
		response.ServerError(c, fmt.Errorf("internal error"))
	}
}

func fanyiByYandex(srclang, dstlang, text string) (any, error) {
	form := url.Values{}
	form.Set("text", text)
	rawUUID := uuid.New().String()
	ucid := strings.ReplaceAll(rawUUID, "-", "")
	form.Set("ucid", ucid)

	if srclang == "auto" {
		// TODO: add auto transfer
	}

	form.Set("lang", fmt.Sprintf("%s-%s", srclang, dstlang))
	form.Set("srv", "android")
	form.Set("format", "text")

	req, _ := http.NewRequest("POST", "https://translate.yandex.net/api/v1/tr.json/translate", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Text []string `json:"text"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	translatedText := ""
	if len(result.Text) > 0 {
		translatedText = result.Text[0]
	}
	return translatedText, nil
}
