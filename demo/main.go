package main

import (
	"apigatewaydemo/demo/config"
	"apigatewaydemo/demo/gateway"
	"flag"
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		conf config.ApiGatewayConfig
	)

	confFilePath := flag.String("c", "config.json", "apigateway configuration file")
	flag.Parse()

	if err := config.ParseConfig(*confFilePath, &conf); err != nil {
		panic(err)
	}

	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		panic("parse seelog.xml error")
	}

	log.ReplaceLogger(logger)

	defer log.Flush()

	log.Info("Seelog Init Success!")

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
