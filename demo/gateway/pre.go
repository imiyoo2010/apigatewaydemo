package gateway

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func BeforeRequestChain(ctx *fasthttp.RequestCtx, enable_pluginList []string) (bool, string) {

	var run_mwlist []Middlerware

	for _, name := range enable_pluginList {

		run_mwlist = mwmaps[name]

		for _, mw := range run_mwlist {

			fmt.Printf("Uri: %s , Middlerware Name: %s", string(ctx.Request.RequestURI()), name)

			result, err := mw.ProcessRequest(ctx, "")

			if result == 0 {
				fmt.Printf("Uri: %s , Middlerware Name: %s, The Result is False, Stop Next Middlerware!", string(ctx.Request.RequestURI()), name)
				return false, fmt.Sprintf("middlerware: %s, error: %s", name, err)
			}
		}
	}

	return true, fmt.Sprintf("middlerware: %s, error: %s", "success", "检测通过")
}
