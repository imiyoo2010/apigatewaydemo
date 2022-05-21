package model

//Api数据表
type Route struct {
	ID       int    	`json:"id"`
	Name     string 	`json:"name"`
	Method   string 	`json:"method"`
	Protocol string 	`json:"protocol"`

	UpstreamSign string  `json:"upstream_sign"`

	GatePath   string 	`json:"gate_path"`
	GateParams string 	`json:"gate_params"`

	BackPath     string `json:"back_path"`
	BackParams   string `json:"back_params"`
	BackPosition string `json:"back_position"`

	//业务字段
	Comment    	string `json:"-"`
	Is_del     	int    `json:"-"`
}

//Upstream数据表
type Upstream struct {
	ID           int    `json:"id"`
	UpstreamName string `json:"upstream_name"`
	UpstreamSign string `json:"upstream_sign"`
	UpstreamIp   string `json:"upstream_ip"`
	Protocol     string `json:"protocol"`

	//业务字段
	Comment      string `json:"-"`
}

//Version数据表
type Version struct {
	ID           	int    	`json:"id"`
	FileName 		string 	`json:"file_name"`
	VersionNum 		int 	`json:"version_num"`
}