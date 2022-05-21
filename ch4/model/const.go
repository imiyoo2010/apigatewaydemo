package model

var RespMsg = map[int]map[string]string{

	//执行成功
	0: {
		"en": "success",
		"cn": "成功",
	},

	//网关错误
	1001: {
		"en": "pre-request error",
		"cn": "请求预处理阶段错误",
	},

	1002: {
		"en": "post-request error",
		"cn": "请求后端执行阶段错误",
	},

	1003: {
		"en": "Reverse reflection error",
		"cn": "反参映射配置文件错误",
	},

	//后端错误
	2001: {
		"en": "backend error",
		"cn": "后端错误",
	},
}
