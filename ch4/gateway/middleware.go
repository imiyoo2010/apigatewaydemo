package gateway

import (
	"apigatewaydemo/ch4/utils"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"path"
	"plugin"
	"strings"
)

type Middlerware interface {
	ProcessRequest(r *fasthttp.RequestCtx, conf interface{}) (int, error) // 对HTTP请求进行处理和判定
	Name() string                                                         // 插件的名称
	//Priority() int 														// 插件的优先级
}

//插件数据库
var mwmaps map[string][]Middlerware

//插件数组
//var mwlist []Middlerware

func init() {

	mwmaps = make(map[string][]Middlerware)

	//mwlist = make([]Middlerware,100)

	solist, err := utils.ListDir("goplugin", "so")

	if err != nil {
		fmt.Printf("goplugin error:%s\n", err)
		return
	}
	for _, sofile := range solist {

		//获取文件的后缀(文件类型)
		fileType := path.Ext(sofile)
		//获取文件名称(不带后缀)
		fileNameOnly := strings.TrimSuffix(sofile, fileType)

		plugin_path := "goplugin" + string(os.PathSeparator) + sofile

		open, err := plugin.Open(plugin_path)

		if err != nil {
			panic(err)
		}

		plugin_inst := strings.ToUpper(fileNameOnly[:1]) + fileNameOnly[1:] + "Plugin"

		whitelistplugin, err := open.Lookup(plugin_inst)

		if mwplugin, ok := whitelistplugin.(Middlerware); ok {

			Regist(mwplugin.Name(), mwplugin)
		}
	}

}

func Regist(name string, mwinstance Middlerware) {

	//mwlist = append(mwlist[:index],append([]Middlerware{mwinstance},mwlist[index:]...)...)

	mwmaps[name] = append(mwmaps[name], mwinstance)
}
