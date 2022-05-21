package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

func (s *Server) HandleAPIResponse(ctx *fasthttp.RequestCtx, code int, msg string, resp interface{}) {

	c := make(map[string]interface{})

	c["code"] = code
	c["msg"] = msg
	c["data"] = resp

	gbody, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error marshalling body:", err.Error())
	}

	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentType("application/json;charset=UTF-8")
	ctx.Response.SetBody(gbody)
}
