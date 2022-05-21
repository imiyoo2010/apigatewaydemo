package gateway

import (
	"apigatewaydemo/ch6/config"
	"apigatewaydemo/ch6/dataload"
	"apigatewaydemo/ch6/model"
	"apigatewaydemo/ch6/utils"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type Server struct {
	srv *GoServ

	httpserver *fasthttp.Server

	router *fasthttprouter.Router

	dataload *dataload.DataLoad

	gatelistenAddr string

	//API映射表解析数据
	cluster  map[string][]string
	gateApis map[string]model.GateAPI
}

func New(conf *config.ApiGatewayConfig) *Server {

	s := new(Server)

	s.srv = new(GoServ)

	s.httpserver = &fasthttp.Server{
		Handler:                            s.mainHandler,
		ErrorHandler:                       nil,
		HeaderReceived:                     nil,
		ContinueHandler:                    nil,
		Name:                               "MyApiGateway",
		Concurrency:                        1024,
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
		DisableKeepalive:                   true,
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

	s.gatelistenAddr = conf.ListenAddr + ":" + strconv.Itoa(conf.ListenPort)

	s.dataload = dataload.New(conf)

	s.gateApis = make(map[string]model.GateAPI)
	s.cluster = make(map[string][]string)

	//初始化加载本地配置文件
	s.LoadLocalAPI()

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
	/*	ToDo:
		2. PostData参数处理
		3、认证参数&请求头参数
	*/
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
	//resp := &ctx.Response
	/*
		req,res := fasthttp.AcquireRequest(),fasthttp.AcquireResponse()
		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(res)
		}()
	*/

	req.SetRequestURI(back_real_path)
	// 设定网关请求的请求方法
	req.Header.SetMethodBytes(method)
	//设定请求Host
	req.Header.SetHost(hosts[0])

	fmt.Println(req)

	res := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}

	now2 := time.Now()

	if err := client.Do(req, res); err != nil {
		//log.Errorf("请求失败:%s", err.Error())
		s.HandleAPIResponse(ctx, 1001, model.RespMsg[1001]["en"], nil)
		return
	}

	resp := res.Body()

	now3 := time.Now()

	var apilog model.GateLog

	apilog.Clientip = ctx.RemoteIP().String()
	apilog.RequestID = utils.GenUUID()
	apilog.Status = res.StatusCode()
	apilog.Url = gate_path
	apilog.UpstreamHost = hosts[0]
	apilog.UpstreamUri = back_real_path
	apilog.Upstream_status = res.StatusCode()
	apilog.Upstream_time = now3.Sub(now2).Seconds()

	fmt.Println(apilog)

	select {

	case s.dataload.LogChan <- &apilog:
		log.Infof("make apilog record:%s", apilog.RequestID)

	default:
		log.Error("apilog record error")

	}

	s.HandleAPIResponse(ctx, 0, model.RespMsg[0]["en"], resp)
}

func (s *Server) mainHandler(ctx *fasthttp.RequestCtx) {

	//1、pre

	//next, err := BeforeRequestChain(ctx, p.conf.PluginList, p.conf)

	//if !next {
	//	HandleBeforeResponse(ctx, 1001, model.RespMsg[1001]["en"]+": "+err)
	//	return
	//}

	/*
		2.Router & Post: 路由映射和请求执行
	*/
	s.router.Handler(ctx)

	//4、error&log

}

func (s *Server) start() {

	log.Info("Apigateway Start...")

	s.srv.Wrap(s.startServer)

}

func (s *Server) startServer() {
	log.Info("ListenAndServe:" + s.gatelistenAddr)

	//初始路由映射
	s.ApiMapping()

	s.httpserver.ListenAndServe(s.gatelistenAddr)

}

func (s *Server) stopServer() {

	s.httpserver.Shutdown()

}

type GoServ struct {
	sync.WaitGroup
}

func (s *GoServ) Wrap(h func()) {
	s.Add(1)
	go func() {
		h()
		s.Done()
	}()
}

func StartServer(conf *config.ApiGatewayConfig) {
	server := New(conf)
	server.start()
}
