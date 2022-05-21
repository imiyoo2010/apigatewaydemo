package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type WhiteList struct {

}

func (w *WhiteList) ProcessRequest(ctx *fasthttp.RequestCtx, conf interface{}) (int, error) {
	fmt.Println("Plugin WhiteList")

	flag := true

	if flag {
		return 1, nil
	}else{
		return 0,nil
	}
}


func (w *WhiteList) Name() string {
	return "whitelist"
}

var WhiteListPlugin WhiteList