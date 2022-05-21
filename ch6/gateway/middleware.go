package gateway

import "github.com/valyala/fasthttp"

type Middlerware interface {
	ProcessRequest(r *fasthttp.RequestCtx, conf interface{}) (int, error)    //  对HTTP请求进行处理和判定
	Name() string                                                            //   插件的名称
}


var mwmaps map[string][]Middlerware


func init() {
	mwmaps = make(map[string][]Middlerware)

	//Regist("test",mwinstance)
}


func Regist(name string, mwinstance Middlerware) {
	mwmaps[name] = append(mwmaps[name], mwinstance)
}