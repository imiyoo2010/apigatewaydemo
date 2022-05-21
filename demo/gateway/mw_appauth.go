package gateway

import (
	"apigatewaydemo/demo/model"
	"apigatewaydemo/demo/utils"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"io/ioutil"
)

type AppAuth struct {
}

func init() {
	Regist("appauth", &AppAuth{})
}

func (a *AppAuth) ProcessRequest(ctx *fasthttp.RequestCtx, conf interface{}) (int, error) {

	// head params
	appId := string(ctx.Request.Header.Peek("AppId"))
	token := string(ctx.Request.Header.Peek("Token"))
	if appId == "" || token == "" {
		return 0, errors.New("AppId/AppSecret is empty")
	}

	// get appSecret
	appAuthList := LoadAppAuth()
	var appSecret string
	for _, v := range appAuthList {
		if v.AppId == appId {
			appSecret = v.AppSecret
			break
		}
	}

	if utils.CheckSign(appId, appSecret, token) {
		return 1, nil //执行后续插件
	} else {
		return 0, errors.New("AppId/AppSecret is not valid") //终止执行后续
	}
}

func (a *AppAuth) Name() string {
	return "appauth"
}

func LoadAppAuth() []model.AuthInfo {
	var r []model.AuthInfo

	data, _ := ioutil.ReadFile("storage/local_auth.json")
	json.Unmarshal(data, &r)
	return r
}
