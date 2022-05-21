package main

import (
	"apigatewaydemo/ch5/server/do"
	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	c.String(200, "Hello World!")
}

func main() {
	// 1.创建路由
	r := gin.Default()
	r.TrustedPlatform = gin.PlatformCloudflare
	// 2.绑定路由规则，执行的函数
	//网关接口

	r.GET("/gateway/apimap", do.GetApiMap)

	r.POST("/gateway/apimap", do.AddApiMap)

	r.POST("/gateway/upstream", do.AddUpstream)

	//前端接口

	r.Run(":8088")
}
