package balance

import (
	"testing"
)

func TestWeightBalance(t *testing.T) {

	w := WeightRoundRobinBalance{}

	//利用不同的域名来模拟加权轮询
	w.Add("www.imiyoo.com","4")
	w.Add("i.imiyoo.com.","2")

	for i:=0; i<6; i++ {

		_, addr := w.Next()

		t.Log("当前请求的域名"+ addr)
	}
}
