package do

import (
	"apigatewaydemo/demo/server/data"
	"apigatewaydemo/demo/server/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

var (
	db *data.Sqlite
)

func init() {
	db = data.NewSqlite()
}

func GetApiMap(c *gin.Context) {

	rs := db.GetRouter()
	us := db.GetUpstreamp()

	cluster_map := make(map[string][]string)

	var api_gate_list []model.Api

	for _, api_talbe := range rs {

		api_gate := model.Api{}

		api_gate.ID = api_talbe.ID
		api_gate.UpStream = api_talbe.UpstreamSign
		api_gate.Service = api_talbe.Protocol
		api_gate.GatePath = api_talbe.GatePath
		api_gate.BackPath = api_talbe.BackPath
		api_gate.Method = append(api_gate.Method, strings.Split(api_talbe.Method, ",")...)

		gateParams := strings.Split(api_talbe.GateParams, ",")
		backParams := strings.Split(api_talbe.BackParams, ",")
		positions := strings.Split(api_talbe.BackPosition, ",")
		for i := 0; i < len(gateParams); i++ {
			param := &model.Param{}
			param.Gate = gateParams[i]
			param.Back = backParams[i]
			param.Position = positions[i]
			api_gate.Params = append(api_gate.Params, *param)
		}

		api_gate_list = append(api_gate_list, api_gate)
	}

	for _, item := range us {
		if item.UpstreamSign == "" {
			continue
		}
		var addr_list []string
		addr_list = append(addr_list, item.UpstreamIp)
		cluster_map[item.UpstreamSign] = addr_list
	}

	c.JSON(200, gin.H{
		"version":  db.GetVersion("api_map"),
		"clusters": cluster_map,
		"apis":     api_gate_list,
	})
}

func AddApiMap(c *gin.Context) {

	//获取POST数据
	body := c.Request.Body

	requestJson := &model.Route{}

	data, err := ioutil.ReadAll(body)

	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据接受失败, 请检查是否格式正确",
		})
		return
	}
	json.Unmarshal(data, &requestJson)

	fmt.Println(requestJson.Name)

	lastid := db.AddRouter(requestJson)

	db.UpdateVersion("api_map")

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  lastid,
	})
}

func AddUpstream(c *gin.Context) {

	//获取POST数据
	body := c.Request.Body

	requestJson := &model.Upstream{}

	data, err := ioutil.ReadAll(body)

	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据接受失败, 请检查是否格式正确",
		})
		return
	}
	json.Unmarshal(data, &requestJson)

	fmt.Println(requestJson.UpstreamSign)

	lastid := db.AddUpstreamp(requestJson)

	db.UpdateVersion("api_map")

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  lastid,
	})
}
