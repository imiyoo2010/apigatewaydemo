package model

//网关接口
type GateAPI struct {
	ID int `json:"id"`

	Method   []string `json:"method"`
	GatePath string   `json:"gate_path"`

	Upstream string `json:"upstream"`
	Service  string `json:"service"`
	BackPath string `json:"back_path"`

	Params []Param `json:"params"`
}

type Param struct {
	Gate_param string `json:"gate"`
	Back_param string `json:"back"`
	Position   string `json:"position"`
}

type RespAPIs struct {
	APIs     []GateAPI           `json:"apis"`
	Clusters map[string][]string `json:"clusters"`
}

type AuthInfo struct {
	Id        int    `json:"id"`
	AppId     string `json:"appid"`
	AppSecret string `json:"appsecret"`
}
