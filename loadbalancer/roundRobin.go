package loadbalancer

import (
	"sync"
)

// RoundRobinLoadBalancer  表示轮询负载均衡器
type RoundRobinLoadBalancer struct {
	nodes        []*WeightedServer
	mutex        sync.Mutex
	totalWeight  int
	currentIndex int
	currentRound int
}

// NewRoundRobinLoadBalancer 创建一个新的负载均衡器
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{}
}

// AddNode 添加一个节点到负载均衡器中
func (lb *RoundRobinLoadBalancer) AddNode(addr string, weight int) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.nodes = append(lb.nodes, &WeightedServer{
		Addr:   addr,
		Weight: weight,
	})
	lb.totalWeight += weight
}

// RemoveNode 从负载均衡器中删除一个节点
func (lb *RoundRobinLoadBalancer) RemoveNode(addr string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	for i, node := range lb.nodes {
		if node.Addr == addr {
			lb.totalWeight -= node.Weight
			lb.nodes = append(lb.nodes[:i], lb.nodes[i+1:]...)
			return
		}
	}
}

// NextNode 返回下一个被选中的节点
func (lb *RoundRobinLoadBalancer) NextNode() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.nodes) == 0 {
		return ""
	}

	for {
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.nodes)
		if lb.currentIndex == 0 {
			lb.currentRound++
			if lb.currentRound >= lb.totalWeight {
				lb.currentRound = 0
			}
		}

		node := lb.nodes[lb.currentIndex]
		if node.Weight >= lb.currentRound {
			return node.Addr
		}
	}
}

func (lb *RoundRobinLoadBalancer) AllNodes() []*WeightedServer {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.nodes
}
