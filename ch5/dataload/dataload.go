package dataload

import (
	"apigatewaydemo/ch5/config"
	"apigatewaydemo/ch5/model"
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
)

type DataLoad struct {
	ApiServerDomain string

	conf *config.ApiGatewayConfig

	dataclient *http.Client
}

func New(conf *config.ApiGatewayConfig) *DataLoad {

	d := new(DataLoad)

	//d.ScgServerDomain = domain

	d.conf = conf

	d.dataclient = &http.Client{}

	return d
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