package gateway

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type Server struct {
	httpserver *fasthttp.Server
}

func New() *Server {

	s := new(Server)

	s.httpserver = &fasthttp.Server{
		Handler: s.mainHandler,

		Name: "MyApiGateway",

		DisableKeepalive: true,
		ReadTimeout:      10,
		WriteTimeout:     10,
	}

	return s
}

func (s *Server) mainHandler(ctx *fasthttp.RequestCtx) {

	//1、pre

	//2、router

	//3、post

	//4、log

	ctx.Response.SetBodyString("Hello,World!")

}

func (s *Server) startServer() {

	//s.httpserver.ListenAndServeUNIX("apigate.sock",0666)

	fmt.Println("MyApiGateway ListenAndServe: :8080")

	s.httpserver.ListenAndServe(":8080")

}

func (s *Server) stopServer() {

	s.httpserver.Shutdown()

}

func StartServer() {
	server := New()
	server.startServer()
}
