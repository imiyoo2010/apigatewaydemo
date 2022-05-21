package gateway

import (
	"apigatewaydemo/ch3/balance"
	"apigatewaydemo/ch3/model"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strings"
	"time"
)

type Server struct {
	httpserver *fasthttp.Server

	router *fasthttprouter.Router

	hostclient map[string]*fasthttp.LBClient

	hostweightclient map[string]*balance.WeightRoundRobinBalance

	//API映射表解析数据
	cluster  map[string][]string
	gateApis map[string]model.GateAPI
}

func New() *Server {

	//初始化
	s := new(Server)

	s.gateApis = make(map[string]model.GateAPI)
	s.cluster = make(map[string][]string)

	s.hostclient = make(map[string]*fasthttp.LBClient)
	s.hostweightclient = make(map[string]*balance.WeightRoundRobinBalance)

	//初始化加载本地配置文件
	s.LoadLocalAPI()

	//s.makeHostClient()

	s.makeWeightHostClient()

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

	fmt.Println(exAPIs)

	for _, item := range exAPIs.APIs {
		s.gateApis[item.GatePath] = item
	}

	s.cluster = exAPIs.Clusters

	fmt.Println(s.gateApis)
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

func (s *Server) makeHostClient() {
	//cluster集群名称，用来确定几组upstream，后端可以创建hostclient
	for _, v := range s.gateApis {
		var isHttps bool

		if v.Service == "https" {
			isHttps = true
		} else {
			isHttps = false
		}

		for _, addr1 := range s.cluster[v.Upstream] {

			var lbc fasthttp.LBClient
			c := &fasthttp.HostClient{
				Addr:  addr1,
				IsTLS: isHttps,
			}

			lbc.Clients = append(lbc.Clients, c)
			lbc.Timeout = 30 * time.Second

			s.hostclient[v.Upstream] = &lbc
		}
	}
}

func (s *Server) makeWeightHostClient() {
	//cluster集群名称，用来确定几组upstream，后端可以创建hostclient
	for _, v := range s.gateApis {
		var isHttps bool

		if v.Service == "https" {
			isHttps = true
		} else {
			isHttps = false
		}

		for _, weightaddr := range s.cluster[v.Upstream] {
			/*
				    用分号来分割和增加权重信息
					{"cluster":{"test":["10.10.10.10:8080;4"]}}
			*/

			item := strings.Split(weightaddr, ";")

			balan := balance.WeightRoundRobinBalance{IsHttps: isHttps}

			if len(item) > 1 {
				addr := item[0]
				weight := item[1]
				balan.Add(addr, weight)
			} else {
				balan.Add(item[0], "1")
			}

			s.hostweightclient[v.Upstream] = &balan
		}
	}
}

func (s *Server) ApiHandler(ctx *fasthttp.RequestCtx) {

	//1、构造网关发送的请求

	//1.1 路径映射，根据网关路径和映射关系获取后端路径
	gate_path := string(ctx.URI().Path())
	back_api, ok := s.gateApis[gate_path]
	if !ok {
		return
	}

	hosts := s.cluster[back_api.Upstream]
	fmt.Println(hosts)

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

	var (
		back_real_path string
	)

	argsQS := args.QueryString()

	method := ctx.Method()

	if len(argsQS) > 0 {
		back_real_path = back_api.BackPath + "?" + string(argsQS)
	} else {
		back_real_path = back_api.BackPath
	}

	req := &ctx.Request

	//设定网关请求的完整路径(带参数)
	req.SetRequestURI(back_real_path)
	// 设定网关请求的请求方法
	req.Header.SetMethodBytes(method)

	res := &ctx.Response

	//client := s.hostclient[back_api.Upstream]
	//加权负载client
	client, addr := s.hostweightclient[back_api.Upstream].Next()

	//设定请求Host
	//req.Header.SetHost(hosts[0])
	req.Header.SetHost(addr)

	if err := client.Do(req, res); err != nil {
		fmt.Printf("请求失败:%s", err.Error())
		s.HandleAPIResponse(ctx, 1001, model.RespMsg[1001]["en"], nil)
		return
	}

	resp := res.Body()

	s.HandleAPIResponse(ctx, 0, model.RespMsg[0]["en"], string(resp))
}

func (s *Server) mainHandler(ctx *fasthttp.RequestCtx) {

	//1、pre

	/*
		2.Router & Post: 路由映射和请求执行
	*/
	s.router.Handler(ctx)

	//4、log

}

func (s *Server) startServer() {

	//log.Info("ListenAndServe: :8088")

	//初始路由映射
	s.ApiMapping()

	s.httpserver.ListenAndServe(":8088")

}

func (s *Server) stopServer() {

	s.httpserver.Shutdown()

}

func StartServer() {
	server := New()
	server.startServer()
}
