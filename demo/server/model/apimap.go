package model

//网关端的路由映射交互数据
type ApiMap struct {
	Version int                    `json:"version"`
	Apis    []Api                  `json:"apis"`
	Hash    string                 `json:"hash"`
	Cluster map[string]interface{} `json:"clusters"`
}

//网关端的API接口结构
type Api struct {
	ID       int       `json:"id"`
	Method   []string  `json:"method"`
	GatePath string    `json:"gate_path"`
	UpStream string    `json:"upstream"`
	Service  string    `json:"service"`
	BackPath string    `json:"back_path"`
	Params   []Param   `json:"params"`
}

//网关端的参数设置
type Param struct {
	Gate     string `json:"gate"`
	Back     string `json:"back"`
	Position string `json:"position"`
}