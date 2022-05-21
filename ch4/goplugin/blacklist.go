package main

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Blacklist struct {
}

func (w *Blacklist) ProcessRequest(ctx *fasthttp.RequestCtx, conf interface{}) (int, error) {

	fmt.Println("Blacklist Plugin Checking...")

	flag := true

	/*
		remote_addr := ctx.RemoteIP()

		if remote_addr in net.IPs blacklist {
			flag = false

			fmt.Printf("Client Blacklist IP: %s",remote_addr.String())
		}
	*/

	if flag {
		return 1, nil //执行后续插件
	} else {
		return 0, errors.New("Blacklist IP") //终止执行后续
	}
}

func (w *Blacklist) Name() string {
	return "Blacklist"
}

var BlacklistPlugin Blacklist
