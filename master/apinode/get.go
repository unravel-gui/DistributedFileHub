package apinode

import (
	"DisHub/master/heartbeat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func get(c *gin.Context) {
	ds := heartbeat.GetApiServers()
	c.JSON(http.StatusOK, ds)
}
