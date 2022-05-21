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

type AuthInfo struct {
	Id        int    `json:"id"`
	AppId     string `json:"appid"`
	AppSecret string `json:"appsecret"`
}

type RespAPIs struct {
	Version  int                 `json:"version"`
	APIs     []GateAPI           `json:"apis"`
	Clusters map[string][]string `json:"clusters"`
}

//网关日志
type GateLog struct {
	RequestID       string  `json:"request_id"`
	Clientip        string  `json:"client_ip"`
	Url             string  `json:"url"`
	Status          int     `json:"status"`
	UpstreamHost    string  `json:"upstream_host"`
	UpstreamUri     string  `json:"upstream_uri"`
	Upstream_status int     `json:"upstream_status"`
	Upstream_time   float64 `json:"upstream_time"`
}
