package main

import (
	"apigatewaydemo/ch4/gateway"
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	gateway.StartServer()

	//http.ListenAndServe("8089",nil)

	Signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, Signals...)

	select {
	case s := <-ch:
		log.Infof("received signal %s: terminating", s)
	}

}
