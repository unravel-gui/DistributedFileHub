package balancer

import (
	"DisHub/common"
	"DisHub/loadbalancer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func get(c *gin.Context) {
	nodeType := c.Param("nodeType")
	nt := common.StringToNodeType(nodeType)
	if nt == -1 {
		c.JSON(http.StatusBadRequest, "unknown nodeType")
		return
	}
	ok, server := loadbalancer.G_LoadBalancerMap.NextNode(nt)
	if !ok {
		c.JSON(http.StatusNotFound, fmt.Sprintf("no load balancer for nodeType: %s", nt.ToString()))
		return
	}
	c.JSON(http.StatusOK, server)
	return
}
