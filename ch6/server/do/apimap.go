package do

import (
	"apigatewaydemo/ch5/server/data"
	"apigatewaydemo/ch5/server/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)


var (
	db  *data.Sqlite
)


func init() {
	db = data.NewSqlite()
}

func GetApiMap(c *gin.Context) {


	rs 	:=db.GetRouter()

	us 	:=db.GetUpstreamp()


	cluster_map :=make(map[string][]string)

	for _, item :=range us {
		if item.UpstreamSign=="" {
			continue
		}
		var addr_list []string
		addr_list = append(addr_list,item.UpstreamIp)
		cluster_map[item.UpstreamSign]=addr_list
	}


	c.JSON(200,gin.H{

		"version":"1",
		"clusters":cluster_map,
		"apis":rs,

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

	c.JSON(200, gin.H{
		"code":0,
		"msg": lastid,
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

	c.JSON(200, gin.H{
		"code":0,
		"msg": lastid,
	})
}