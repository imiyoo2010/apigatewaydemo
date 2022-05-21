package gateway

import (
	"apigatewaydemo/ch4/model"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io/ioutil"
)

type Server struct {
	httpserver *fasthttp.Server

	router *fasthttprouter.Router

	//API映射表解析数据
	cluster  map[string][]string
	gateApis map[string]model.GateAPI
}

func New() *Server {

	s := new(Server)

	s.httpserver = &fasthttp.Server{
		Handler:                            s.mainHandler,
		ErrorHandler:                       nil,
		HeaderReceived:                     nil,
		ContinueHandler:                    nil,
		Name:                               "MyApiGateway",
		Concurrency:                        0,
		ReadBufferSize:                     0,
		WriteBufferSize:                    0,
		ReadTimeout:                        0,
		WriteTimeout:                       0,
		IdleTimeout:                        0,
		MaxConnsPerIP:                      0,
		MaxRequestsPerConn:                 0,
		MaxKeepaliveDuration:               0,
		MaxIdleWorkerDuration:              0,
		TCPKeepalivePeriod:                 0,
		MaxRequestBodySize:                 0,
		DisableKeepalive:                   false,
		TCPKeepalive:                       false,
		ReduceMemoryUsage:                  false,
		GetOnly:                            false,
		DisablePreParseMultipartForm:       false,
		LogAllErrors:                       false,
		SecureErrorLogMessage:              false,
		DisableHeaderNamesNormalizing:      false,
		SleepWhenConcurrencyLimitsExceeded: 0,
		NoDefaultServerHeader:              false,
		NoDefaultDate:                      false,
		NoDefaultContentType:               false,
		KeepHijackedConns:                  false,
		CloseOnShutdown:                    false,
		StreamRequestBody:                  false,
		ConnState:                          nil,
		Logger:                             nil,
		TLSConfig:                          nil,
	}

	s.router = fasthttprouter.New()

	return s
}

func (s *Server) LoadLocalAPI() {

	var exAPIs model.RespAPIs

	data, err := ioutil.ReadFile("storage/local_apis.json")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	json.Unmarshal(data, &exAPIs)

	for _, item := range exAPIs.APIs {
		s.gateApis[item.GatePath] = item
	}
}

func (s *Server) ApiMapping() {

	//LoadTestAPI
	s.router.Handle("GET", "/ping", s.EchoTest)

	//LoadLocalAPI
	for _, item := range s.gateApis {
		for _, method := range item.Method {
			s.router.Handle(method, item.GatePath, s.ApiHandler)
		}
	}
}

func (s *Server) EchoTest(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.SetBody([]byte("pong"))

}

func (s *Server) ApiHandler(ctx *fasthttp.RequestCtx) {

	//1、构造网关发送的请求

	//1.1 路径映射，根据网关路径和映射关系获取后端路径
	gate_path := string(ctx.URI().Path())
	back_api, ok := s.gateApis[gate_path]
	if !ok {
		return
	}

	//1.2 参数映射，根据网关参数和映射关系获取后端参数
	/*
		1. QueryString参数处理
	*/
	args := ctx.URI().QueryArgs()
	for _, item := range back_api.Params {
		if item.Position == "querystring" {
			if args.Has(item.Gate_param) {
				argsValue := args.Peek(item.Gate_param)
				args.Del(item.Gate_param)
				args.AddBytesV(item.Back_param, argsValue)
			}
		}
	}

	req := &ctx.Request
	res := fasthttp.AcquireResponse()

	var client *fasthttp.Client

	//var c *fasthttp.LBClient

	if err := client.Do(req, res); err != nil {
		//log.Errorf("请求失败:%s", err.Error())
		s.HandleAPIResponse(ctx, 1001, model.RespMsg[1001]["en"], nil)
		return
	}

	resp := res.Body()

	s.HandleAPIResponse(ctx, 0, model.RespMsg[0]["en"], resp)
}

func (s *Server) mainHandler(ctx *fasthttp.RequestCtx) {

	//1、pre
	enable_pluginlist := []string{"Blacklist", "appauth"}

	next, err := BeforeRequestChain(ctx, enable_pluginlist)

	if !next {
		HandleBeforeResponse(ctx, 1001, model.RespMsg[1001]["en"]+": "+err, "")
		return
	}

	/*
		2.Router & Post: 路由映射和请求执行
	*/
	s.router.Handler(ctx)

	//4、error&log

}

func (s *Server) startServer() {

	//log.Info("ListenAndServe: :8088")
	s.httpserver.ListenAndServe(":8080")

}

func (s *Server) stopServer() {

	s.httpserver.Shutdown()

}

func StartServer() {
	server := New()
	server.startServer()
}
