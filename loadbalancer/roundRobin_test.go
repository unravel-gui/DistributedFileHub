package loadbalancer

import (
	"fmt"
	"testing"
)

func TestRoundRobinLoadBalancer(t *testing.T) {
	loadBalancer := NewRoundRobinLoadBalancer()

	// 添加节点
	loadBalancer.AddNode("Node1", 1)
	loadBalancer.AddNode("Node2", 2)
	loadBalancer.AddNode("Node3", 3)

	// 进行一定数量的请求
	for i := 0; i < 10; i++ {
		selectedNode := loadBalancer.NextNode()
		fmt.Printf("Request %d goes to %s\n", i+1, selectedNode)
	}
	fmt.Println()
	// 添加新节点
	loadBalancer.AddNode("Node4", 4)

	for i := 0; i < 10; i++ {
		selectedNode := loadBalancer.NextNode()
		fmt.Printf("Request %d goes to %s\n", i, selectedNode)
	}
	fmt.Println()
	// 删除节点
	loadBalancer.RemoveNode("Node2")

	fmt.Println()
	// 进行一定数量的请求
	for i := 0; i < 10; i++ {
		selectedNode := loadBalancer.NextNode()
		fmt.Printf("Request %d goes to %s\n", i, selectedNode)
	}
}
