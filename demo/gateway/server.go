package gateway

import (
	"apigatewaydemo/demo/config"
	"apigatewaydemo/demo/dataload"
	"apigatewaydemo/demo/model"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io/ioutil"

	log "github.com/cihub/seelog"
)

type Server struct {
	srv *GoServ

	httpserver *fasthttp.Server

	router *fasthttprouter.Router

	dataload *dataload.DataLoad

	conf *config.ApiGatewayConfig

	gatelistenAddr string

	apiVersion int

	//API映射表解析数据
	cluster  map[string][]string
	gateApis map[string]model.GateAPI
}

func New(conf *config.ApiGatewayConfig) *Server {

	s := new(Server)

	s.srv = new(GoServ)

	s.conf = conf

	s.dataload = dataload.New(conf)

	s.gatelistenAddr = conf.ListenAddr + ":" + strconv.Itoa(conf.ListenPort)

	s.httpserver = &fasthttp.Server{
		Handler:                            s.mainHandler,
		ErrorHandler:                       nil,
		HeaderReceived:                     nil,
		ContinueHandler:                    nil,
		Name:                               "MyApiGateway",
		Concurrency:                        800,
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

	s.gateApis = make(map[string]model.GateAPI)
	s.cluster = make(map[string][]string)

	//初始化加载本地配置文件
	s.LoadLocalAPI()

	//动态更新配置
	go s.UpdateStorage()

	return s
}

func (s *Server) UpdateStorage() {

	for {

		isRestart := false

		lastApiVersion := s.apiVersion
		httpapis := s.dataload.GetApiMapping()

		if httpapis != nil && (httpapis.Version >= lastApiVersion+1) {

			api_content, err := json.Marshal(httpapis)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}

			s.apiVersion = httpapis.Version

			//写入本地apis文件并进行服务重启生效
			log.Infof("Apis is changed, Update and Save Apis.")
			ioutil.WriteFile("storage/local_apis.json", api_content, 0666)

			isRestart = true
		}

		if isRestart {
			//restart
			s.restartServer()
		}

		time.Sleep(time.Second * 5)
	}
}

func (s *Server) LoadLocalAPI() {

	var exAPIs model.RespAPIs

	data, err := ioutil.ReadFile("storage/local_apis.json")
	if err != nil {
		log.Error("Error:", err.Error())
	}
	json.Unmarshal(data, &exAPIs)

	for _, item := range exAPIs.APIs {
		s.gateApis[item.GatePath] = item
	}

	s.apiVersion = exAPIs.Version

	s.cluster = exAPIs.Clusters

}

func (s *Server) ApiMapping() {

	//LoadTestAPI
	s.router.Handle("GET", "/ping", s.EchoTest)

	//LoadAdminAPI
	s.router.Handle("GET", "/admin/restart", s.RestartTest)

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

func (s *Server) RestartTest(ctx *fasthttp.RequestCtx) {

	s.restartServer()

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

	req.SetRequestURI(back_real_path)
	// 设定网关请求的请求方法
	req.Header.SetMethodBytes(method)
	//设定请求Host
	req.Header.SetHost(hosts[0])

	fmt.Println(req)

	res := fasthttp.AcquireResponse()

	client := &fasthttp.Client{ReadTimeout: 10}

	if err := client.Do(req, res); err != nil {
		//log.Errorf("请求失败:%s", err.Error())
		//s.HandleAPIResponse(ctx, 1001, model.RespMsg[1001]["en"], nil, now2.Sub(now1).Seconds())
		return
	}

	/*
		ToDo: Demo Don't Need Content-Type && Gzip

			bodyType := res.Header.Peek("Content-Type")
			isJson := false
			if strings.Contains(string(bodyType), "application/json") {
				isJson = true
			}
	*/

	isJson := false
	resp := res.Body()

	if isJson {
		body := make(map[string]interface{})
		json.Unmarshal(resp, &body)
		s.HandleAPIResponse(ctx, 0, model.RespMsg[0]["en"], body)
	} else {
		s.HandleAPIResponse(ctx, 0, model.RespMsg[0]["en"], string(resp))
	}
}

func (s *Server) mainHandler(ctx *fasthttp.RequestCtx) {

	//1、pre
	next, err := BeforeRequestChain(ctx, s.conf.PluginList)

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

func (s *Server) start() {

	log.Info("Apigateway Start...")

	s.srv.Wrap(s.startServer)

}

func (s *Server) startServer() {

	log.Info("Apigateway ListenAndServe:" + s.gatelistenAddr)

	//初始路由映射
	s.ApiMapping()

	s.httpserver.ListenAndServe(s.gatelistenAddr)

}

func (s *Server) stopServer() {

	log.Info("Apigateway Shutdown...")

	s.httpserver.Shutdown()

}

func (s *Server) restartServer() {

	log.Info("Apigateway Restart...")

	s.stopServer()

	cmd := exec.Command("./control.sh", "restart")
	cmd.Start()

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
