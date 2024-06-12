package loadbalancer

import (
	"DisHub/common/utils"
	"math/rand"
	"sync"
)

// RandomRobinLoadBalancer  表示轮询负载均衡器
type RandomRobinLoadBalancer struct {
	nodes       []*WeightedServer
	mutex       sync.Mutex
	weights     []int
	totalWeight int
	rand        *rand.Rand
}

// NewRandomRobinLoadBalancer 创建一个新的负载均衡器
func NewRandomRobinLoadBalancer() *RandomRobinLoadBalancer {
	return &RandomRobinLoadBalancer{
		rand: rand.New(rand.NewSource(utils.GetNow().UnixNano())),
	}
}

func (lb *RandomRobinLoadBalancer) initWeightCache() {
	lb.weights = nil
	sum := 0
	for _, node := range lb.nodes {
		sum += node.Weight
		lb.weights = append(lb.weights, sum)
	}
}

// AddNode 添加一个节点到负载均衡器中
func (lb *RandomRobinLoadBalancer) AddNode(addr string, weight int) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.nodes = append(lb.nodes, &WeightedServer{
		Addr:   addr,
		Weight: weight,
	})
	lb.totalWeight += weight
	lb.initWeightCache()
}

// RemoveNode 从负载均衡器中删除一个节点
func (lb *RandomRobinLoadBalancer) RemoveNode(addr string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	for i, node := range lb.nodes {
		if node.Addr == addr {
			lb.totalWeight -= node.Weight
			lb.nodes = append(lb.nodes[:i], lb.nodes[i+1:]...)
			lb.initWeightCache()
			return
		}
	}
}

// NextNode 返回下一个被选中的节点
func (lb *RandomRobinLoadBalancer) NextNode() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.nodes) == 0 {
		return ""
	}

	target := lb.rand.Intn(lb.totalWeight)

	for i, node := range lb.nodes {
		if lb.weights[i] >= target {
			return node.Addr
		}
	}

	return lb.nodes[0].Addr
}

func (lb *RandomRobinLoadBalancer) AllNodes() []*WeightedServer {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.nodes
}
