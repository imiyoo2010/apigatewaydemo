package main

import (
	"apigatewaydemo/ch5/config"
	"apigatewaydemo/ch5/gateway"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/cihub/seelog"
)

func main() {

	var (
		conf config.ApiGatewayConfig
	)

	confFilePath := flag.String("c", "./config.json", "apigateway configuration file")

	flag.Parse()

	if err := config.ParseConfig(*confFilePath, &conf); err != nil {
		panic(err)
	}

	gateway.StartServer(&conf)
	//http.ListenAndServe("8089",nil)

	Signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, Signals...)

	select {
	case s := <-ch:
		log.Infof("received signal %s: terminating", s)
	}
}
