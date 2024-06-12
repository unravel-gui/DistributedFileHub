package loadbalancer

import (
	"DisHub/common"
	"DisHub/config"
	"sync"
)

type LBMap struct {
	lbm   map[common.NodeType]*LoadBalancer
	mutex sync.Mutex
}

var G_LoadBalancerMap = new(LBMap)

func (lbm *LBMap) Load() {
	lbm.lbm = make(map[common.NodeType]*LoadBalancer)
	storage, retry := config.GetLoadBalancerConfig()
	lbm.lbm[common.MASTER] = NewLoadBalancer(storage, retry)
	lbm.lbm[common.DATANODE] = NewLoadBalancer(storage, retry)
	lbm.lbm[common.APINODE] = NewLoadBalancer(storage, retry)
	lbm.lbm[common.USERNODE] = NewLoadBalancer(storage, retry)
}

func (lbMap *LBMap) ChangeStorage(storage int) {
	lbMap.mutex.Lock()
	defer lbMap.mutex.Unlock()
	for _, v := range lbMap.lbm {
		v.ChangeStorage(storage)
	}
	return
}

func (lbMap *LBMap) NextNode(nodeType common.NodeType) (bool, string) {
	lb, existed := lbMap.lbm[nodeType]
	if !existed {
		return false, ""
	}
	return true, lb.NextNode()
}

func (lbMap *LBMap) AddNode(nodeType common.NodeType, addr string, weight int) bool {
	lb, existed := lbMap.lbm[nodeType]
	if !existed {
		return false
	}
	lb.AddNode(addr, weight)
	return true
}

func (lbMap *LBMap) RemoveNode(nodeType common.NodeType, addr string) bool {
	lb, existed := lbMap.lbm[nodeType]
	if !existed {
		return false
	}
	lb.RemoveNode(addr)
	return true
}
