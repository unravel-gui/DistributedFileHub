package loadbalancer

const (
	RANDOMROBIN int = iota
	ROUNDROBIN
)

type WeightedServer struct {
	Addr   string
	Weight int
}

type LoadBalancerI interface {
	NextNode() string
	AddNode(addr string, weight int)
	RemoveNode(addr string)
	AllNodes() []*WeightedServer
}

type LoadBalancer struct {
	lb      LoadBalancerI
	storage int
	retry   int
}

func NewLoadBalancer(storage, retry int) *LoadBalancer {
	lb := new(LoadBalancer)
	lb.storage = storage
	lb.retry = retry
	switch storage {
	case 0:
		lb.lb = NewRandomRobinLoadBalancer()
	case 1:
		lb.lb = NewRoundRobinLoadBalancer()
	default:
		return nil
	}
	return lb
}

func (lb LoadBalancer) ChangeStorage(storage int) bool {
	wss := lb.lb.AllNodes()
	switch storage {
	case 0:
		lb.lb = NewRandomRobinLoadBalancer()
	case 1:
		lb.lb = NewRoundRobinLoadBalancer()
	default:
		return false
	}
	for _, server := range wss {
		lb.lb.AddNode(server.Addr, server.Weight)
	}
	return true
}

func (lb *LoadBalancer) NextNode() string {
	return lb.lb.NextNode()
}

func (lb *LoadBalancer) AddNode(addr string, weight int) {
	lb.lb.AddNode(addr, weight)
}

func (lb *LoadBalancer) RemoveNode(addr string) {
	lb.lb.RemoveNode(addr)
}
