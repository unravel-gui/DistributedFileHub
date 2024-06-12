package master

import (
	"DisHub/master/heartbeat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func get(c *gin.Context) {
	ds := heartbeat.GetMasterServers()
	c.JSON(http.StatusOK, ds)
}
