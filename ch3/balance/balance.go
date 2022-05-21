package balance

import (
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
)

type WeightRoundRobinBalance struct {
	IsHttps bool

	curIndex int

	rss []*WeightNode
}

type WeightNode struct {
	weight int // 配置的权重，即在配置文件或初始化时约定好的每个节点的权重

	currentWeight int //节点当前权重，会一直变化

	effectiveWeight int //有效权重，初始值为weight, 通讯过程中发现节点异常，则-1 ，之后再次选取本节点，调用成功一次则+1，直达恢复到weight 。 用于健康检查，处理异常节点，降低其权重。

	addr string // 服务器addr

	client *fasthttp.HostClient //当前节点的连接
}

func (r *WeightRoundRobinBalance) isHealthy(Node *WeightNode) bool {
	//利用网络连接的状态来判断健康状态
	return true
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("params len need 2")
	}
	addr := params[0]
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	hostclient := &fasthttp.HostClient{
		Addr:  addr,
		IsTLS: r.IsHttps,
	}

	node := &WeightNode{
		weight:          int(parInt),
		effectiveWeight: int(parInt), // 初始化時有效权重 = 配置权重值
		currentWeight:   int(parInt), // 初始化時当前权重 = 配置权重值
		addr:            addr,
		client:          hostclient,
	}

	r.rss = append(r.rss, node)

	return nil
}

func (r *WeightRoundRobinBalance) Next() (*fasthttp.HostClient, string) {
	if len(r.rss) == 0 {
		return nil, ""
	}

	totalWeight := 0
	var maxWeightNode *WeightNode
	for key, node := range r.rss {
		//step 1 统计所有有效权重之和
		totalWeight += node.effectiveWeight

		//step 2 变更节点当前权重为的节点当前权重+节点有效权重
		node.currentWeight += node.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if r.isHealthy(node) {
			if node.effectiveWeight < node.weight {
				node.effectiveWeight++
			}

		} else {
			node.effectiveWeight--
		}

		//step 4 选择最大临时权重点节点
		if maxWeightNode == nil || maxWeightNode.currentWeight < node.currentWeight {
			maxWeightNode = node
			r.curIndex = key
		}
	}

	//step 5 变更临时权重为 临时权重-有效权重之和
	maxWeightNode.currentWeight -= totalWeight

	//step 6 返回当前选择的client
	return maxWeightNode.client, maxWeightNode.addr
}
