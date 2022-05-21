package main

import (
	"apigatewaydemo/ch1/gateway"
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	gateway.StartServer()
	
	//通过监听管道，避免程序退出
	Signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, Signals...)

	select {
	case s := <-ch:
		log.Infof("received signal %s: terminating", s)
	}

}
