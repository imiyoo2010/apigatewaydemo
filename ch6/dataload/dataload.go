package dataload

import (
	"apigatewaydemo/ch6/config"
	"apigatewaydemo/ch6/dbstore"
	"apigatewaydemo/ch6/model"
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"time"
)

type DataLoad struct {
	ApiServerDomain string

	conf *config.ApiGatewayConfig

	dataclient *http.Client

	logCollector *dbstore.ESOutput

	LogChan chan *model.GateLog
}

func New(conf *config.ApiGatewayConfig) *DataLoad {

	d := new(DataLoad)

	//d.ScgServerDomain = domain

	d.conf = conf

	d.dataclient = &http.Client{}

	d.LogChan = make(chan *model.GateLog, 1000)

	d.logCollector = dbstore.NewESClient(conf.ESUrl)

	go d.SendGateLog()

	return d
}

func (d *DataLoad) SendGateLog() {

	for {
		select {
		case apilog := <-d.LogChan:
			fmt.Println(apilog)
			msg, _ := json.Marshal(apilog)
			d.logCollector.CollectApiLog(string(msg))
		}

		time.Sleep(time.Second * 5)
	}
}

func (d *DataLoad) GetApiMapping() *model.RespAPIs {
	//秒级更新，每隔5秒进行处理
	var gateapis model.RespAPIs

	url := d.conf.ApimapConfig

	res, err := d.dataclient.Get(url)

	if err != nil {
		log.Errorf("GetApiMapping err:%s", err)
		return &gateapis
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode == 200 {
		json.Unmarshal(body, &gateapis)
	}

	return &gateapis
}
