package loadbalancer

import (
	"fmt"
	"testing"
)

func TestRandomRobinLoadBalancer(t *testing.T) {
	lb := NewRandomRobinLoadBalancer()

	// 添加节点
	lb.AddNode("Node1", 1)
	lb.AddNode("Node2", 2)
	lb.AddNode("Node3", 3)

	// 进行一定数量的请求
	for i := 0; i < 10; i++ {
		selectedNode := lb.NextNode()
		fmt.Printf("Request %d goes to %s\n", i+1, selectedNode)
	}
	fmt.Println()
	// 添加新节点
	lb.AddNode("Node4", 4)
	for i := 0; i < 10; i++ {
		selectedNode := lb.NextNode()
		fmt.Printf("Request %d goes to %s\n", i+1, selectedNode)
	}
	fmt.Println()
	// 删除节点
	lb.RemoveNode("Node2")

	for i := 0; i < 10; i++ {
		selectedNode := lb.NextNode()
		fmt.Printf("Request %d goes to %s\n", i+1, selectedNode)
	}
}
